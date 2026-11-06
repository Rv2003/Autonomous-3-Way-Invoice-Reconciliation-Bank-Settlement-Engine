package email

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	repo "github.com/Rv2003/Autonomous-3-Way-Invoice-Reconciliation-Bank-Settlement-Engine/BACKEND/internal/adapters/postgresql/sqlc"
	"github.com/Rv2003/Autonomous-3-Way-Invoice-Reconciliation-Bank-Settlement-Engine/BACKEND/internal/env"
	"github.com/Rv2003/Autonomous-3-Way-Invoice-Reconciliation-Bank-Settlement-Engine/BACKEND/internal/gmailoauth"
	appjson "github.com/Rv2003/Autonomous-3-Way-Invoice-Reconciliation-Bank-Settlement-Engine/BACKEND/internal/json"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"golang.org/x/oauth2"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

const gmailOAuthStateCookie = "gmail_oauth_state"

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

func (h *handler) StartGmailOAuth(w http.ResponseWriter, r *http.Request) {
	config, err := gmailoauth.Config(gmail.GmailReadonlyScope)
	if err != nil {
		http.Error(w, "Gmail OAuth is not configured", http.StatusInternalServerError)
		return
	}
	if config.RedirectURL == "" {
		http.Error(w, "GMAIL_REDIRECT_URL is not configured", http.StatusInternalServerError)
		return
	}

	stateBytes := make([]byte, 32)
	if _, err := rand.Read(stateBytes); err != nil {
		http.Error(w, "Could not start Gmail authorization", http.StatusInternalServerError)
		return
	}
	state := base64.RawURLEncoding.EncodeToString(stateBytes)
	http.SetCookie(w, &http.Cookie{
		Name:     gmailOAuthStateCookie,
		Value:    state,
		Path:     "/email/gmail/callback",
		HttpOnly: true,
		Secure:   r.TLS != nil || strings.EqualFold(os.Getenv("COOKIE_SECURE"), "true"),
		SameSite: http.SameSiteLaxMode,
		MaxAge:   600,
	})
	http.Redirect(w, r, config.AuthCodeURL(state, oauth2.AccessTypeOffline, oauth2.ApprovalForce), http.StatusFound)
}

func (h *handler) CompleteGmailOAuth(w http.ResponseWriter, r *http.Request) {
	stateCookie, err := r.Cookie(gmailOAuthStateCookie)
	if err != nil {
		http.Error(w, "Gmail authorization state is missing; start the connection again", http.StatusBadRequest)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:     gmailOAuthStateCookie,
		Value:    "",
		Path:     "/email/gmail/callback",
		HttpOnly: true,
		Secure:   r.TLS != nil || strings.EqualFold(os.Getenv("COOKIE_SECURE"), "true"),
		SameSite: http.SameSiteLaxMode,
		MaxAge:   -1,
		Expires:  time.Unix(0, 0),
	})

	state := r.URL.Query().Get("state")
	if state == "" || subtle.ConstantTimeCompare([]byte(state), []byte(stateCookie.Value)) != 1 {
		http.Error(w, "Invalid Gmail authorization state", http.StatusBadRequest)
		return
	}
	if r.URL.Query().Get("error") != "" {
		http.Error(w, "Google did not authorize Gmail access", http.StatusBadRequest)
		return
	}
	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "Google authorization code is missing", http.StatusBadRequest)
		return
	}

	config, err := gmailoauth.Config(gmail.GmailReadonlyScope)
	if err != nil {
		http.Error(w, "Gmail OAuth is not configured", http.StatusInternalServerError)
		return
	}
	token, err := config.Exchange(r.Context(), code)
	if err != nil {
		slog.Error("could not exchange Gmail OAuth code", "error", err)
		http.Error(w, "Could not complete Gmail authorization", http.StatusBadGateway)
		return
	}

	gmailService, err := gmail.NewService(
		r.Context(),
		option.WithHTTPClient(config.Client(r.Context(), token)),
	)
	if err != nil {
		http.Error(w, "Could not connect to Gmail", http.StatusBadGateway)
		return
	}
	profile, err := gmailService.Users.GetProfile("me").Context(r.Context()).Do()
	if err != nil {
		slog.Error("could not read authorized Gmail profile", "error", err)
		http.Error(w, "Could not verify the authorized Gmail account", http.StatusBadGateway)
		return
	}

	_, err = h.service.CreateMailbox(r.Context(), repo.CreateEmailMailboxParams{
		Provider:     "gmail",
		EmailAddress: profile.EmailAddress,
	})
	if err != nil {
		var pgErr interface{ SQLState() string }
		if !errors.As(err, &pgErr) || pgErr.SQLState() != "23505" {
			slog.Error("could not register Gmail mailbox", "error", err)
			http.Error(w, "Could not register Gmail mailbox", http.StatusInternalServerError)
			return
		}
		mailbox, getErr := h.service.GetMailboxByEmail(r.Context(), "gmail", profile.EmailAddress)
		if getErr != nil || mailbox.Status != "active" {
			http.Error(w, "The existing Gmail mailbox is not active", http.StatusConflict)
			return
		}
	}

	if err := gmailoauth.SaveToken(token); err != nil {
		slog.Error("could not save Gmail OAuth token", "error", err)
		http.Error(w, "Could not securely save Gmail authorization", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, env.GetString("GMAIL_FRONTEND_URL", "http://localhost:5173/"), http.StatusSeeOther)
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
	HistoryID    string `json:"historyId"`
}

func (h *handler) GmailWebhook(w http.ResponseWriter, r *http.Request) {
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
	notification.HistoryID = strings.TrimSpace(notification.HistoryID)

	if notification.EmailAddress == "" {
		http.Error(w, "emailAddress is required", http.StatusBadRequest)
		return
	}

	if notification.HistoryID == "" {
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
			HistoryID:       notification.HistoryID,
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

func (h *handler) ListGmailMessages(w http.ResponseWriter, r *http.Request) {
	limit := int64(10)
	if value := r.URL.Query().Get("limit"); value != "" {
		parsedLimit, err := strconv.ParseInt(value, 10, 64)
		if err != nil || parsedLimit < 1 || parsedLimit > 50 {
			http.Error(w, "limit must be between 1 and 50", http.StatusBadRequest)
			return
		}
		limit = parsedLimit
	}

	inbox, err := gmailoauth.FetchInbox(r.Context(), limit)
	if err != nil {
		slog.Error("could not fetch Gmail inbox", "error", err)
		http.Error(w, "Could not fetch Gmail inbox. Check OAuth configuration and token.", http.StatusBadGateway)
		return
	}

	mailbox, err := h.service.GetMailboxByEmail(r.Context(), "gmail", inbox.Email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			http.Error(w, "Gmail account is not registered as a mailbox", http.StatusConflict)
			return
		}
		slog.Error("could not verify Gmail mailbox", "error", err)
		http.Error(w, "Could not verify Gmail mailbox", http.StatusInternalServerError)
		return
	}
	if mailbox.Status != "active" {
		http.Error(w, "Gmail mailbox is not active", http.StatusConflict)
		return
	}

	appjson.Write(w, http.StatusOK, inbox)
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
