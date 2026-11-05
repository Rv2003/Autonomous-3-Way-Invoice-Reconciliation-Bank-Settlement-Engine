-- +goose Up
-- +goose StatementBegin

CREATE TABLE IF NOT EXISTS email_mailboxes (
    id               UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    provider         TEXT NOT NULL CHECK (provider IN ('gmail')),
    email_address    TEXT NOT NULL,
    history_id       TEXT,
    watch_expires_at TIMESTAMPTZ,
    status           TEXT NOT NULL DEFAULT 'active'
                     CHECK (status IN ('active', 'disabled', 'reauthorization_required')),
    created_at       TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at       TIMESTAMPTZ NOT NULL DEFAULT now(),

    CONSTRAINT uq_email_mailboxes_provider_email
        UNIQUE (provider, email_address)
);

CREATE INDEX IF NOT EXISTS idx_email_mailboxes_status
    ON email_mailboxes (status);


CREATE TABLE IF NOT EXISTS email_ingestion_events (
    id                UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    mailbox_id        UUID NOT NULL REFERENCES email_mailboxes(id) ON DELETE CASCADE,
    pubsub_message_id TEXT NOT NULL,
    history_id        TEXT NOT NULL,
    status            TEXT NOT NULL DEFAULT 'received'
                      CHECK (status IN ('received', 'processing', 'processed', 'failed')),
    attempts          INTEGER NOT NULL DEFAULT 0
                      CHECK (attempts >= 0),
    last_error        TEXT,
    received_at       TIMESTAMPTZ NOT NULL DEFAULT now(),
    processed_at      TIMESTAMPTZ,

    CONSTRAINT uq_email_ingestion_pubsub_message
        UNIQUE (mailbox_id, pubsub_message_id)
);

CREATE INDEX IF NOT EXISTS idx_email_ingestion_mailbox_history
    ON email_ingestion_events (mailbox_id, history_id);

CREATE INDEX IF NOT EXISTS idx_email_ingestion_status
    ON email_ingestion_events (status);

CREATE INDEX IF NOT EXISTS idx_email_ingestion_received_at
    ON email_ingestion_events (received_at DESC);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS email_ingestion_events;
DROP TABLE IF EXISTS email_mailboxes;

-- +goose StatementEnd
