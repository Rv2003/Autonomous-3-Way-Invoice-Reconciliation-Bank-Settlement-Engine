package email

import (
	"context"

	repo "github.com/Rv2003/Autonomous-3-Way-Invoice-Reconciliation-Bank-Settlement-Engine/BACKEND/internal/adapters/postgresql/sqlc"
	"github.com/jackc/pgx/v5/pgtype"
)

type Service interface {
	CreateMailbox(
		ctx context.Context,
		params repo.CreateEmailMailboxParams,
	) (repo.EmailMailbox, error)

	GetMailboxByEmail(
		ctx context.Context,
		provider string,
		emailAddress string,
	) (repo.EmailMailbox, error)

	UpdateMailboxWatch(
		ctx context.Context,
		id pgtype.UUID,
		historyID string,
		watchExpiresAt pgtype.Timestamptz,
	) (repo.EmailMailbox, error)

	CreateIngestionEvent(
		ctx context.Context,
		params repo.CreateEmailIngestionEventParams,
	) (repo.EmailIngestionEvent, error)

	GetIngestionEvent(
		ctx context.Context,
		id pgtype.UUID,
	) (repo.EmailIngestionEvent, error)

	ListIngestionEvents(
		ctx context.Context,
	) ([]repo.EmailIngestionEvent, error)

	MarkIngestionEventProcessing(
		ctx context.Context,
		id pgtype.UUID,
	) (repo.EmailIngestionEvent, error)

	MarkIngestionEventProcessed(
		ctx context.Context,
		id pgtype.UUID,
	) (repo.EmailIngestionEvent, error)

	MarkIngestionEventFailed(
		ctx context.Context,
		id pgtype.UUID,
		lastError string,
	) (repo.EmailIngestionEvent, error)
}

type svc struct {
	repo repo.Querier
}

func NewService(repo repo.Querier) Service {
	return &svc{
		repo: repo,
	}
}

func (s *svc) CreateMailbox(
	ctx context.Context,
	params repo.CreateEmailMailboxParams,
) (repo.EmailMailbox, error) {
	return s.repo.CreateEmailMailbox(ctx, params)
}

func (s *svc) GetMailboxByEmail(
	ctx context.Context,
	provider string,
	emailAddress string,
) (repo.EmailMailbox, error) {
	return s.repo.GetEmailMailboxByEmail(
		ctx,
		repo.GetEmailMailboxByEmailParams{
			Provider:     provider,
			EmailAddress: emailAddress,
		},
	)
}

func (s *svc) UpdateMailboxWatch(
	ctx context.Context,
	id pgtype.UUID,
	historyID string,
	watchExpiresAt pgtype.Timestamptz,
) (repo.EmailMailbox, error) {
	return s.repo.UpdateEmailMailboxWatch(
		ctx,
		repo.UpdateEmailMailboxWatchParams{
			ID: id,
			HistoryID: pgtype.Text{
				String: historyID,
				Valid:  true,
			},
			WatchExpiresAt: watchExpiresAt,
		},
	)
}

func (s *svc) CreateIngestionEvent(
	ctx context.Context,
	params repo.CreateEmailIngestionEventParams,
) (repo.EmailIngestionEvent, error) {
	return s.repo.CreateEmailIngestionEvent(ctx, params)
}

func (s *svc) GetIngestionEvent(
	ctx context.Context,
	id pgtype.UUID,
) (repo.EmailIngestionEvent, error) {
	return s.repo.GetEmailIngestionEvent(ctx, id)
}

func (s *svc) ListIngestionEvents(
	ctx context.Context,
) ([]repo.EmailIngestionEvent, error) {
	return s.repo.ListEmailIngestionEvents(ctx)
}

func (s *svc) MarkIngestionEventProcessing(
	ctx context.Context,
	id pgtype.UUID,
) (repo.EmailIngestionEvent, error) {
	return s.repo.MarkEmailIngestionEventProcessing(ctx, id)
}

func (s *svc) MarkIngestionEventProcessed(
	ctx context.Context,
	id pgtype.UUID,
) (repo.EmailIngestionEvent, error) {
	return s.repo.MarkEmailIngestionEventProcessed(ctx, id)
}

func (s *svc) MarkIngestionEventFailed(
	ctx context.Context,
	id pgtype.UUID,
	lastError string,
) (repo.EmailIngestionEvent, error) {
	return s.repo.MarkEmailIngestionEventFailed(
		ctx,
		repo.MarkEmailIngestionEventFailedParams{
			ID: id,
			LastError: pgtype.Text{
				String: lastError,
				Valid:  true,
			},
		},
	)
}
