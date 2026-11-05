package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
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

	token, err := authorize(ctx, config)
	if err != nil {
		log.Fatalf("authorization failed: %v", err)
	}

	if err := saveToken(tokenFile, token); err != nil {
		log.Fatalf("failed to save token: %v", err)
	}

	fmt.Println()
	fmt.Println("Gmail authorization completed successfully.")
	fmt.Printf("Token saved to: %s\n", tokenFile)
}

func authorize(ctx context.Context, config *oauth2.Config) (*oauth2.Token, error) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return nil, fmt.Errorf("failed to create local listener: %w", err)
	}
	defer listener.Close()

	config.RedirectURL = "http://" + listener.Addr().String()

	authURL := config.AuthCodeURL(
		"awirbse-gmail-auth",
		oauth2.AccessTypeOffline,
	)

	fmt.Println("Open the following URL in your browser:")
	fmt.Println()
	fmt.Println(authURL)
	fmt.Println()
	fmt.Println("Sign in using: awirbse.devtest@gmail.com")
	fmt.Println()
	fmt.Println("Waiting for Google's callback...")

	result := make(chan *oauth2.Token, 1)
	errResult := make(chan error, 1)

	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("error") != "" {
			errResult <- fmt.Errorf(
				"Google returned OAuth error: %s",
				r.URL.Query().Get("error"),
			)
			http.Error(w, "Authorization failed.", http.StatusBadRequest)
			return
		}

		code := r.URL.Query().Get("code")
		if code == "" {
			errResult <- fmt.Errorf("authorization code missing")
			http.Error(w, "Authorization code missing.", http.StatusBadRequest)
			return
		}

		token, err := config.Exchange(ctx, code)
		if err != nil {
			errResult <- fmt.Errorf("failed to exchange authorization code: %w", err)
			http.Error(w, "Authorization failed.", http.StatusInternalServerError)
			return
		}

		fmt.Fprintln(w, "Authorization successful. You can close this window.")

		result <- token
	})

	server := &http.Server{
		Handler: mux,
	}

	go func() {
		if err := server.Serve(listener); err != nil && err != http.ErrServerClosed {
			errResult <- err
		}
	}()

	select {
	case token := <-result:
		_ = server.Shutdown(context.Background())
		return token, nil

	case err := <-errResult:
		_ = server.Shutdown(context.Background())
		return nil, err
	}
}

func saveToken(path string, token *oauth2.Token) error {
	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return err
	}

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	return json.NewEncoder(file).Encode(token)
}

