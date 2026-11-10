package email

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"

	repo "github.com/Rv2003/Autonomous-3-Way-Invoice-Reconciliation-Bank-Settlement-Engine/BACKEND/internal/adapters/postgresql/sqlc"
	appjson "github.com/Rv2003/Autonomous-3-Way-Invoice-Reconciliation-Bank-Settlement-Engine/BACKEND/internal/json"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type handler struct {
	service Service
}

func NewHandler(service Service) *handler {
	return &handler{
		service: service,
	}
}

type createMailboxRequest struct {
	Provider     string `json:"provider"`
	EmailAddress string `json:"email_address"`
}

func (h *handler) CreateMailbox(w http.ResponseWriter, r *http.Request) {
	var req createMailboxRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON request body", http.StatusBadRequest)
		return
	}

	req.Provider = strings.ToLower(strings.TrimSpace(req.Provider))
	req.EmailAddress = strings.TrimSpace(req.EmailAddress)

	if req.Provider != "gmail" {
		http.Error(w, "provider must be gmail", http.StatusBadRequest)
		return
	}

	if req.EmailAddress == "" {
		http.Error(w, "email_address is required", http.StatusBadRequest)
		return
	}

	mailbox, err := h.service.CreateMailbox(
		r.Context(),
		repo.CreateEmailMailboxParams{
			Provider:     req.Provider,
			EmailAddress: req.EmailAddress,
		},
	)
	if err != nil {
		var pgErr interface {
			SQLState() string
		}

		if errors.As(err, &pgErr) && pgErr.SQLState() == "23505" {
			http.Error(w, "email mailbox already exists", http.StatusConflict)
			return
		}

		http.Error(w, "could not create email mailbox", http.StatusInternalServerError)
		return
	}

	appjson.Write(w, http.StatusCreated, mailbox)
}

type gmailPubSubRequest struct {
	Message gmailPubSubMessage `json:"message"`
}

type gmailPubSubMessage struct {
	MessageID string `json:"messageId"`
	Data      string `json:"data"`
}

type gmailNotification struct {
	EmailAddress string `json:"emailAddress"`
	HistoryID    int64  `json:"historyId"`
}

func (h *handler) GmailWebhook(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "could not read request body", http.StatusBadRequest)
		return
	}

	log.Printf("GMAIL WEBHOOK RAW BODY: %s", string(body))

	r.Body = io.NopCloser(bytes.NewReader(body))

	var req gmailPubSubRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON request body", http.StatusBadRequest)
		return
	}

	if req.Message.MessageID == "" {
		http.Error(w, "message.messageId is required", http.StatusBadRequest)
		return
	}

	if req.Message.Data == "" {
		http.Error(w, "message.data is required", http.StatusBadRequest)
		return
	}

	data, err := base64.StdEncoding.DecodeString(req.Message.Data)
	if err != nil {
		http.Error(w, "invalid message.data encoding", http.StatusBadRequest)
		return
	}

	var notification gmailNotification

	if err := json.Unmarshal(data, &notification); err != nil {
		http.Error(w, "invalid Gmail notification data", http.StatusBadRequest)
		return
	}

	notification.EmailAddress = strings.TrimSpace(notification.EmailAddress)

	if notification.EmailAddress == "" {
		http.Error(w, "emailAddress is required", http.StatusBadRequest)
		return
	}

	if notification.HistoryID == 0 {
		http.Error(w, "historyId is required", http.StatusBadRequest)
		return
	}

	mailbox, err := h.service.GetMailboxByEmail(
		r.Context(),
		"gmail",
		notification.EmailAddress,
	)
	if err != nil {
		http.Error(w, "mailbox not found", http.StatusNotFound)
		return
	}

	event, err := h.service.CreateIngestionEvent(
		r.Context(),
		repo.CreateEmailIngestionEventParams{
			MailboxID:       mailbox.ID,
			PubsubMessageID: req.Message.MessageID,
			HistoryID:       strconv.FormatInt(notification.HistoryID, 10),
		},
	)
	if err != nil {
		var pgErr interface {
			SQLState() string
		}

		if errors.As(err, &pgErr) && pgErr.SQLState() == "23505" {
			w.WriteHeader(http.StatusOK)
			return
		}

		http.Error(w, "could not create ingestion event", http.StatusInternalServerError)
		return
	}

	appjson.Write(w, http.StatusAccepted, event)
}

func (h *handler) ListIngestionEvents(w http.ResponseWriter, r *http.Request) {
	events, err := h.service.ListIngestionEvents(r.Context())
	if err != nil {
		http.Error(w, "could not list ingestion events", http.StatusInternalServerError)
		return
	}

	appjson.Write(w, http.StatusOK, events)
}

func (h *handler) ProcessIngestionEvent(w http.ResponseWriter, r *http.Request) {
	idString := chi.URLParam(r, "id")

	parsedID, err := uuid.Parse(idString)
	if err != nil {
		http.Error(w, "invalid ingestion event id", http.StatusBadRequest)
		return
	}

	id := pgtype.UUID{
		Bytes: parsedID,
		Valid: true,
	}

	event, err := h.service.MarkIngestionEventProcessing(
		r.Context(),
		id,
	)
	if err != nil {
		http.Error(w, "could not process ingestion event", http.StatusConflict)
		return
	}

	appjson.Write(w, http.StatusOK, event)
}

func (h *handler) CompleteIngestionEvent(w http.ResponseWriter, r *http.Request) {
	idString := chi.URLParam(r, "id")

	parsedID, err := uuid.Parse(idString)
	if err != nil {
		http.Error(w, "invalid ingestion event id", http.StatusBadRequest)
		return
	}

	id := pgtype.UUID{
		Bytes: parsedID,
		Valid: true,
	}

	event, err := h.service.MarkIngestionEventProcessed(
		r.Context(),
		id,
	)
	if err != nil {
		http.Error(w, "could not complete ingestion event", http.StatusConflict)
		return
	}

	appjson.Write(w, http.StatusOK, event)
}

func (h *handler) FailIngestionEvent(w http.ResponseWriter, r *http.Request) {
	idString := chi.URLParam(r, "id")

	parsedID, err := uuid.Parse(idString)
	if err != nil {
		http.Error(w, "invalid ingestion event id", http.StatusBadRequest)
		return
	}

	lastError := strings.TrimSpace(r.URL.Query().Get("error"))
	if lastError == "" {
		http.Error(w, "error query parameter is required", http.StatusBadRequest)
		return
	}

	id := pgtype.UUID{
		Bytes: parsedID,
		Valid: true,
	}

	event, err := h.service.MarkIngestionEventFailed(
		r.Context(),
		id,
		lastError,
	)
	if err != nil {
		http.Error(w, "could not fail ingestion event", http.StatusConflict)
		return
	}

	appjson.Write(w, http.StatusOK, event)
}
