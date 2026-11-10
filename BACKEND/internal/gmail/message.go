package gmail

import (
    "context"
    "fmt"

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
