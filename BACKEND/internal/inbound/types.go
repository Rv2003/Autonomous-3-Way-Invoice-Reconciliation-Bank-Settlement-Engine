package inbound

import repo "github.com/Rv2003/Autonomous-3-Way-Invoice-Reconciliation-Bank-Settlement-Engine/BACKEND/internal/adapters/postgresql/sqlc"

type CreateInboundDocumentPayload struct {
	Source            string
	ExternalMessageID string
	Sender            string
	MediaID           string
	MimeType          string
	OriginalFilename  string
	StoragePath       string
	Sha256            string
	SizeBytes         int64
	Status            string
}

func (payload CreateInboundDocumentPayload) toRepoParams() repo.CreateInboundDocumentParams {
	return repo.CreateInboundDocumentParams{
		Source:            payload.Source,
		ExternalMessageID: payload.ExternalMessageID,
		Sender:            payload.Sender,
		MediaID:           payload.MediaID,
		MimeType:          payload.MimeType,
		OriginalFilename:  payload.OriginalFilename,
		StoragePath:       payload.StoragePath,
		Sha256:            payload.Sha256,
		SizeBytes:         payload.SizeBytes,
		Status:            payload.Status,
	}
}