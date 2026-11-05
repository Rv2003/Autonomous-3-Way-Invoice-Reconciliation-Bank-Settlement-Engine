package main

import (
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	repo "github.com/Rv2003/Autonomous-3-Way-Invoice-Reconciliation-Bank-Settlement-Engine/BACKEND/internal/adapters/postgresql/sqlc"
	"github.com/Rv2003/Autonomous-3-Way-Invoice-Reconciliation-Bank-Settlement-Engine/BACKEND/internal/env"
	"github.com/Rv2003/Autonomous-3-Way-Invoice-Reconciliation-Bank-Settlement-Engine/BACKEND/internal/pdf"
	"github.com/Rv2003/Autonomous-3-Way-Invoice-Reconciliation-Bank-Settlement-Engine/BACKEND/internal/telegram"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
)

type application struct {
	config config
	db     *pgxpool.Pool
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

	// A good base middleware stack
	r.Use(middleware.RequestID)
	r.Use(middleware.ClientIPFromRemoteAddr)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(localFrontendCORS)

	// Set a timeout value on the request context (ctx)
	r.Use(middleware.Timeout(60 * time.Second))

	// Public route: users can log in without an existing token.
	r.Post("/auth/login", app.login)

	// Protected routes: a valid JWT is required.
	r.Group(func(r chi.Router) {
		r.Use(app.authenticate)
		r.Get("/auth/me", app.me)
	})

	// RBAC test routes
	r.Group(func(r chi.Router) {
		// Only administrators can access this route.
		r.With(app.requireRole("administrator")).Get(
			"/admin/test",
			func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte("Administrator access granted"))
			},
		)

		// Both administrators and clerks can access this route.
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

	pdfService := pdf.NewService(repo.New(app.db))
	pdfHandler := pdf.NewHandler(pdfService)
	r.Get("/pdf", pdfHandler.ListPdf)
	r.Post("/pdf/submit", pdfHandler.CreatePdf)

	// Telegram Intake Webhook Handlers
	botToken := env.GetString("TELEGRAM_BOT_TOKEN", "")
	tgClient := telegram.NewClient(botToken)
	tgHandler := &telegram.Handler{
		BotToken:   botToken,
		UploadDir:  env.GetString("UPLOAD_DIR", "uploads"),
		Client:     tgClient,
		Store:      telegram.NewPgStore(app.db),
		PdfService: pdfService,
	}

	r.Post("/webhooks/telegram", tgHandler.Receive)
	r.Post("/webhooks/telegram/", tgHandler.Receive)
	r.Get("/telegram/feed", tgHandler.GetFeed)

	return r
}

func localFrontendCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if origin != "" {
			parsed, err := url.Parse(origin)
			if err != nil || (parsed.Scheme != "http" && parsed.Scheme != "https") ||
				parsed.User != nil || parsed.Path != "" || parsed.RawQuery != "" || parsed.Fragment != "" ||
				!isLoopbackHost(parsed.Hostname()) {
				http.Error(w, "origin is not allowed", http.StatusForbidden)
				return
			}

			w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
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
	return strings.EqualFold(host, "localhost") || host == "127.0.0.1" || host == "::1"
}

func (app *application) run(h http.Handler) error {
	srv := &http.Server{
		Addr:         app.config.addr,
		Handler:      h,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute,
	}
	log.Printf("Server has started %s", app.config.addr)
	return srv.ListenAndServe()
}
