package gmail

import (
    "context"
    "fmt"

    gmailapi "google.golang.org/api/gmail/v1"
)

type HistoryResult struct {
    HistoryID  uint64
    MessageIDs []string
}

func (c *Client) ListNewMessages(
    ctx context.Context,
    startHistoryID uint64,
) (*HistoryResult, error) {
    if startHistoryID == 0 {
        return nil, fmt.Errorf("start history ID must be greater than zero")
    }

    result := &HistoryResult{}

    call := c.service.Users.History.
        List("me").
        StartHistoryId(startHistoryID).
        HistoryTypes("messageAdded")

    err := call.Pages(ctx, func(response *gmailapi.ListHistoryResponse) error {
        if response.HistoryId > result.HistoryID {
            result.HistoryID = response.HistoryId
        }

        for _, history := range response.History {
            for _, added := range history.MessagesAdded {
                if added.Message != nil && added.Message.Id != "" {
                    result.MessageIDs = append(
                        result.MessageIDs,
                        added.Message.Id,
                    )
                }
            }
        }

        return nil
    })

    if err != nil {
        return nil, fmt.Errorf("failed to list Gmail history: %w", err)
    }

    return result, nil
}
