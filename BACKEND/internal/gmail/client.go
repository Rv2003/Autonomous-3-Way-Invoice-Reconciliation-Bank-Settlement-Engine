package gmail

import (
    "context"
    "encoding/json"
    "fmt"
    "os"

    "golang.org/x/oauth2"
    "golang.org/x/oauth2/google"
    gmailapi "google.golang.org/api/gmail/v1"
    "google.golang.org/api/option"
)

const (
    credentialsFile = "credentials/client_secret_841501621403-f8qd9akhvh1f68le2p953rq0rtjmrv60.apps.googleusercontent.com.json"
    tokenFile       = "tokens/gmail-token.json"
)

type Client struct {
    service *gmailapi.Service
}

func NewClient(ctx context.Context) (*Client, error) {
    credentials, err := os.ReadFile(credentialsFile)
    if err != nil {
        return nil, fmt.Errorf("failed to read OAuth credentials: %w", err)
    }

    config, err := google.ConfigFromJSON(
        credentials,
        gmailapi.GmailReadonlyScope,
    )
    if err != nil {
        return nil, fmt.Errorf("failed to parse OAuth credentials: %w", err)
    }

    tokenData, err := os.ReadFile(tokenFile)
    if err != nil {
        return nil, fmt.Errorf("failed to read Gmail token: %w", err)
    }

    var token oauth2.Token
    if err := json.Unmarshal(tokenData, &token); err != nil {
        return nil, fmt.Errorf("failed to parse Gmail token: %w", err)
    }

    httpClient := config.Client(ctx, &token)

    service, err := gmailapi.NewService(
        ctx,
        option.WithHTTPClient(httpClient),
    )
    if err != nil {
        return nil, fmt.Errorf("failed to create Gmail service: %w", err)
    }

    return &Client{
        service: service,
    }, nil
}

func (c *Client) Service() *gmailapi.Service {
    return c.service
}
