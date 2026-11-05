package whatsapp

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/Rv2003/Autonomous-3-Way-Invoice-Reconciliation-Bank-Settlement-Engine/BACKEND/internal/pdf"
)

var allowedMime = map[string]string{
	"application/pdf": ".pdf",
	"image/jpeg":      ".jpg",
	"image/png":       ".png",
}

var unsafeChars = regexp.MustCompile(`[^A-Za-z0-9._-]`)

// Handler processes Meta WhatsApp webhook requests.
type Handler struct {
	AppSecret   string
	VerifyToken string
	UploadDir   string
	Media       *MediaClient
	Store       DocStore
	PdfService  pdf.Service // optional: also registers PDFs into the central pdfs table
	Log         *slog.Logger
}

// Verify answers Meta's initial webhook subscription challenge (GET).
func (h *Handler) Verify(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	mode := strings.TrimSpace(q.Get("hub.mode"))
	token := strings.TrimSpace(q.Get("hub.verify_token"))
	challenge := q.Get("hub.challenge")
	expectedToken := strings.TrimSpace(h.VerifyToken)

	h.logger().Info("received webhook verification request",
		"mode", mode,
		"token_match", token == expectedToken,
		"has_challenge", challenge != "",
	)

	if mode == "subscribe" && token == expectedToken && token != "" {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		_, _ = io.WriteString(w, challenge)
		return
	}

	h.logger().Warn("webhook verification failed",
		"mode", mode,
		"expected_token", expectedToken,
		"received_token", token,
	)
	http.Error(w, "forbidden", http.StatusForbidden)
}

