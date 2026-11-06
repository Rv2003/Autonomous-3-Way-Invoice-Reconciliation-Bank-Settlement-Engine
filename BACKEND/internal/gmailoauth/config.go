package gmailoauth

import (
	"fmt"

	"github.com/Rv2003/Autonomous-3-Way-Invoice-Reconciliation-Bank-Settlement-Engine/BACKEND/internal/env"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

func Config(scope string) (*oauth2.Config, error) {
	clientID := env.GetString("GMAIL_CLIENT_ID", "")
	clientSecret := env.GetString("GMAIL_CLIENT_SECRET", "")
	if clientID == "" || clientSecret == "" {
		return nil, fmt.Errorf("GMAIL_CLIENT_ID and GMAIL_CLIENT_SECRET must be set")
	}

	return &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  env.GetString("GMAIL_REDIRECT_URL", "http://localhost:8080/email/gmail/callback"),
		Scopes:       []string{scope},
		Endpoint:     google.Endpoint,
	}, nil
}
