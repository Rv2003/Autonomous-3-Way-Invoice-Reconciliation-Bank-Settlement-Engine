package pdf

import (
	"context"

	repo "github.com/Rv2003/Autonomous-3-Way-Invoice-Reconciliation-Bank-Settlement-Engine/BACKEND/internal/adapters/postgresql/sqlc"
)

type Service interface {
	ListPdf(ctx context.Context) ([]repo.ListPdfRow, error)
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
