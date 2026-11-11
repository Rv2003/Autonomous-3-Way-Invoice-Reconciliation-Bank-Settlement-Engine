package pdf

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	repo "github.com/Rv2003/Autonomous-3-Way-Invoice-Reconciliation-Bank-Settlement-Engine/BACKEND/internal/adapters/postgresql/sqlc"
	"github.com/jackc/pgx/v5/pgtype"
)

var ErrInvalidDocumentID = errors.New("invalid document ID")

type Service interface {
	ListPdf(ctx context.Context) ([]repo.ListPdfRow, error)
	GetPdf(ctx context.Context, id string) (repo.Pdf, error)
	CreatePdf(ctx context.Context, payload CreatePdfPayload) (repo.Pdf, error)
	CreatePdfFromBytes(
		ctx context.Context,
		filename string,
		data []byte,
	) (repo.Pdf, error)
	DeletePdf(ctx context.Context, id string) (string, error)
}

type svc struct {
	repo repo.Querier
}

func NewService(repo repo.Querier) Service {
	return &svc{repo: repo}
}

func (s *svc) ListPdf(ctx context.Context) ([]repo.ListPdfRow, error) {
	return s.repo.ListPdf(ctx)
}

func (s *svc) GetPdf(ctx context.Context, id string) (repo.Pdf, error) {
	var pdfID pgtype.UUID

	if err := pdfID.Scan(id); err != nil {
		return repo.Pdf{}, ErrInvalidDocumentID
	}

	return s.repo.GetPdf(ctx, pdfID)
}

func (s *svc) CreatePdf(
	ctx context.Context,
	payload CreatePdfPayload,
) (repo.Pdf, error) {
	return s.repo.CreatePdf(ctx, payload.toRepoParams())
}

func (s *svc) CreatePdfFromBytes(
	ctx context.Context,
	filename string,
	data []byte,
) (repo.Pdf, error) {
	if len(data) == 0 {
		return repo.Pdf{}, fmt.Errorf("PDF data is empty")
	}

	if filepath.Ext(filename) == "" {
		filename += ".pdf"
	}

	if err := os.MkdirAll(pdfUploadDir, 0750); err != nil {
		return repo.Pdf{}, fmt.Errorf(
			"create PDF upload directory: %w",
			err,
		)
	}

	hash := sha256.Sum256(data)
	sha256Hex := hex.EncodeToString(hash[:])

	dst, err := os.CreateTemp(
		pdfUploadDir,
		"gmail-pdf-*.pdf",
	)
	if err != nil {
		return repo.Pdf{}, fmt.Errorf(
			"create PDF storage file: %w",
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

		return repo.Pdf{}, fmt.Errorf(
			"write PDF data: %w",
			err,
		)
	}

	if err := dst.Close(); err != nil {
		cleanup()

		return repo.Pdf{}, fmt.Errorf(
			"close PDF storage file: %w",
			err,
		)
	}

	pdf, err := s.CreatePdf(
		ctx,
		CreatePdfPayload{
			OriginalName: filepath.Base(filename),
			StoredPath:   filepath.ToSlash(storedPath),
			SizeBytes:    int64(len(data)),
			Sha256:       sha256Hex,
		},
	)
	if err != nil {
		cleanup()

		return repo.Pdf{}, fmt.Errorf(
			"save PDF metadata: %w",
			err,
		)
	}

	return pdf, nil
}

func (s *svc) DeletePdf(
	ctx context.Context,
	id string,
) (string, error) {
	var pdfID pgtype.UUID

	if err := pdfID.Scan(id); err != nil {
		return "", ErrInvalidDocumentID
	}

	return s.repo.DeletePdf(ctx, pdfID)
}
