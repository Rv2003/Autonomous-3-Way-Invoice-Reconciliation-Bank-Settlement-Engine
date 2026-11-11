package email

import (
	"context"
	"fmt"
	"strconv"

	"github.com/Rv2003/Autonomous-3-Way-Invoice-Reconciliation-Bank-Settlement-Engine/BACKEND/internal/gmail"
	"github.com/jackc/pgx/v5/pgtype"
)

type Processor struct {
	service     Service
	gmailClient *gmail.Client
}

func NewProcessor(
	service Service,
	gmailClient *gmail.Client,
) *Processor {
	return &Processor{
		service:     service,
		gmailClient: gmailClient,
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

		// 5b. Extract the plain-text email body.
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

		// 5c. Find attachments in the Gmail message.
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

		// Temporary checkpoint.
		//
		// Next stage:
		// Download PDF attachments and connect them
		// to the existing PDF storage/processing pipeline.
		_ = body
		_ = attachments
	}

	// 6. All discovered messages were processed successfully.
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
