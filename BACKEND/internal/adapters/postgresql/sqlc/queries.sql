-- name: ListPdf :many
SELECT id, original_name, size_bytes, status, uploaded_at
FROM pdfs
ORDER BY uploaded_at DESC;


-- name: GetPdf :one
SELECT *
FROM pdfs
WHERE id = $1;


-- name: CreatePdf :one
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
RETURNING *;

-- name: DeletePdf :one
DELETE FROM pdfs
WHERE id::text = $1
RETURNING stored_path;

-- name: UpdatePdfStatus :exec
UPDATE pdfs
SET
    status = $2,
    processed_at = now()
WHERE id = $1;


-- name: CreateEmailMailbox :one
INSERT INTO email_mailboxes (
    provider,
    email_address
)
VALUES (
    $1,
    $2
)
RETURNING *;


-- name: GetEmailMailboxByEmail :one
SELECT *
FROM email_mailboxes
WHERE provider = $1
  AND email_address = $2;

-- name: UpdateEmailMailboxWatch :one
UPDATE email_mailboxes
SET
    history_id = $2,
    watch_expires_at = $3,
    updated_at = NOW()
WHERE id = $1
RETURNING *;



-- name: CreateEmailIngestionEvent :one
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
RETURNING *;


-- name: ListEmailIngestionEvents :many
SELECT *
FROM email_ingestion_events
ORDER BY received_at DESC;


-- name: MarkEmailIngestionEventProcessing :one
UPDATE email_ingestion_events
SET
    status = 'processing',
    attempts = attempts + 1
WHERE id = $1
  AND status = 'received'
RETURNING *;


-- name: MarkEmailIngestionEventProcessed :one
UPDATE email_ingestion_events
SET
    status = 'processed',
    processed_at = now()
WHERE id = $1
  AND status = 'processing'
RETURNING *;

-- name: MarkEmailIngestionEventFailed :one
UPDATE email_ingestion_events
SET
    status = 'failed',
    last_error = $2
WHERE id = $1
  AND status = 'processing'
RETURNING *;