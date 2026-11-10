package main

import (
	"context"
	"log"
	"strconv"
	"time"

	repo "github.com/Rv2003/Autonomous-3-Way-Invoice-Reconciliation-Bank-Settlement-Engine/BACKEND/internal/adapters/postgresql/sqlc"
	"github.com/Rv2003/Autonomous-3-Way-Invoice-Reconciliation-Bank-Settlement-Engine/BACKEND/internal/email"
	"github.com/Rv2003/Autonomous-3-Way-Invoice-Reconciliation-Bank-Settlement-Engine/BACKEND/internal/env"
	"github.com/Rv2003/Autonomous-3-Way-Invoice-Reconciliation-Bank-Settlement-Engine/BACKEND/internal/gmail"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	gmailProvider = "gmail"
	topicName     = "projects/awirbse-email-integration/topics/gmail-email-events"
)

func main() {
	ctx := context.Background()

	// Connect to PostgreSQL using the existing project configuration.
	dsn := env.GetString(
		"GOOSE_DBSTRING",
		"host=localhost user=invoice_app password=password123 dbname=invoice_reconciliation sslmode=disable",
	)

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		log.Fatalf("failed to create database connection pool: %v", err)
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	log.Println("Connected to PostgreSQL.")

	// Reuse the existing email service.
	repository := repo.New(pool)
	emailService := email.NewService(repository)

	// Reuse the existing Gmail OAuth/client implementation.
	gmailClient, err := gmail.NewClient(ctx)
	if err != nil {
		log.Fatalf("failed to create Gmail client: %v", err)
	}

	// Identify the authenticated Gmail account.
	profile, err := gmailClient.Service().Users.GetProfile("me").Do()
	if err != nil {
		log.Fatalf("failed to get Gmail profile: %v", err)
	}

	log.Printf("Authenticated Gmail account: %s", profile.EmailAddress)

	// Find the mailbox registered in PostgreSQL.
	mailbox, err := emailService.GetMailboxByEmail(
		ctx,
		gmailProvider,
		profile.EmailAddress,
	)
	if err != nil {
		log.Fatalf(
			"failed to find Gmail mailbox %s in PostgreSQL: %v",
			profile.EmailAddress,
			err,
		)
	}

	log.Printf("Mailbox found: %s", mailbox.EmailAddress)

	// Register Gmail push notifications through Google Cloud Pub/Sub.
	watch, err := gmailClient.StartWatch(ctx, topicName)
	if err != nil {
		log.Fatalf("failed to start Gmail watch: %v", err)
	}

	// Gmail returns the expiration as Unix milliseconds.
	expiresAt := time.UnixMilli(watch.Expiration)

	// Gmail history IDs are identifiers, not timestamps.
	historyID := strconv.FormatUint(watch.HistoryID, 10)

	// Persist the Gmail watch state in PostgreSQL.
	updatedMailbox, err := emailService.UpdateMailboxWatch(
		ctx,
		mailbox.ID,
		historyID,
		pgtype.Timestamptz{
			Time:  expiresAt,
			Valid: true,
		},
	)
	if err != nil {
		log.Fatalf("failed to persist Gmail watch state: %v", err)
	}

	log.Println("Gmail watch started successfully.")
	log.Printf("Mailbox: %s", updatedMailbox.EmailAddress)
	log.Printf("History ID: %s", updatedMailbox.HistoryID.String)
	log.Printf(
		"Watch expires: %s",
		updatedMailbox.WatchExpiresAt.Time.Format(time.RFC3339),
	)
	log.Printf("Pub/Sub topic: %s", topicName)
}
