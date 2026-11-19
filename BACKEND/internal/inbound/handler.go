package inbound

import (
	"errors"
	"fmt"
	"log"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) ViewDocument(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	idStr = strings.TrimPrefix(idStr, "email-")
	idStr = strings.TrimPrefix(idStr, "tg-")

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid inbound document ID", http.StatusBadRequest)
		return
	}

	doc, err := h.service.GetDocument(r.Context(), id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			http.Error(w, "inbound document not found", http.StatusNotFound)
			return
		}
		log.Printf("load inbound document metadata: %v", err)
		http.Error(w, "could not load inbound document", http.StatusInternalServerError)
		return
	}

	filePath := filepath.Clean(filepath.FromSlash(doc.StoragePath))
	file, err := os.Open(filePath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			http.Error(w, "document file not found on disk", http.StatusNotFound)
			return
		}
		log.Printf("open inbound document %q: %v", filePath, err)
		http.Error(w, "could not open document", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		log.Printf("stat inbound document %q: %v", filePath, err)
		http.Error(w, "could not open document", http.StatusInternalServerError)
		return
	}
	if !info.Mode().IsRegular() {
		log.Printf("inbound document is not a regular file: %q", filePath)
		http.Error(w, "could not open document", http.StatusInternalServerError)
		return
	}

	contentType := doc.MimeType
	if contentType == "" || strings.EqualFold(contentType, "application/octet-stream") {
		contentType = "application/pdf"
	}

	filename := doc.OriginalFilename
	if filename == "" {
		filename = fmt.Sprintf("document-%d.pdf", doc.ID)
	}

	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Content-Disposition", mime.FormatMediaType("inline", map[string]string{
		"filename": filename,
	}))
	w.Header().Set("X-Content-Type-Options", "nosniff")
	http.ServeContent(w, r, filename, info.ModTime(), file)
}

func (h *Handler) DeleteDocument(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	idStr = strings.TrimPrefix(idStr, "email-")
	idStr = strings.TrimPrefix(idStr, "tg-")

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid inbound document ID", http.StatusBadRequest)
		return
	}

	row, err := h.service.DeleteDocument(r.Context(), id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			http.Error(w, "inbound document not found", http.StatusNotFound)
			return
		}
		log.Printf("delete inbound document metadata: %v", err)
		http.Error(w, "could not delete document", http.StatusInternalServerError)
		return
	}

	if row.StoragePath != "" {
		cleanPath := filepath.Clean(filepath.FromSlash(row.StoragePath))
		_ = os.Remove(cleanPath)
	}

	w.WriteHeader(http.StatusNoContent)
}