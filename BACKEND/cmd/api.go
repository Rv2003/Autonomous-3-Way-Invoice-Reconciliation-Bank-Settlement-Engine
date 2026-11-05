package main

import (
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	repo "github.com/Rv2003/Autonomous-3-Way-Invoice-Reconciliation-Bank-Settlement-Engine/BACKEND/internal/adapters/postgresql/sqlc"

	"github.com/Rv2003/Autonomous-3-Way-Invoice-Reconciliation-Bank-Settlement-Engine/BACKEND/internal/email"
	"github.com/Rv2003/Autonomous-3-Way-Invoice-Reconciliation-Bank-Settlement-Engine/BACKEND/internal/pdf"

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

	r.Group(func(r chi.Router) {
		r.Use(app.requireRole("administrator"))
		r.Get("/users", app.listUsers)
		r.Post("/users", app.createUser)
		r.Patch("/users/{id}/role", app.updateUserRole)
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

	pdfService := pdf.NewService(repo.New(app.db))
	pdfHandler := pdf.NewHandler(pdfService)
	emailService := email.NewService(repo.New(app.db))
	emailHandler := email.NewHandler(emailService)
	r.Get("/pdf", pdfHandler.ListPdf)
	r.Post("/pdf/submit", pdfHandler.CreatePdf)

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

			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			w.Header().Set("Access-Control-Allow-Credentials", "true")
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
		WriteTimeout: 30 * time.Second,
		ReadTimeout:  10 * time.Second,
		IdleTimeout:  time.Minute,
	}
	log.Printf("Server has started %s", app.config.addr)
	return srv.ListenAndServe()
}
