package telegram

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
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

// Handler processes Telegram Bot webhook events.
type Handler struct {
	BotToken   string
	UploadDir  string
	Client     *Client
	Store      DocStore
	PdfService pdf.Service
	Log        *slog.Logger
}

// Receive handles incoming Telegram webhook updates (POST /webhooks/telegram).
func (h *Handler) Receive(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(io.LimitReader(r.Body, 1<<20))
	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	var update Update
	if err := json.Unmarshal(body, &update); err != nil {
		h.logger().Warn("invalid telegram update payload", "err", err)
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	// Always return 200 OK to Telegram so it does not hammer the server with retries
	w.WriteHeader(http.StatusOK)

	if update.Message == nil {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	if err := h.handleMessage(ctx, update.Message); err != nil {
		h.logger().Error("failed to process telegram message", "err", err)
	}
}

func (h *Handler) handleMessage(ctx context.Context, msg *Message) error {
	var (
		fileID       string
		filename     string
		mimeType     string
		isPhoto      bool
	)

	// Check if message is a document (e.g. PDF invoice)
	if msg.Document != nil {
		fileID = msg.Document.FileID
		filename = msg.Document.FileName
		mimeType = msg.Document.MimeType
	} else if len(msg.Photo) > 0 {
		// Take the highest resolution photo version
		bestPhoto := msg.Photo[len(msg.Photo)-1]
		fileID = bestPhoto.FileID
		filename = fmt.Sprintf("invoice_photo_%d.jpg", msg.MessageID)
		mimeType = "image/jpeg"
		isPhoto = true
	} else {
		// Text message or status
		if strings.HasPrefix(msg.Text, "/start") && msg.From != nil {
			_ = h.Client.SendMessage(ctx, msg.From.ID,
				"👋 Welcome to Autonomous 3-Way Invoice Reconciliation Bot!\n\n"+
					"Please attach and send your supplier Invoice PDF (or clear photo) here to queue it for 3-way matching.")
		}
		return nil
	}

	ext := filepath.Ext(filename)
	if ext == "" {
		if val, ok := allowedMime[mimeType]; ok {
			ext = val
		} else {
			ext = ".pdf"
		}
		filename += ext
	}

	msgIDStr := strconv.FormatInt(msg.MessageID, 10)

	// Idempotency check
	if h.Store != nil {
		exists, err := h.Store.Exists(ctx, "telegram", msgIDStr)
		if err != nil {
			return fmt.Errorf("check existing telegram message: %w", err)
		}
		if exists {
			h.logger().Info("duplicate telegram delivery ignored", "message_id", msgIDStr)
			return nil
		}
	}

	// Resolve sender name
	sender := "Telegram User"
	if msg.From != nil {
		if msg.From.Username != "" {
			sender = "@" + msg.From.Username
		} else if msg.From.FirstName != "" {
			sender = msg.From.FirstName
			if msg.From.LastName != "" {
				sender += " " + msg.From.LastName
			}
		}
	}

	// Download file bytes from Telegram
	var data []byte
	if strings.HasPrefix(fileID, "mock-") {
		// Mock testing mode
		data = []byte("%PDF-1.4 Fake Test Invoice Document Generated for Telegram Intake\n%EOF")
	} else {
		filePath, err := h.Client.GetFilePath(ctx, fileID)
		if err != nil {
			return fmt.Errorf("resolve telegram file path: %w", err)
		}

		data, err = h.Client.Download(ctx, filePath)
		if err != nil {
			return fmt.Errorf("download telegram file bytes: %w", err)
		}
	}

	uploadDir := h.UploadDir
	if uploadDir == "" {
		uploadDir = "uploads"
	}
	targetDir := filepath.Join(uploadDir, "telegram")
	if err := os.MkdirAll(targetDir, 0750); err != nil {
		return fmt.Errorf("create storage directory: %w", err)
	}

	safeName := unsafeChars.ReplaceAllString(filename, "_")
	savedPath := filepath.Join(targetDir, msgIDStr+"_"+safeName)
	if err := os.WriteFile(savedPath, data, 0640); err != nil {
		return fmt.Errorf("write telegram file to disk: %w", err)
	}

	sum := sha256.Sum256(data)
	shaHex := hex.EncodeToString(sum[:])

	// Record in PostgreSQL inbound_documents
	if h.Store != nil {
		_, err := h.Store.Insert(ctx, InboundDocument{
			Source:            "telegram",
			ExternalMessageID: msgIDStr,
			Sender:            sender,
			MediaID:           fileID,
			MimeType:          mimeType,
			OriginalFilename:  filename,
			StoragePath:       filepath.ToSlash(savedPath),
			SHA256:            shaHex,
			SizeBytes:         int64(len(data)),
			Status:            "RECEIVED",
		})
		if err != nil {
			return fmt.Errorf("insert inbound document: %w", err)
		}
	}

	// Sync with central pdfs table if it is a PDF
	if strings.EqualFold(ext, ".pdf") && h.PdfService != nil {
		storedName := filename
		if !strings.HasPrefix(strings.ToLower(filename), "[telegram]") {
			storedName = fmt.Sprintf("[telegram] %s", filename)
		}
		_, err := h.PdfService.CreatePdf(ctx, pdf.CreatePdfPayload{
			OriginalName: storedName,
			StoredPath:   filepath.ToSlash(savedPath),
			SizeBytes:    int64(len(data)),
			Sha256:       shaHex,
		})
		if err != nil {
			h.logger().Warn("could not record into pdfs table (duplicate or error)", "err", err)
		}
	}

	h.logger().Info("stored telegram document", "path", savedPath, "sender", sender)

	// Send confirmation message back to the sender on Telegram
	if msg.From != nil {
		replyText := fmt.Sprintf("✅ Document received: *%s*\nStatus: Queued for 3-Way Reconciliation\nHash: `%s...`",
			filename, shaHex[:12])
		if isPhoto {
			replyText = "📸 Invoice photo received!\nStatus: Queued for OCR & 3-Way Reconciliation"
		}
		_ = h.Client.SendMessage(ctx, msg.From.ID, replyText)
	}

	return nil
}

// GetFeed returns recent documents received via Telegram as JSON.
func (h *Handler) GetFeed(w http.ResponseWriter, r *http.Request) {
	if h.Store == nil {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte("[]"))
		return
	}

	docs, err := h.Store.ListRecent(r.Context(), "telegram", 10)
	if err != nil {
		h.logger().Error("failed to list recent telegram documents", "err", err)
		http.Error(w, "could not load telegram documents", http.StatusInternalServerError)
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

