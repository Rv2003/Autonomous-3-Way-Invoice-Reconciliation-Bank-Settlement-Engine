package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/Rv2003/Autonomous-3-Way-Invoice-Reconciliation-Bank-Settlement-Engine/BACKEND/internal/adapters/postgresql/migrations"
	repo "github.com/Rv2003/Autonomous-3-Way-Invoice-Reconciliation-Bank-Settlement-Engine/BACKEND/internal/adapters/postgresql/sqlc"
	"github.com/Rv2003/Autonomous-3-Way-Invoice-Reconciliation-Bank-Settlement-Engine/BACKEND/internal/email"
	"github.com/Rv2003/Autonomous-3-Way-Invoice-Reconciliation-Bank-Settlement-Engine/BACKEND/internal/env"
	"github.com/Rv2003/Autonomous-3-Way-Invoice-Reconciliation-Bank-Settlement-Engine/BACKEND/internal/gmail"
	"github.com/Rv2003/Autonomous-3-Way-Invoice-Reconciliation-Bank-Settlement-Engine/BACKEND/internal/inbound"
	"github.com/Rv2003/Autonomous-3-Way-Invoice-Reconciliation-Bank-Settlement-Engine/BACKEND/internal/pdf"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
)

func main() {
	ctx := context.Background()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	if _, err := jwtSecret(); err != nil {
		logger.Error("authentication is not configured", "error", err)
		os.Exit(1)
	}

	cfg := config{
		addr: ":8080",
		db: dbConfig{
			dsn: env.GetString(
				"GOOSE_DBSTRING",
				"host=localhost user=invoice_app password=password123 dbname=invoice_reconciliation sslmode=disable",
			),
		},
	}

	pool, err := pgxpool.New(ctx, cfg.db.dsn)
	if err != nil {
		logger.Error(
			"failed to create database connection pool",
			"error",
			err,
		)
		os.Exit(1)
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		logger.Error(
			"failed to connect to database",
			"error",
			err,
		)
		os.Exit(1)
	}

	logger.Info("connected to database")

	migrationDB := stdlib.OpenDBFromPool(pool)
	defer migrationDB.Close()

	if err := migrations.Up(ctx, migrationDB); err != nil {
		logger.Error(
			"failed to apply database migrations",
			"error",
			err,
		)
		os.Exit(1)
	}

	gmailClient, err := gmail.NewClient(ctx)
	if err != nil {
		logger.Error(
			"failed to initialize Gmail client",
			"error",
			err,
		)
		os.Exit(1)
	}

	repository := repo.New(pool)

	emailService := email.NewService(repository)

	// Ensure default project mailbox exists in PostgreSQL
	adminMailbox := env.GetString("GMAIL_MAILBOX_ADDRESS", "awirbse.devtest@gmail.com")
	_, _ = emailService.CreateMailbox(ctx, repo.CreateEmailMailboxParams{
		Provider:     "gmail",
		EmailAddress: adminMailbox,
	})

	pdfService := pdf.NewService(repository)
	inboundService := inbound.NewService(repository)

	emailProcessor := email.NewProcessor(
		emailService,
		gmailClient,
		pdfService,
		inboundService,
	)

	api := application{
		config:         cfg,
		db:             pool,
		gmailClient:    gmailClient,
		emailService:   emailService,
		pdfService:     pdfService,
		inboundService: inboundService,
		emailProcessor: emailProcessor,
	}

	if err := api.run(api.mount()); err != nil {
		slog.Error(
			"Server has failed to start",
			"error",
			err,
		)
		os.Exit(1)
	}
}