package gmailoauth

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/Rv2003/Autonomous-3-Way-Invoice-Reconciliation-Bank-Settlement-Engine/BACKEND/internal/env"
	"golang.org/x/oauth2"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

const defaultTokenFile = "tokens/gmail-token.json"

type Inbox struct {
	Email    string         `json:"email"`
	Messages []InboxMessage `json:"messages"`
}

type InboxMessage struct {
	ID             string    `json:"id"`
	From           string    `json:"from"`
	Subject        string    `json:"subject"`
	Snippet        string    `json:"snippet"`
	ReceivedAt     time.Time `json:"received_at"`
	HasAttachments bool      `json:"has_attachments"`
}

func NewService(ctx context.Context) (*gmail.Service, error) {
	config, err := Config(gmail.GmailReadonlyScope)
	if err != nil {
		return nil, err
	}

	tokenPath, err := resolveTokenPath()
	if err != nil {
		return nil, err
	}
	tokenData, err := os.ReadFile(tokenPath)
	if err != nil {
		return nil, fmt.Errorf("read Gmail token file %q: %w", tokenPath, err)
	}

	var token oauth2.Token
	if err := json.Unmarshal(tokenData, &token); err != nil {
		return nil, fmt.Errorf("parse Gmail token file: %w", err)
	}

	service, err := gmail.NewService(ctx, option.WithHTTPClient(config.Client(ctx, &token)))
	if err != nil {
		return nil, fmt.Errorf("create Gmail API client: %w", err)
	}

	return service, nil
}

func SaveToken(token *oauth2.Token) error {
	tokenPath, err := resolveTokenPath()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(tokenPath), 0700); err != nil {
		return fmt.Errorf("create Gmail token directory: %w", err)
	}

	data, err := json.Marshal(token)
	if err != nil {
		return fmt.Errorf("encode Gmail token: %w", err)
	}

	if err := os.WriteFile(tokenPath, data, 0600); err != nil {
		return fmt.Errorf("save Gmail token: %w", err)
	}
	return nil
}

func resolveTokenPath() (string, error) {
	if configuredPath := env.GetString("GMAIL_TOKEN_FILE", ""); configuredPath != "" {
		return configuredPath, nil
	}

	candidates := []string{
		defaultTokenFile,
		filepath.Join("..", defaultTokenFile),
		filepath.Join("..", "..", "BACKEND", defaultTokenFile),
	}
	for _, candidate := range candidates {
		if _, err := os.Stat(candidate); err == nil {
			return candidate, nil
		} else if !os.IsNotExist(err) {
			return "", fmt.Errorf("check Gmail token file %q: %w", candidate, err)
		}
	}

	return defaultTokenFile, nil
}

func FetchInbox(ctx context.Context, limit int64) (Inbox, error) {
	service, err := NewService(ctx)
	if err != nil {
		return Inbox{}, err
	}

	profile, err := service.Users.GetProfile("me").Context(ctx).Do()
	if err != nil {
		return Inbox{}, fmt.Errorf("get Gmail profile: %w", err)
	}

	list, err := service.Users.Messages.List("me").
		LabelIds("INBOX").
		MaxResults(limit).
		Context(ctx).
		Do()
	if err != nil {
		return Inbox{}, fmt.Errorf("list Gmail inbox messages: %w", err)
	}

	inbox := Inbox{
		Email:    profile.EmailAddress,
		Messages: make([]InboxMessage, 0, len(list.Messages)),
	}
	for _, listedMessage := range list.Messages {
		message, err := service.Users.Messages.Get("me", listedMessage.Id).
			Format("full").
			MetadataHeaders("From", "Subject").
			Context(ctx).
			Do()
		if err != nil {
			return Inbox{}, fmt.Errorf("get Gmail message %q: %w", listedMessage.Id, err)
		}

		inbox.Messages = append(inbox.Messages, InboxMessage{
			ID:             message.Id,
			From:           headerValue(message.Payload, "From"),
			Subject:        headerValue(message.Payload, "Subject"),
			Snippet:        message.Snippet,
			ReceivedAt:     time.UnixMilli(message.InternalDate).UTC(),
			HasAttachments: hasAttachments(message.Payload),
		})
	}

	return inbox, nil
}

func headerValue(payload *gmail.MessagePart, name string) string {
	if payload == nil {
		return ""
	}
	for _, header := range payload.Headers {
		if strings.EqualFold(header.Name, name) {
			return header.Value
		}
	}
	return ""
}

func hasAttachments(part *gmail.MessagePart) bool {
	if part == nil {
		return false
	}
	if part.Filename != "" || (part.Body != nil && part.Body.AttachmentId != "") {
		return true
	}
	for _, child := range part.Parts {
		if hasAttachments(child) {
			return true
		}
	}
	return false
}
