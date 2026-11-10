package pdf

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

const (
	maxPDFStorageSize = 25 << 20
	pdfStorageDir     = "uploads"
)

type Storage struct {
	dir string
}

func NewStorage(dir string) *Storage {
	if dir == "" {
		dir = pdfStorageDir
	}

	return &Storage{
		dir: dir,
	}
}

func (s *Storage) Store(filename string, data []byte) (CreatePdfPayload, error) {
	return s.StoreReader(filename, bytes.NewReader(data), int64(len(data)))
}

func (s *Storage) StoreReader(filename string, reader io.Reader, maxSize int64) (CreatePdfPayload, error) {
	originalName := filepath.Base(strings.ReplaceAll(filename, "\\", "/"))

	if originalName == "." ||
		originalName == "" ||
		!strings.EqualFold(filepath.Ext(originalName), ".pdf") {
		return CreatePdfPayload{}, fmt.Errorf("uploaded file must have a .pdf extension")
	}

	if reader == nil {
		return CreatePdfPayload{}, fmt.Errorf("PDF data is missing")
	}

	if maxSize <= 0 || maxSize > maxPDFStorageSize {
		maxSize = maxPDFStorageSize
	}

	if err := os.MkdirAll(s.dir, 0750); err != nil {
		return CreatePdfPayload{}, fmt.Errorf("create PDF storage directory: %w", err)
	}

	dst, err := os.CreateTemp(s.dir, "pdf-*.pdf")
	if err != nil {
		return CreatePdfPayload{}, fmt.Errorf("create PDF storage file: %w", err)
	}

	storedPath := dst.Name()

	cleanup := func() {
		if err := os.Remove(storedPath); err != nil && !errors.Is(err, os.ErrNotExist) {
			// Best-effort cleanup.
		}
	}

	hash := sha256.New()

	limitedReader := io.LimitReader(reader, maxSize+1)

	size, err := io.Copy(io.MultiWriter(dst, hash), limitedReader)
	if err != nil {
		_ = dst.Close()
		cleanup()
		return CreatePdfPayload{}, fmt.Errorf("store PDF: %w", err)
	}

	if err := dst.Close(); err != nil {
		cleanup()
		return CreatePdfPayload{}, fmt.Errorf("close stored PDF: %w", err)
	}

	if size == 0 {
		cleanup()
		return CreatePdfPayload{}, fmt.Errorf("uploaded PDF is empty")
	}

	if size > maxSize {
		cleanup()
		return CreatePdfPayload{}, fmt.Errorf("PDF exceeds the allowed size limit")
	}

	return CreatePdfPayload{
		OriginalName: originalName,
		StoredPath:   filepath.ToSlash(storedPath),
		SizeBytes:    size,
		Sha256:       hex.EncodeToString(hash.Sum(nil)),
	}, nil
}

func (s *Storage) Remove(storedPath string) error {
	cleanPath := filepath.Clean(filepath.FromSlash(storedPath))

	relativePath, err := filepath.Rel(s.dir, cleanPath)
	if err != nil ||
		filepath.IsAbs(cleanPath) ||
		relativePath == ".." ||
		relativePath == "." ||
		strings.HasPrefix(relativePath, ".."+string(filepath.Separator)) {
		return fmt.Errorf("path is outside PDF storage directory")
	}

	if err := os.Remove(cleanPath); err != nil && !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("remove PDF: %w", err)
	}

	return nil
}
