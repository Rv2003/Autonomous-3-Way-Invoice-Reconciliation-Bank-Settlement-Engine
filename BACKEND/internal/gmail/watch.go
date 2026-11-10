package gmail

import (
	"context"
	"fmt"

	gmailapi "google.golang.org/api/gmail/v1"
)

type WatchResult struct {
	HistoryID  uint64
	Expiration int64
}

func (c *Client) StartWatch(
	ctx context.Context,
	topicName string,
) (*WatchResult, error) {
	if topicName == "" {
		return nil, fmt.Errorf("topic name must not be empty")
	}

	response, err := c.service.Users.Watch("me", &gmailapi.WatchRequest{
		TopicName: topicName,
	}).Do()
	if err != nil {
		return nil, fmt.Errorf("failed to start Gmail watch: %w", err)
	}

	return &WatchResult{
		HistoryID:  response.HistoryId,
		Expiration: response.Expiration,
	}, nil
}
