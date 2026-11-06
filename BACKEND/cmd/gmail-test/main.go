package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"

	_ "github.com/jackc/pgx/v5/stdlib"
)

const (
	credentialsFile = "credentials/client_secret_841501621403-f8qd9akhvh1f68le2p953rq0rtjmrv60.apps.googleusercontent.com.json"
	tokenFile       = "tokens/gmail-token.json"
)

func main() {
	ctx := context.Background()

	credentials, err := os.ReadFile(credentialsFile)
	if err != nil {
		log.Fatalf("failed to read OAuth credentials: %v", err)
	}

	config, err := google.ConfigFromJSON(
		credentials,
		gmail.GmailReadonlyScope,
	)
	if err != nil {
		log.Fatalf("failed to parse OAuth credentials: %v", err)
	}

	tokenData, err := os.ReadFile(tokenFile)
	if err != nil {
		log.Fatalf("failed to read Gmail token: %v", err)
	}

	var token oauth2.Token
	if err := json.Unmarshal(tokenData, &token); err != nil {
		log.Fatalf("failed to parse Gmail token: %v", err)
	}

	client := config.Client(ctx, &token)

	service, err := gmail.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("failed to create Gmail service: %v", err)
	}

	profile, err := service.Users.GetProfile("me").Do()
	if err != nil {
		log.Fatalf("failed to access Gmail API: %v", err)
	}

	fmt.Println("Gmail API connection successful.")
	fmt.Printf("Authenticated account: %s\n", profile.EmailAddress)

	dbString := os.Getenv("GOOSE_DBSTRING")
	if dbString == "" {
		log.Fatal("GOOSE_DBSTRING is not set")
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
