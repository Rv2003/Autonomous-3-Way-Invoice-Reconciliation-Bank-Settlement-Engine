package main

import (
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/Rv2003/Autonomous-3-Way-Invoice-Reconciliation-Bank-Settlement-Engine/BACKEND/internal/email"
	"github.com/Rv2003/Autonomous-3-Way-Invoice-Reconciliation-Bank-Settlement-Engine/BACKEND/internal/env"
	"github.com/Rv2003/Autonomous-3-Way-Invoice-Reconciliation-Bank-Settlement-Engine/BACKEND/internal/gmail"
	"github.com/Rv2003/Autonomous-3-Way-Invoice-Reconciliation-Bank-Settlement-Engine/BACKEND/internal/inbound"
	"github.com/Rv2003/Autonomous-3-Way-Invoice-Reconciliation-Bank-Settlement-Engine/BACKEND/internal/pdf"
	"github.com/Rv2003/Autonomous-3-Way-Invoice-Reconciliation-Bank-Settlement-Engine/BACKEND/internal/telegram"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
)

type application struct {
	config         config
	db             *pgxpool.Pool
	gmailClient    *gmail.Client
	emailService   email.Service
	pdfService     pdf.Service
	inboundService inbound.Service
	emailProcessor *email.Processor
}

type config struct {
	addr string
	db   dbConfig
}

type dbConfig struct {
	dsn string
}

func (app *application) mount() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.ClientIPFromRemoteAddr)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(localFrontendCORS)

	r.Use(middleware.Timeout(60 * time.Second))

	r.Post("/auth/login", app.login)
	r.Post("/auth/logout", app.logout)

	r.Group(func(r chi.Router) {
		r.Use(app.authenticate)

		r.Get("/auth/me", app.me)
		r.Get("/auth/profile", app.me)
		r.Patch("/auth/profile", app.updateProfile)
		r.Patch("/auth/password", app.changePassword)
	})

	emailHandler := email.NewHandler(
		app.emailService,
		app.emailProcessor,
		app.inboundService,
	)

	r.Post("/webhooks/gmail", emailHandler.GmailWebhook)
	r.Get("/email/feed", emailHandler.GetFeed)
	r.Get("/email/mailbox", emailHandler.GetMailbox)
	r.Post("/email/sync", emailHandler.SyncInbox)

	inboundHandler := inbound.NewHandler(app.inboundService)
	r.Get("/inbound/{id}/file", inboundHandler.ViewDocument)
	r.Delete("/inbound/{id}", inboundHandler.DeleteDocument)

	// Telegram Intake Webhook & Feed Handlers
	botToken := env.GetString("TELEGRAM_BOT_TOKEN", "")
	tgClient := telegram.NewClient(botToken)
	tgHandler := &telegram.Handler{
		BotToken:   botToken,
		UploadDir:  env.GetString("UPLOAD_DIR", "uploads"),
		Client:     tgClient,
		Store:      telegram.NewPgStore(app.db),
		PdfService: app.pdfService,
	}

	r.Post("/webhooks/telegram", tgHandler.Receive)
	r.Post("/webhooks/telegram/", tgHandler.Receive)
	r.Get("/telegram/feed", tgHandler.GetFeed)

	r.Group(func(r chi.Router) {
		r.Use(app.requireRole("administrator"))

		r.Get("/users", app.listUsers)
		r.Post("/users", app.createUser)
		r.Patch("/users/{id}/role", app.updateUserRole)

		r.Post("/email/mailboxes", emailHandler.CreateMailbox)
		r.Get(
			"/email/ingestion-events",
			emailHandler.ListIngestionEvents,
		)
		r.Patch(
			"/email/ingestion-events/{id}/processing",
			emailHandler.ProcessIngestionEvent,
		)
		r.Patch(
			"/email/ingestion-events/{id}/processed",
			emailHandler.CompleteIngestionEvent,
		)
		r.Patch(
			"/email/ingestion-events/{id}/failed",
			emailHandler.FailIngestionEvent,
		)
	})

	r.Group(func(r chi.Router) {
		r.With(app.requireRole("administrator")).Get(
			"/admin/test",
			func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte("Administrator access granted"))
			},
		)

		r.With(app.requireRole("administrator", "clerk")).Get(
			"/staff/test",
			func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte("Staff access granted"))
			},
		)
	})

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hi"))
	})

	pdfHandler := pdf.NewHandler(app.pdfService)

	r.Get("/pdf", pdfHandler.ListPdf)
	r.Get("/pdf/{id}/file", pdfHandler.ViewPdf)
	r.Post("/pdf/submit", pdfHandler.CreatePdf)
	r.Delete("/pdf/{id}", pdfHandler.DeletePdf)

	return r
}

func localFrontendCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")

		if origin != "" {
			parsed, err := url.Parse(origin)

			if err != nil ||
				(parsed.Scheme != "http" &&
					parsed.Scheme != "https") ||
				parsed.User != nil ||
				parsed.Path != "" ||
				parsed.RawQuery != "" ||
				parsed.Fragment != "" ||
				!isLoopbackHost(parsed.Hostname()) {

				http.Error(
					w,
					"origin is not allowed",
					http.StatusForbidden,
				)
				return
			}

			w.Header().Set(
				"Access-Control-Allow-Origin",
				origin,
			)

			w.Header().Set(
				"Access-Control-Allow-Methods",
				"GET, POST, PATCH, DELETE, OPTIONS",
			)

			w.Header().Set(
				"Access-Control-Allow-Headers",
				"Content-Type, Authorization",
			)

			w.Header().Set(
				"Access-Control-Allow-Credentials",
				"true",
			)

			w.Header().Add("Vary", "Origin")
		}

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func isLoopbackHost(host string) bool {
	return strings.EqualFold(host, "localhost") ||
		host == "127.0.0.1" ||
		host == "::1"
}

func (app *application) run(h http.Handler) error {
	srv := &http.Server{
		Addr:         app.config.addr,
		Handler:      h,
		WriteTimeout: 30 * time.Second,
		ReadTimeout:  10 * time.Second,
		IdleTimeout:  time.Minute,
	}

	log.Printf(
		"Server has started %s",
		app.config.addr,
	)

	return srv.ListenAndServe()
}