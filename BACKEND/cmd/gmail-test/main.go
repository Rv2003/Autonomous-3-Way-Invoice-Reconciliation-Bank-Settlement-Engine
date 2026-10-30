package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
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
}
