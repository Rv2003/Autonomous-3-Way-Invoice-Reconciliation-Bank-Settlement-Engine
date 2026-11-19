package email

import (
	"context"
	"fmt"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/Rv2003/Autonomous-3-Way-Invoice-Reconciliation-Bank-Settlement-Engine/BACKEND/internal/gmail"
	"github.com/Rv2003/Autonomous-3-Way-Invoice-Reconciliation-Bank-Settlement-Engine/BACKEND/internal/inbound"
	"github.com/Rv2003/Autonomous-3-Way-Invoice-Reconciliation-Bank-Settlement-Engine/BACKEND/internal/pdf"
	"github.com/jackc/pgx/v5/pgtype"
)

type Processor struct {
	service        Service
	gmailClient    *gmail.Client
	pdfService     pdf.Service
	inboundService inbound.Service
}

func NewProcessor(
	service Service,
	gmailClient *gmail.Client,
	pdfService pdf.Service,
	inboundService inbound.Service,
) *Processor {
	return &Processor{
		service:        service,
		gmailClient:    gmailClient,
		pdfService:     pdfService,
		inboundService: inboundService,
	}
}

func (p *Processor) Process(
	ctx context.Context,
	eventID pgtype.UUID,
) error {
	// 1. Load the ingestion event.
	event, err := p.service.GetIngestionEvent(ctx, eventID)
	if err != nil {
		return fmt.Errorf("get ingestion event: %w", err)
	}

	// 2. Mark the event as processing.
	_, err = p.service.MarkIngestionEventProcessing(ctx, eventID)
	if err != nil {
		return fmt.Errorf(
			"mark ingestion event processing: %w",
			err,
		)
	}

	// 3. Convert Gmail history ID from the database string
	//    into the uint64 required by the Gmail History API.
	historyID, err := strconv.ParseUint(
		event.HistoryID,
		10,
		64,
	)
	if err != nil {
		_, _ = p.service.MarkIngestionEventFailed(
			ctx,
			eventID,
			fmt.Sprintf("invalid history ID: %v", err),
		)

		return fmt.Errorf(
			"parse history ID %q: %w",
			event.HistoryID,
			err,
		)
	}

	// 4. Ask Gmail for messages added since this history ID.
	history, err := p.gmailClient.ListNewMessages(
		ctx,
		historyID,
	)
	if err != nil {
		_, _ = p.service.MarkIngestionEventFailed(
			ctx,
			eventID,
			err.Error(),
		)

		return fmt.Errorf(
			"list new Gmail messages: %w",
			err,
		)
	}

	// 5. Process every new Gmail message.
	for _, messageID := range history.MessageIDs {
		// Skip if already processed or deleted
		if existing, err := p.inboundService.GetDocumentByExternalID(ctx, "email", messageID); err == nil && existing.ID != 0 {
			continue
		}

		// 5a. Retrieve the complete Gmail message.
		message, err := p.gmailClient.GetMessage(
			ctx,
			messageID,
		)
		if err != nil {
			_, _ = p.service.MarkIngestionEventFailed(
				ctx,
				eventID,
				err.Error(),
			)

			return fmt.Errorf(
				"get Gmail message %q: %w",
				messageID,
				err,
			)
		}

		// 5b. Extract the sender.
		sender, err := gmail.ExtractSender(message)
		if err != nil {
			_, _ = p.service.MarkIngestionEventFailed(
				ctx,
				eventID,
				err.Error(),
			)

			return fmt.Errorf(
				"extract Gmail sender %q: %w",
				messageID,
				err,
			)
		}

		// 5c. Extract the plain-text email body.
		body, err := gmail.ExtractBody(message)
		if err != nil {
			_, _ = p.service.MarkIngestionEventFailed(
				ctx,
				eventID,
				err.Error(),
			)

			return fmt.Errorf(
				"extract Gmail message body %q: %w",
				messageID,
				err,
			)
		}

		// 5d. Find attachments in the Gmail message.
		attachments, err := gmail.ListAttachments(message)
		if err != nil {
			_, _ = p.service.MarkIngestionEventFailed(
				ctx,
				eventID,
				err.Error(),
			)

			return fmt.Errorf(
				"list Gmail attachments %q: %w",
				messageID,
				err,
			)
		}

		// The email body is currently extracted successfully.
		// It will be used later when email metadata/content
		// is persisted into the reconciliation workflow.
		_ = body

		// 5e. Process PDF attachments.
		for _, attachment := range attachments {
			if !isPDFAttachment(
				attachment.Filename,
				attachment.MimeType,
			) {
				continue
			}

			// Download the actual attachment bytes from Gmail.
			data, err := p.gmailClient.DownloadAttachment(
				ctx,
				message.Id,
				attachment.AttachmentID,
			)
			if err != nil {
				_, _ = p.service.MarkIngestionEventFailed(
					ctx,
					eventID,
					err.Error(),
				)

				return fmt.Errorf(
					"download PDF attachment %q from Gmail message %q: %w",
					attachment.Filename,
					messageID,
					err,
				)
			}

			if len(data) == 0 {
				err := fmt.Errorf(
					"PDF attachment %q is empty",
					attachment.Filename,
				)

				_, _ = p.service.MarkIngestionEventFailed(
					ctx,
					eventID,
					err.Error(),
				)

				return err
			}

			// Store generic inbound-document metadata and file.
			_, err = p.inboundService.CreateDocument(
				ctx,
				inbound.CreateInboundDocumentPayload{
					Source:            "email",
					ExternalMessageID: message.Id,
					Sender:            sender,
					MediaID:           attachment.AttachmentID,
					MimeType:          attachment.MimeType,
					OriginalFilename:  attachment.Filename,
					Status:            "RECEIVED",
				},
				data,
			)
			if err != nil {
				_, _ = p.service.MarkIngestionEventFailed(
					ctx,
					eventID,
					err.Error(),
				)

				return fmt.Errorf(
					"store inbound document %q: %w",
					attachment.Filename,
					err,
				)
			}

			// Preserve the existing PDF storage path.
			_, err = p.pdfService.CreatePdfFromBytes(
				ctx,
				attachment.Filename,
				data,
			)
			if err != nil {
				_, _ = p.service.MarkIngestionEventFailed(
					ctx,
					eventID,
					err.Error(),
				)

				return fmt.Errorf(
					"store PDF attachment %q: %w",
					attachment.Filename,
					err,
				)
			}
		}
	}

	// 6. All discovered messages and PDF attachments
	//    were processed successfully.
	_, err = p.service.MarkIngestionEventProcessed(
		ctx,
		eventID,
	)
	if err != nil {
		return fmt.Errorf(
			"mark ingestion event processed: %w",
			err,
		)
	}

	return nil
}

