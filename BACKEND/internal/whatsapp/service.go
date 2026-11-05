package whatsapp

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

// DocStore is the storage abstraction for recorded inbound documents.
type DocStore interface {
	Exists(ctx context.Context, source, externalMessageID string) (bool, error)
	Insert(ctx context.Context, d InboundDocument) (bool, error)
	ListRecent(ctx context.Context, source string, limit int) ([]InboundDocument, error)
}

// PgStore implements DocStore on PostgreSQL using pgxpool.
type PgStore struct {
	DB *pgxpool.Pool
}

// NewPgStore returns a new PostgreSQL store.
func NewPgStore(db *pgxpool.Pool) *PgStore {
	return &PgStore{DB: db}
}

// Exists checks if an inbound document with the given source and external ID was already recorded.
func (s *PgStore) Exists(ctx context.Context, source, externalMessageID string) (bool, error) {
	var exists bool
	err := s.DB.QueryRow(ctx,
		`SELECT EXISTS (
			SELECT 1 FROM inbound_documents
			WHERE source = $1 AND external_message_id = $2
		)`,
		source, externalMessageID,
	).Scan(&exists)
	return exists, err
}

// Insert saves metadata to inbound_documents. Returns true if inserted, false on conflict (duplicate).
func (s *PgStore) Insert(ctx context.Context, d InboundDocument) (bool, error) {
	tag, err := s.DB.Exec(ctx, `
		INSERT INTO inbound_documents
		  (source, external_message_id, sender, media_id, mime_type,
		   original_filename, storage_path, sha256, size_bytes)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		ON CONFLICT (source, external_message_id) DO NOTHING`,
		d.Source, d.ExternalMessageID, d.Sender, d.MediaID, d.MimeType,
		d.OriginalFilename, d.StoragePath, d.SHA256, d.SizeBytes,
	)
	if err != nil {
		return false, err
	}
	return tag.RowsAffected() == 1, nil
}

// ListRecent returns the most recently received documents from the given source.
func (s *PgStore) ListRecent(ctx context.Context, source string, limit int) ([]InboundDocument, error) {
	if limit <= 0 {
		limit = 10
	}
	rows, err := s.DB.Query(ctx, `
		SELECT id, source, external_message_id, sender, media_id, mime_type,
		       original_filename, storage_path, sha256, size_bytes, status, received_at
		FROM inbound_documents
		WHERE source = $1
		ORDER BY received_at DESC
		LIMIT $2`,
		source, limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var docs []InboundDocument
	for rows.Next() {
		var d InboundDocument
		err := rows.Scan(
			&d.ID,
			&d.Source,
			&d.ExternalMessageID,
			&d.Sender,
			&d.MediaID,
			&d.MimeType,
			&d.OriginalFilename,
			&d.StoragePath,
			&d.SHA256,
			&d.SizeBytes,
			&d.Status,
			&d.ReceivedAt,
		)
		if err != nil {
			return nil, err
		}
		docs = append(docs, d)
	}
	if docs == nil {
		docs = []InboundDocument{}
	}
	return docs, rows.Err()
}
