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

-- name: GetEmailIngestionEvent :one
SELECT *
FROM email_ingestion_events
WHERE id = $1;


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

-- name: CreateInboundDocument :one
INSERT INTO inbound_documents (
    source,
    external_message_id,
    sender,
    media_id,
    mime_type,
    original_filename,
    storage_path,
    sha256,
    size_bytes,
    status
)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6,
    $7,
    $8,
    $9,
    $10
)
RETURNING *;

-- name: ListInboundDocumentsBySource :many
SELECT *
FROM inbound_documents
WHERE source = $1 AND status != 'DELETED'
ORDER BY received_at DESC
LIMIT $2;

-- name: GetInboundDocument :one
SELECT *
FROM inbound_documents
WHERE id = $1;

-- name: GetInboundDocumentByExternalID :one
SELECT *
FROM inbound_documents
WHERE source = $1 AND external_message_id = $2
LIMIT 1;

-- name: GetActiveEmailMailbox :one
SELECT *
FROM email_mailboxes
WHERE provider = $1 AND status = 'active'
ORDER BY updated_at DESC
LIMIT 1;

-- name: ListRecentInboundDocuments :many
SELECT *
FROM inbound_documents
WHERE status != 'DELETED'
ORDER BY received_at DESC
LIMIT $1;

-- name: DeleteInboundDocument :one
UPDATE inbound_documents
SET status = 'DELETED'
WHERE id = $1
RETURNING storage_path, sha256;

-- name: DeletePdfBySha256 :many
DELETE FROM pdfs
WHERE sha256 = $1
RETURNING stored_path;