func (p *Processor) SyncRecent(
	ctx context.Context,
	maxMessages int,
) (int, error) {
	if maxMessages <= 0 {
		maxMessages = 10
	}

	messages, err := p.gmailClient.Service().Users.Messages.
		List("me").
		MaxResults(int64(maxMessages)).
		Do()
	if err != nil {
		return 0, fmt.Errorf("list gmail messages: %w", err)
	}

	count := 0
	for _, msgRef := range messages.Messages {
		// Check if this message was already ingested or deleted
		if existing, err := p.inboundService.GetDocumentByExternalID(ctx, "email", msgRef.Id); err == nil && existing.ID != 0 {
			continue
		}

		message, err := p.gmailClient.GetMessage(ctx, msgRef.Id)
		if err != nil {
			continue
		}

		sender, _ := gmail.ExtractSender(message)
		attachments, err := gmail.ListAttachments(message)
		if err != nil {
			continue
		}

		for _, attachment := range attachments {
			if !isPDFAttachment(attachment.Filename, attachment.MimeType) {
				continue
			}

			data, err := p.gmailClient.DownloadAttachment(
				ctx,
				message.Id,
				attachment.AttachmentID,
			)
			if err != nil || len(data) == 0 {
				continue
			}

			_, err = p.inboundService.CreateDocument(
				ctx,
				inbound.CreateInboundDocumentPayload{
					Source:            "email",
					ExternalMessageID: message.Id,
					Sender:            sender,
					MediaID:           attachment.AttachmentID,
					MimeType:          attachment.MimeType,
					OriginalFilename:  attachment.Filename,
					Status:            "RECEIVED",
				},
				data,
			)
			if err != nil {
				// Duplicate or error, continue
				continue
			}

			_, _ = p.pdfService.CreatePdfFromBytes(
				ctx,
				attachment.Filename,
				data,
			)

			count++
		}
	}

	return count, nil
}

func isPDFAttachment(
	filename string,
	mimeType string,
) bool {
	if strings.EqualFold(
		strings.TrimSpace(mimeType),
		"application/pdf",
	) {
		return true
	}

	return strings.EqualFold(
		filepath.Ext(filename),
		".pdf",
	)
}