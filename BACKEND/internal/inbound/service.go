package inbound

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"

	repo "github.com/Rv2003/Autonomous-3-Way-Invoice-Reconciliation-Bank-Settlement-Engine/BACKEND/internal/adapters/postgresql/sqlc"
)

const inboundUploadDir = "uploads/inbound"

type Service interface {
	CreateDocument(
		ctx context.Context,
		payload CreateInboundDocumentPayload,
		data []byte,
	) (repo.InboundDocument, error)
	ListRecent(
		ctx context.Context,
		source string,
		limit int32,
	) ([]repo.InboundDocument, error)
	GetDocument(
		ctx context.Context,
		id int64,
	) (repo.InboundDocument, error)
	GetDocumentByExternalID(
		ctx context.Context,
		source string,
		externalID string,
	) (repo.InboundDocument, error)
	DeleteDocument(
		ctx context.Context,
		id int64,
	) (repo.DeleteInboundDocumentRow, error)
}

type svc struct {
	repo repo.Querier
}

func NewService(repo repo.Querier) Service {
	return &svc{
		repo: repo,
	}
}

func (s *svc) CreateDocument(
	ctx context.Context,
	payload CreateInboundDocumentPayload,
	data []byte,
) (repo.InboundDocument, error) {
	if len(data) == 0 {
		return repo.InboundDocument{}, fmt.Errorf(
			"inbound document data is empty",
		)
	}

	if err := os.MkdirAll(inboundUploadDir, 0750); err != nil {
		return repo.InboundDocument{}, fmt.Errorf(
			"create inbound upload directory: %w",
			err,
		)
	}

	hash := sha256.Sum256(data)
	sha256Hex := hex.EncodeToString(hash[:])

	filename := filepath.Base(payload.OriginalFilename)

	if filename == "." || filename == "" {
		filename = "inbound-document"
	}

	dst, err := os.CreateTemp(
		inboundUploadDir,
		"inbound-*"+filepath.Ext(filename),
	)
	if err != nil {
		return repo.InboundDocument{}, fmt.Errorf(
			"create inbound storage file: %w",
			err,
		)
	}

	storedPath := dst.Name()

	cleanup := func() {
		_ = os.Remove(storedPath)
	}

	if _, err := dst.Write(data); err != nil {
		_ = dst.Close()
		cleanup()

		return repo.InboundDocument{}, fmt.Errorf(
			"write inbound document: %w",
			err,
		)
	}

	if err := dst.Close(); err != nil {
		cleanup()

		return repo.InboundDocument{}, fmt.Errorf(
			"close inbound storage file: %w",
			err,
		)
	}

	payload.OriginalFilename = filename
	payload.StoragePath = filepath.ToSlash(storedPath)
	payload.Sha256 = sha256Hex
	payload.SizeBytes = int64(len(data))

	if payload.Status == "" {
		payload.Status = "RECEIVED"
	}

	document, err := s.repo.CreateInboundDocument(
		ctx,
		payload.toRepoParams(),
	)
	if err != nil {
		cleanup()

		return repo.InboundDocument{}, fmt.Errorf(
			"save inbound document metadata: %w",
			err,
		)
	}

	return document, nil
}

func (s *svc) ListRecent(
	ctx context.Context,
	source string,
	limit int32,
) ([]repo.InboundDocument, error) {
	if limit <= 0 {
		limit = 10
	}

	if source == "" {
		return s.repo.ListRecentInboundDocuments(ctx, limit)
	}

	return s.repo.ListInboundDocumentsBySource(
		ctx,
		repo.ListInboundDocumentsBySourceParams{
			Source: source,
			Limit:  limit,
		},
	)
}

func (s *svc) GetDocument(
	ctx context.Context,
	id int64,
) (repo.InboundDocument, error) {
	return s.repo.GetInboundDocument(ctx, id)
}

func (s *svc) GetDocumentByExternalID(
	ctx context.Context,
	source string,
	externalID string,
) (repo.InboundDocument, error) {
	return s.repo.GetInboundDocumentByExternalID(
		ctx,
		repo.GetInboundDocumentByExternalIDParams{
			Source:            source,
			ExternalMessageID: externalID,
		},
	)
}

func (s *svc) DeleteDocument(
	ctx context.Context,
	id int64,
) (repo.DeleteInboundDocumentRow, error) {
	row, err := s.repo.DeleteInboundDocument(ctx, id)
	if err != nil {
		return row, err
	}

	// Delete matching entry from pdfs table by SHA256 if present
	if row.Sha256 != "" {
		paths, err := s.repo.DeletePdfBySha256(ctx, row.Sha256)
		if err == nil {
			for _, p := range paths {
				_ = os.Remove(p)
			}
		}
	}

	return row, nil
}