package gmail

import (
	"context"
	"fmt"
	"strings"

	gmailapi "google.golang.org/api/gmail/v1"
)

func (c *Client) GetMessage(
	ctx context.Context,
	messageID string,
) (*gmailapi.Message, error) {
	if messageID == "" {
		return nil, fmt.Errorf("message ID must not be empty")
	}

	message, err := c.service.Users.Messages.
		Get("me", messageID).
		Format("full").
		Do()

	if err != nil {
		return nil, fmt.Errorf("failed to get Gmail message: %w", err)
	}

	return message, nil
}

func ExtractSender(message *gmailapi.Message) (string, error) {
	if message == nil {
		return "", fmt.Errorf("message must not be nil")
	}

	if message.Payload == nil {
		return "", fmt.Errorf("message payload must not be nil")
	}

	for _, header := range message.Payload.Headers {
		if strings.EqualFold(
			strings.TrimSpace(header.Name),
			"From",
		) {
			sender := strings.TrimSpace(header.Value)

			if sender == "" {
				return "", fmt.Errorf("Gmail From header is empty")
			}

			return sender, nil
		}
	}

	return "", fmt.Errorf("Gmail From header not found")
}