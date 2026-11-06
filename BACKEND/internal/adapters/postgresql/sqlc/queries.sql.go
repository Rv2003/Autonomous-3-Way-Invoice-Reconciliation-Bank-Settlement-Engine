package repo

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const createEmailIngestionEvent = `
INSERT INTO email_ingestion_events (
    mailbox_id,
    pubsub_message_id,
    history_id
)
VALUES (
    $1,
    $2,
    $3
)
RETURNING id, mailbox_id, pubsub_message_id, history_id, status, attempts, last_error, received_at, processed_at
`

type CreateEmailIngestionEventParams struct {
	MailboxID       pgtype.UUID `json:"mailbox_id"`
	PubsubMessageID string      `json:"pubsub_message_id"`
	HistoryID       string      `json:"history_id"`
}

func (q *Queries) CreateEmailIngestionEvent(ctx context.Context, arg CreateEmailIngestionEventParams) (EmailIngestionEvent, error) {
	row := q.db.QueryRow(ctx, createEmailIngestionEvent, arg.MailboxID, arg.PubsubMessageID, arg.HistoryID)
	var i EmailIngestionEvent
	err := row.Scan(
		&i.ID,
		&i.MailboxID,
		&i.PubsubMessageID,
		&i.HistoryID,
		&i.Status,
		&i.Attempts,
		&i.LastError,
		&i.ReceivedAt,
		&i.ProcessedAt,
	)
	return i, err
}

const createEmailMailbox = `
INSERT INTO email_mailboxes (
    provider,
    email_address
)
VALUES (
    $1,
    $2
)
RETURNING id, provider, email_address, history_id, watch_expires_at, status, created_at, updated_at
`

type CreateEmailMailboxParams struct {
	Provider     string `json:"provider"`
	EmailAddress string `json:"email_address"`
}

