package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/Rv2003/Autonomous-3-Way-Invoice-Reconciliation-Bank-Settlement-Engine/BACKEND/internal/env"
	"github.com/Rv2003/Autonomous-3-Way-Invoice-Reconciliation-Bank-Settlement-Engine/BACKEND/internal/gmailoauth"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	ctx := context.Background()

	service, err := gmailoauth.NewService(ctx)
	if err != nil {
		log.Fatalf("failed to configure Gmail API: %v", err)
	}

	profile, err := service.Users.GetProfile("me").Do()
	if err != nil {
		log.Fatalf("failed to access Gmail API: %v", err)
	}

	fmt.Println("Gmail API connection successful.")
	fmt.Printf("Authenticated account: %s\n", profile.EmailAddress)

	dbString := env.GetString("GOOSE_DBSTRING", "")
	if dbString == "" {
		log.Fatal("GOOSE_DBSTRING is not set; add it to BACKEND/.env or set it in the environment")
	}

	db, err := sql.Open("pgx", dbString)
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	var (
		mailboxID string
		email     string
		status    string
	)

	err = db.QueryRowContext(
		ctx,
		`SELECT id, email_address, status
		 FROM email_mailboxes
		 WHERE provider = 'gmail'
		   AND email_address = $1`,
		profile.EmailAddress,
	).Scan(&mailboxID, &email, &status)

	if err != nil {
		if err == sql.ErrNoRows {
			log.Fatalf("Gmail account is not registered in email_mailboxes: %s", profile.EmailAddress)
		}

		log.Fatalf("failed to query email mailbox: %v", err)
	}

	fmt.Println("Mailbox found in PostgreSQL.")
	fmt.Printf("Mailbox ID: %s\n", mailboxID)
	fmt.Printf("Email address: %s\n", email)
	fmt.Printf("Status: %s\n", status)
}
