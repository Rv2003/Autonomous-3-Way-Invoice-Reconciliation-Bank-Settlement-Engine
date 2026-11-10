package gmail

import (
	"context"
	"encoding/base64"
	"fmt"

	gmailapi "google.golang.org/api/gmail/v1"
)

type AttachmentInfo struct {
	Filename     string
	MimeType     string
	AttachmentID string
	Size         int64
}

func ListAttachments(message *gmailapi.Message) ([]AttachmentInfo, error) {
	if message == nil {
		return nil, fmt.Errorf("message must not be nil")
	}

	var attachments []AttachmentInfo

	findAttachments(message.Payload, &attachments)

	return attachments, nil
}

func findAttachments(
	part *gmailapi.MessagePart,
	attachments *[]AttachmentInfo,
) {
	if part == nil {
		return
	}

	if part.Filename != "" && part.Body != nil {
		*attachments = append(*attachments, AttachmentInfo{
			Filename:     part.Filename,
			MimeType:     part.MimeType,
			AttachmentID: part.Body.AttachmentId,
			Size:         part.Body.Size,
		})
	}

	for _, child := range part.Parts {
		findAttachments(child, attachments)
	}
}

func (c *Client) DownloadAttachment(
	ctx context.Context,
	messageID string,
	attachmentID string,
) ([]byte, error) {
	if messageID == "" {
		return nil, fmt.Errorf("message ID must not be empty")
	}

	if attachmentID == "" {
		return nil, fmt.Errorf("attachment ID must not be empty")
	}

	attachment, err := c.service.Users.Messages.Attachments.
		Get("me", messageID, attachmentID).
		Context(ctx).
		Do()
	if err != nil {
		return nil, fmt.Errorf("failed to download Gmail attachment: %w", err)
	}

	if attachment.Data == "" {
		return nil, fmt.Errorf("Gmail attachment contains no data")
	}

	data, err := base64.URLEncoding.DecodeString(attachment.Data)
	if err != nil {
		return nil, fmt.Errorf("failed to decode Gmail attachment data: %w", err)
	}

	return data, nil
}
