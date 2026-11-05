package pdf

import repo "github.com/Rv2003/Autonomous-3-Way-Invoice-Reconciliation-Bank-Settlement-Engine/BACKEND/internal/adapters/postgresql/sqlc"

type CreatePdfPayload struct {
	OriginalName string
	StoredPath   string
	SizeBytes    int64
	Sha256       string
}

func (payload CreatePdfPayload) toRepoParams() repo.CreatePdfParams {
	return repo.CreatePdfParams{
		OriginalName: payload.OriginalName,
		StoredPath:   payload.StoredPath,
		SizeBytes:    payload.SizeBytes,
		Sha256:       payload.Sha256,
	}
}
