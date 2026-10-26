package main

import (
	"log"
	"net/http"
	"time"

	"github.com/Rv2003/Autonomous-3-Way-Invoice-Reconciliation-Bank-Settlement-Engine/BACKEND/internal/pdf"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type application struct {
	config config
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
  r.Use(middleware.ClientIPFromRemoteAddr) // pick one ClientIPFrom* based on your infra, see below
  r.Use(middleware.Logger)
  r.Use(middleware.Recoverer)

  // Set a timeout value on the request context (ctx), that will signal
  // through ctx.Done() that the request has timed out and further
  // processing should be stopped.
  r.Use(middleware.Timeout(60 * time.Second))

  r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("hi"))
  })
	//http.ListenAndServe(":3333", r)
pdfHandler:=pdf.NewHandler(nil)
r.Get("/pdf",pdfHandler.ListPdf)
	
return r
}


func (app *application)  run(h http.Handler) error{
srv:= &http.Server{
	Addr:app.config.addr,
	Handler: h,
	WriteTimeout: time.Second*30,
	ReadTimeout: time.Second*10,
	IdleTimeout: time.Minute,
}
log.Printf("Server has started %s", app.config.addr)
return srv.ListenAndServe()

}