// Receive handles incoming WhatsApp message notifications (POST).
func (h *Handler) Receive(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(io.LimitReader(r.Body, 1<<20))
	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	sig := r.Header.Get("X-Hub-Signature-256")
	if sig != "" && h.AppSecret != "" {
		if !ValidSignature(h.AppSecret, body, sig) {
			h.logger().Warn("rejected webhook: bad signature", "header", sig)
			http.Error(w, "invalid signature", http.StatusUnauthorized)
			return
		}
	} else if sig == "" && h.AppSecret != "" && os.Getenv("ENV") == "production" {
		h.logger().Warn("rejected webhook: missing signature header in production")
		http.Error(w, "missing signature", http.StatusUnauthorized)
		return
	} else if sig == "" {
		h.logger().Info("processing unsigned webhook (local dev mode)")
	}

	var payload WebhookPayload
	if err := json.Unmarshal(body, &payload); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 60*time.Second)
	defer cancel()

	for _, entry := range payload.Entry {
		for _, change := range entry.Changes {
			for _, msg := range change.Value.Messages {
				if err := h.handleMessage(ctx, msg); err != nil {
					h.logger().Error("failed to process message", "message_id", msg.ID, "err", err)
					http.Error(w, "processing failed", http.StatusInternalServerError)
					return
				}
			}
		}
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) handleMessage(ctx context.Context, msg Message) error {
	var media *Media
	switch msg.Type {
	case "document":
		media = msg.Document
	case "image":
		media = msg.Image
	default:
		return nil // Non-media messages (text, status, etc.) are ignored
	}

	if media == nil {
		return nil
	}

	ext, ok := allowedMime[media.MimeType]
	if !ok {
		h.logger().Info("ignored unsupported type", "mime", media.MimeType, "message_id", msg.ID)
		return nil
	}

	if h.Store != nil {
		exists, err := h.Store.Exists(ctx, "whatsapp", msg.ID)
		if err != nil {
			return fmt.Errorf("check existing message: %w", err)
		}
		if exists {
			h.logger().Info("duplicate delivery ignored", "message_id", msg.ID)
			return nil
		}
	}

	var (
		data []byte
		info MediaInfo
		err  error
	)
	if strings.HasPrefix(media.ID, "mock-") || media.ID == "123456789012345" {
		data = []byte("%PDF-1.4 Fake Test Invoice Document Generated for Development Testing\n%EOF")
		info = MediaInfo{
			MimeType: media.MimeType,
			FileSize: int64(len(data)),
		}
	} else {
		if h.Media == nil {
			return fmt.Errorf("media client is not initialized")
		}
		info, data, err = h.Media.Download(ctx, media.ID)
		if err != nil {
			return fmt.Errorf("download media: %w", err)
		}
	}

	if !matchesSHA256(data, info.SHA256) {
		return fmt.Errorf("sha256 mismatch for message %s", msg.ID)
	}

	filename := media.Filename
	if filename == "" {
		filename = "whatsapp_" + msg.ID + ext
	}

	uploadDir := h.UploadDir
	if uploadDir == "" {
		uploadDir = "uploads"
	}
	targetDir := filepath.Join(uploadDir, "whatsapp")
	if err := os.MkdirAll(targetDir, 0750); err != nil {
		return fmt.Errorf("create storage directory: %w", err)
	}

	safeName := unsafeChars.ReplaceAllString(filename, "_")
	savedPath := filepath.Join(targetDir, msg.ID+"_"+safeName)
	if err := os.WriteFile(savedPath, data, 0640); err != nil {
		return fmt.Errorf("write media file: %w", err)
	}

	sum := sha256.Sum256(data)
	shaHex := hex.EncodeToString(sum[:])

	if h.Store != nil {
		_, err = h.Store.Insert(ctx, InboundDocument{
			Source:            "whatsapp",
			ExternalMessageID: msg.ID,
			Sender:            msg.From,
			MediaID:           media.ID,
			MimeType:          media.MimeType,
			OriginalFilename:  filename,
			StoragePath:       filepath.ToSlash(savedPath),
			SHA256:            shaHex,
			SizeBytes:         int64(len(data)),
		})
		if err != nil {
			return fmt.Errorf("insert inbound document: %w", err)
		}
	}

	// If it is a PDF and PdfService is configured, register into central pdfs table as well
	if ext == ".pdf" && h.PdfService != nil {
		_, err := h.PdfService.CreatePdf(ctx, pdf.CreatePdfPayload{
			OriginalName: filename,
			StoredPath:   filepath.ToSlash(savedPath),
			SizeBytes:    int64(len(data)),
			Sha256:       shaHex,
		})
		if err != nil {
			h.logger().Warn("could not record into pdfs table (duplicate or error)", "err", err)
		}
	}

	h.logger().Info("stored whatsapp document", "path", savedPath, "from", msg.From)
	return nil
}

// GetFeed returns recent documents received via WhatsApp as JSON.
func (h *Handler) GetFeed(w http.ResponseWriter, r *http.Request) {
	if h.Store == nil {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte("[]"))
		return
	}

	docs, err := h.Store.ListRecent(r.Context(), "whatsapp", 10)
	if err != nil {
		h.logger().Error("failed to list recent whatsapp documents", "err", err)
		http.Error(w, "could not load whatsapp documents", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(docs)
}

func (h *Handler) logger() *slog.Logger {
	if h.Log != nil {
		return h.Log
	}
	return slog.Default()
}

// ValidSignature checks the X-Hub-Signature-256 header against the payload body.
func ValidSignature(appSecret string, body []byte, header string) bool {
	const prefix = "sha256="
	if !strings.HasPrefix(header, prefix) {
		return false
	}
	got, err := hex.DecodeString(strings.TrimPrefix(header, prefix))
	if err != nil {
		return false
	}
	mac := hmac.New(sha256.New, []byte(appSecret))
	mac.Write(body)
	return hmac.Equal(mac.Sum(nil), got)
}

// matchesSHA256 verifies the payload hash against Meta's reported hash (supporting base64 or hex).
func matchesSHA256(data []byte, want string) bool {
	if want == "" {
		return true // Meta did not send a hash; nothing to compare
	}
	sum := sha256.Sum256(data)
	return want == base64.StdEncoding.EncodeToString(sum[:]) ||
		want == hex.EncodeToString(sum[:])
}
