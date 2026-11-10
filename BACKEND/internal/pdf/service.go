package pdf

import (
	"context"
	"errors"

	repo "github.com/Rv2003/Autonomous-3-Way-Invoice-Reconciliation-Bank-Settlement-Engine/BACKEND/internal/adapters/postgresql/sqlc"
	"github.com/jackc/pgx/v5/pgtype"
)

var ErrInvalidDocumentID = errors.New("invalid document ID")

type Service interface {
	ListPdf(ctx context.Context) ([]repo.ListPdfRow, error)
	GetPdf(ctx context.Context, id string) (repo.Pdf, error)
	CreatePdf(ctx context.Context, payload CreatePdfPayload) (repo.Pdf, error)
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

func (s *svc) CreatePdf(ctx context.Context, payload CreatePdfPayload) (repo.Pdf, error) {
	return s.repo.CreatePdf(ctx, payload.toRepoParams())
}
func (s *svc) DeletePdf(ctx context.Context, id string) (string, error) {
	var pdfID pgtype.UUID
	if err := pdfID.Scan(id); err != nil {
		return "", ErrInvalidDocumentID
	}

	return s.repo.DeletePdf(ctx, pdfID)
}
