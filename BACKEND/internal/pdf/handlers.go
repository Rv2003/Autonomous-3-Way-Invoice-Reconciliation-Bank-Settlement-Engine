package pdf

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/Rv2003/Autonomous-3-Way-Invoice-Reconciliation-Bank-Settlement-Engine/BACKEND/internal/json"
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
