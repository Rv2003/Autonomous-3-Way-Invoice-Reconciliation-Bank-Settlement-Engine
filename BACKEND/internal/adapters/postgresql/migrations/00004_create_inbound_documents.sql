-- +goose Up
-- Inbound document storage migration.
CREATE TABLE IF NOT EXISTS inbound_documents (
    id                  BIGSERIAL PRIMARY KEY,
    source              TEXT        NOT NULL,           -- 'telegram' | 'whatsapp' | 'email' | 'upload'
    external_message_id TEXT        NOT NULL,           -- Message ID from platform
    sender              TEXT        NOT NULL,           -- Sender username or phone
    media_id            TEXT        NOT NULL,           -- Platform file id
    mime_type           TEXT        NOT NULL,
    original_filename   TEXT        NOT NULL DEFAULT '',
    storage_path        TEXT        NOT NULL,
    sha256              TEXT        NOT NULL,
    size_bytes          BIGINT      NOT NULL,
    status              TEXT        NOT NULL DEFAULT 'RECEIVED',
    received_at         TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (source, external_message_id)                -- Idempotency constraint
);

CREATE INDEX IF NOT EXISTS idx_inbound_documents_received_at ON inbound_documents (received_at DESC);
CREATE INDEX IF NOT EXISTS idx_inbound_documents_sha256 ON inbound_documents (sha256);

-- +goose Down
DROP TABLE IF EXISTS inbound_documents;
