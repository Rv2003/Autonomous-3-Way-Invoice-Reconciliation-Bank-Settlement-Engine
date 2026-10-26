-- +goose Up
-- +goose StatementBegin
CREATE TABLE pdfs (
    id             UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    original_name  TEXT        NOT NULL,
    stored_path    TEXT        NOT NULL,
    size_bytes     BIGINT      NOT NULL CHECK (size_bytes > 0),
    sha256         TEXT        NOT NULL UNIQUE,
    status         TEXT        NOT NULL DEFAULT 'uploaded'
                   CHECK (status IN ('uploaded', 'processing', 'done', 'failed')),
    doc_type       TEXT,
    extracted_text TEXT,
    extracted_data JSONB,
    error_message  TEXT,
    uploaded_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
    processed_at   TIMESTAMPTZ
);

CREATE INDEX idx_pdfs_status      ON pdfs (status);
CREATE INDEX idx_pdfs_uploaded_at ON pdfs (uploaded_at DESC);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS pdfs;
-- +goose StatementEnd