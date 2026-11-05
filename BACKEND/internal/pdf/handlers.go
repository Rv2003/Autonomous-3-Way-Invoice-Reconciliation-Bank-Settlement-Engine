package pdf

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"log"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/Rv2003/Autonomous-3-Way-Invoice-Reconciliation-Bank-Settlement-Engine/BACKEND/internal/json"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
)

const (
	maxPDFUploadSize = 25 << 20
	pdfUploadDir     = "uploads"
)

type handler struct {
	service Service
}

func NewHandler(service Service) *handler {
	return &handler{
		service: service,
	}
}

func (h *handler) ListPdf(w http.ResponseWriter, r *http.Request) {
	pdfs, err := h.service.ListPdf(r.Context())
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.Write(w, http.StatusOK, pdfs)
}

func (h *handler) ViewPdf(w http.ResponseWriter, r *http.Request) {
	pdf, err := h.service.GetPdf(r.Context(), chi.URLParam(r, "id"))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			http.Error(w, "Document not found", http.StatusNotFound)
			return
		}
		if errors.Is(err, ErrInvalidDocumentID) {
			http.Error(w, "Invalid document ID", http.StatusBadRequest)
			return
		}
		log.Printf("load PDF metadata: %v", err)
		http.Error(w, "Could not load document", http.StatusInternalServerError)
		return
	}

	filePath, err := safeUploadPath(pdf.StoredPath)
	if err != nil {
		log.Printf("refusing to serve PDF outside upload directory: %q", pdf.StoredPath)
		http.Error(w, "Could not open document", http.StatusInternalServerError)
		return
	}
	file, err := os.Open(filePath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			http.Error(w, "Document file not found", http.StatusNotFound)
			return
		}
		log.Printf("open PDF %q: %v", filePath, err)
		http.Error(w, "Could not open document", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		log.Printf("stat PDF %q: %v", filePath, err)
		http.Error(w, "Could not open document", http.StatusInternalServerError)
		return
	}
	if !info.Mode().IsRegular() {
		log.Printf("PDF is not a regular file: %q", filePath)
		http.Error(w, "Could not open document", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", mime.FormatMediaType("inline", map[string]string{
		"filename": pdf.OriginalName,
	}))
	w.Header().Set("X-Content-Type-Options", "nosniff")
	http.ServeContent(w, r, pdf.OriginalName, info.ModTime(), file)
}

func (h *handler) DeletePdf(w http.ResponseWriter, r *http.Request) {
	storedPath, err := h.service.DeletePdf(r.Context(), chi.URLParam(r, "id"))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			http.Error(w, "Document not found", http.StatusNotFound)
			return
		}
		log.Printf("delete PDF metadata: %v", err)
		http.Error(w, "Could not delete document", http.StatusInternalServerError)
		return
	}

	cleanPath, err := safeUploadPath(storedPath)
	if err != nil {
		log.Printf("refusing to remove PDF outside upload directory: %q", storedPath)
		http.Error(w, "Document was removed from the list, but its file could not be safely deleted", http.StatusInternalServerError)
		return
	}
	if err := os.Remove(cleanPath); err != nil && !errors.Is(err, os.ErrNotExist) {
		log.Printf("remove deleted PDF file %q: %v", cleanPath, err)
		http.Error(w, "Document was removed from the list, but its file could not be deleted", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func safeUploadPath(storedPath string) (string, error) {
	cleanPath := filepath.Clean(filepath.FromSlash(storedPath))
	relativePath, err := filepath.Rel(pdfUploadDir, cleanPath)
	if err != nil || filepath.IsAbs(cleanPath) || relativePath == ".." ||
		relativePath == "." || strings.HasPrefix(relativePath, ".."+string(filepath.Separator)) {
		return "", fmt.Errorf("path is outside upload directory")
	}
	return cleanPath, nil
}

func (h *handler) CreatePdf(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, maxPDFUploadSize)
	if err := r.ParseMultipartForm(1 << 20); err != nil {
		var maxBytesErr *http.MaxBytesError
		if errors.As(err, &maxBytesErr) {
			http.Error(w, "PDF upload exceeds the 25 MB limit", http.StatusRequestEntityTooLarge)
			return
		}
		http.Error(w, "invalid multipart form; upload the PDF in the file field", http.StatusBadRequest)
		return
	}
	if r.MultipartForm != nil {
		defer r.MultipartForm.RemoveAll()
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "missing PDF file in the file field", http.StatusBadRequest)
		return
	}
	defer file.Close()

	originalName := filepath.Base(strings.ReplaceAll(header.Filename, "\\", "/"))
	if originalName == "." || originalName == "" || !strings.EqualFold(filepath.Ext(originalName), ".pdf") {
		http.Error(w, "uploaded file must have a .pdf extension", http.StatusBadRequest)
		return
	}

	if err := os.MkdirAll(pdfUploadDir, 0750); err != nil {
		log.Printf("create PDF upload directory: %v", err)
		http.Error(w, "could not store PDF", http.StatusInternalServerError)
		return
	}

	dst, err := os.CreateTemp(pdfUploadDir, "pdf-*.pdf")
	if err != nil {
		log.Printf("create PDF upload file: %v", err)
		http.Error(w, "could not store PDF", http.StatusInternalServerError)
		return
	}
	storedPath := dst.Name()
	cleanup := func() {
		if err := os.Remove(storedPath); err != nil && !errors.Is(err, os.ErrNotExist) {
			log.Printf("remove unrecorded PDF %q: %v", storedPath, err)
		}
	}

	hash := sha256.New()
	size, copyErr := io.Copy(io.MultiWriter(dst, hash), file)
	closeErr := dst.Close()
	if copyErr != nil || closeErr != nil {
		cleanup()
		if copyErr != nil {
			log.Printf("store uploaded PDF: %v", copyErr)
		} else {
			log.Printf("close stored PDF: %v", closeErr)
		}
		http.Error(w, "could not store PDF", http.StatusInternalServerError)
		return
	}
	if size == 0 {
		cleanup()
		http.Error(w, "uploaded PDF is empty", http.StatusBadRequest)
		return
	}

	payload := CreatePdfPayload{
		OriginalName: originalName,
		StoredPath:   filepath.ToSlash(storedPath),
		SizeBytes:    size,
		Sha256:       hex.EncodeToString(hash.Sum(nil)),
	}
	pdf, err := h.service.CreatePdf(r.Context(), payload)
	if err != nil {
		cleanup()
		log.Println(err)
		http.Error(w, fmt.Sprintf("could not save PDF metadata: %v", err), http.StatusInternalServerError)
		return
	}

	json.Write(w, http.StatusCreated, pdf)
}