func (q *Queries) CreateEmailMailbox(ctx context.Context, arg CreateEmailMailboxParams) (EmailMailbox, error) {
	row := q.db.QueryRow(ctx, createEmailMailbox, arg.Provider, arg.EmailAddress)
	var i EmailMailbox
	err := row.Scan(
		&i.ID,
		&i.Provider,
		&i.EmailAddress,
		&i.HistoryID,
		&i.WatchExpiresAt,
		&i.Status,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const createPdf = `
INSERT INTO pdfs (
    original_name,
    stored_path,
    size_bytes,
    sha256
)
VALUES (
    $1,
    $2,
    $3,
    $4
)
RETURNING id, original_name, stored_path, size_bytes, sha256, status, doc_type, extracted_text, extracted_data, error_message, uploaded_at, processed_at
`

type CreatePdfParams struct {
	OriginalName string `json:"original_name"`
	StoredPath   string `json:"stored_path"`
	SizeBytes    int64  `json:"size_bytes"`
	Sha256       string `json:"sha256"`
}

func (q *Queries) CreatePdf(ctx context.Context, arg CreatePdfParams) (Pdf, error) {
	row := q.db.QueryRow(ctx, createPdf,
		arg.OriginalName,
		arg.StoredPath,
		arg.SizeBytes,
		arg.Sha256,
	)
	var i Pdf
	err := row.Scan(
		&i.ID,
		&i.OriginalName,
		&i.StoredPath,
		&i.SizeBytes,
		&i.Sha256,
		&i.Status,
		&i.DocType,
		&i.ExtractedText,
		&i.ExtractedData,
		&i.ErrorMessage,
		&i.UploadedAt,
		&i.ProcessedAt,
	)
	return i, err
}

const getEmailMailboxByEmail = `
SELECT id, provider, email_address, history_id, watch_expires_at, status, created_at, updated_at
FROM email_mailboxes
WHERE provider = $1
  AND email_address = $2
`

type GetEmailMailboxByEmailParams struct {
	Provider     string `json:"provider"`
	EmailAddress string `json:"email_address"`
}

func (q *Queries) GetEmailMailboxByEmail(ctx context.Context, arg GetEmailMailboxByEmailParams) (EmailMailbox, error) {
	row := q.db.QueryRow(ctx, getEmailMailboxByEmail, arg.Provider, arg.EmailAddress)
	var i EmailMailbox
	err := row.Scan(
		&i.ID,
		&i.Provider,
		&i.EmailAddress,
		&i.HistoryID,
		&i.WatchExpiresAt,
		&i.Status,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const deletePdf = `
DELETE FROM pdfs
WHERE id::text = $1
RETURNING stored_path
`

func (q *Queries) DeletePdf(ctx context.Context, id string) (string, error) {
	var storedPath string
	err := q.db.QueryRow(ctx, deletePdf, id).Scan(&storedPath)
	return storedPath, err
}

const getPdf = `
SELECT id, original_name, stored_path, size_bytes, sha256, status, doc_type, extracted_text, extracted_data, error_message, uploaded_at, processed_at
FROM pdfs
WHERE id = $1
`

func (q *Queries) GetPdf(ctx context.Context, id pgtype.UUID) (Pdf, error) {
	row := q.db.QueryRow(ctx, getPdf, id)
	var i Pdf
	err := row.Scan(
		&i.ID,
		&i.OriginalName,
		&i.StoredPath,
		&i.SizeBytes,
		&i.Sha256,
		&i.Status,
		&i.DocType,
		&i.ExtractedText,
		&i.ExtractedData,
		&i.ErrorMessage,
		&i.UploadedAt,
		&i.ProcessedAt,
	)
	return i, err
}

const listEmailIngestionEvents = `
SELECT id, mailbox_id, pubsub_message_id, history_id, status, attempts, last_error, received_at, processed_at
FROM email_ingestion_events
ORDER BY received_at DESC
`

func (q *Queries) ListEmailIngestionEvents(ctx context.Context) ([]EmailIngestionEvent, error) {
	rows, err := q.db.Query(ctx, listEmailIngestionEvents)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []EmailIngestionEvent
	for rows.Next() {
		var i EmailIngestionEvent
		if err := rows.Scan(
			&i.ID,
			&i.MailboxID,
			&i.PubsubMessageID,
			&i.HistoryID,
			&i.Status,
			&i.Attempts,
			&i.LastError,
			&i.ReceivedAt,
			&i.ProcessedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listPdf = `
SELECT id, original_name, size_bytes, status, uploaded_at
FROM pdfs
ORDER BY uploaded_at DESC
`

type ListPdfRow struct {
	ID           pgtype.UUID        `json:"id"`
	OriginalName string             `json:"original_name"`
	SizeBytes    int64              `json:"size_bytes"`
	Status       string             `json:"status"`
	UploadedAt   pgtype.Timestamptz `json:"uploaded_at"`
}

func (q *Queries) ListPdf(ctx context.Context) ([]ListPdfRow, error) {
	rows, err := q.db.Query(ctx, listPdf)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ListPdfRow
	for rows.Next() {
		var i ListPdfRow
		if err := rows.Scan(
			&i.ID,
			&i.OriginalName,
			&i.SizeBytes,
			&i.Status,
			&i.UploadedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const markEmailIngestionEventFailed = `
UPDATE email_ingestion_events
SET
    status = 'failed',
    last_error = $2
WHERE id = $1
  AND status = 'processing'
RETURNING id, mailbox_id, pubsub_message_id, history_id, status, attempts, last_error, received_at, processed_at
`

type MarkEmailIngestionEventFailedParams struct {
	ID        pgtype.UUID `json:"id"`
	LastError pgtype.Text `json:"last_error"`
}

func (q *Queries) MarkEmailIngestionEventFailed(ctx context.Context, arg MarkEmailIngestionEventFailedParams) (EmailIngestionEvent, error) {
	row := q.db.QueryRow(ctx, markEmailIngestionEventFailed, arg.ID, arg.LastError)
	var i EmailIngestionEvent
	err := row.Scan(
		&i.ID,
		&i.MailboxID,
		&i.PubsubMessageID,
		&i.HistoryID,
		&i.Status,
		&i.Attempts,
		&i.LastError,
		&i.ReceivedAt,
		&i.ProcessedAt,
	)
	return i, err
}

const markEmailIngestionEventProcessed = `
UPDATE email_ingestion_events
SET
    status = 'processed',
    processed_at = now()
WHERE id = $1
  AND status = 'processing'
RETURNING id, mailbox_id, pubsub_message_id, history_id, status, attempts, last_error, received_at, processed_at
`

func (q *Queries) MarkEmailIngestionEventProcessed(ctx context.Context, id pgtype.UUID) (EmailIngestionEvent, error) {
	row := q.db.QueryRow(ctx, markEmailIngestionEventProcessed, id)
	var i EmailIngestionEvent
	err := row.Scan(
		&i.ID,
		&i.MailboxID,
		&i.PubsubMessageID,
		&i.HistoryID,
		&i.Status,
		&i.Attempts,
		&i.LastError,
		&i.ReceivedAt,
		&i.ProcessedAt,
	)
	return i, err
}

const markEmailIngestionEventProcessing = `
UPDATE email_ingestion_events
SET
    status = 'processing',
    attempts = attempts + 1
WHERE id = $1
  AND status = 'received'
RETURNING id, mailbox_id, pubsub_message_id, history_id, status, attempts, last_error, received_at, processed_at
`

func (q *Queries) MarkEmailIngestionEventProcessing(ctx context.Context, id pgtype.UUID) (EmailIngestionEvent, error) {
	row := q.db.QueryRow(ctx, markEmailIngestionEventProcessing, id)
	var i EmailIngestionEvent
	err := row.Scan(
		&i.ID,
		&i.MailboxID,
		&i.PubsubMessageID,
		&i.HistoryID,
		&i.Status,
		&i.Attempts,
		&i.LastError,
		&i.ReceivedAt,
		&i.ProcessedAt,
	)
	return i, err
}

const updatePdfStatus = `
UPDATE pdfs
SET
    status = $2,
    processed_at = now()
WHERE id = $1
`

type UpdatePdfStatusParams struct {
	ID     pgtype.UUID `json:"id"`
	Status string      `json:"status"`
}

func (q *Queries) UpdatePdfStatus(ctx context.Context, arg UpdatePdfStatusParams) error {
	_, err := q.db.Exec(ctx, updatePdfStatus, arg.ID, arg.Status)
	return err
}
