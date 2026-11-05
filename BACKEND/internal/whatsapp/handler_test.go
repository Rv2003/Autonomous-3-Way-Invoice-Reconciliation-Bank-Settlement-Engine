package whatsapp

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
)

const (
	testSecret = "test-app-secret-1234567890"
	testVerify = "test-verify-token"
	testToken  = "test-access-token"
)

type fakeStore struct {
	mu   sync.Mutex
	rows map[string]InboundDocument
}

func newFakeStore() *fakeStore {
	return &fakeStore{rows: make(map[string]InboundDocument)}
}

func (f *fakeStore) Exists(_ context.Context, src, id string) (bool, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	_, ok := f.rows[src+":"+id]
	return ok, nil
}

func (f *fakeStore) Insert(_ context.Context, d InboundDocument) (bool, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	key := d.Source + ":" + d.ExternalMessageID
	if _, ok := f.rows[key]; ok {
		return false, nil
	}
	f.rows[key] = d
	return true, nil
}

func (f *fakeStore) ListRecent(_ context.Context, src string, limit int) ([]InboundDocument, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	var docs []InboundDocument
	for _, row := range f.rows {
		if row.Source == src {
			docs = append(docs, row)
		}
	}
	return docs, nil
}

func fakeGraph(t *testing.T, file []byte) *httptest.Server {
	t.Helper()
	sum := sha256.Sum256(file)
	b64 := base64.StdEncoding.EncodeToString(sum[:])

	mux := http.NewServeMux()
	mux.HandleFunc("/media-123", func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer "+testToken {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"url":       "http://" + r.Host + "/download/media-123",
			"mime_type": "application/pdf",
			"sha256":    b64,
			"file_size": len(file),
		})
	})
	mux.HandleFunc("/download/media-123", func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer "+testToken {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		_, _ = w.Write(file)
	})

	srv := httptest.NewServer(mux)
	t.Cleanup(srv.Close)
	return srv
}

func sign(body []byte) string {
	m := hmac.New(sha256.New, []byte(testSecret))
	m.Write(body)
	return "sha256=" + hex.EncodeToString(m.Sum(nil))
}

func documentPayload(msgID string) []byte {
	return fmt.Appendf(nil, `{
		"object": "whatsapp_business_account",
		"entry": [{
			"changes": [{
				"field": "messages",
				"value": {
					"messages": [{
						"from": "94771234567",
						"id": "%s",
						"type": "document",
						"document": {
							"id": "media-123",
							"mime_type": "application/pdf",
							"filename": "Invoice_9921.pdf"
						}
					}]
				}
			}]
		}]
	}`, msgID)
}

func newHandler(t *testing.T, graphURL string, store DocStore) *Handler {
	return &Handler{
		AppSecret:   testSecret,
		VerifyToken: testVerify,
		UploadDir:   t.TempDir(),
		Media: &MediaClient{
			BaseURL: graphURL,
			Token:   testToken,
			HTTP:    http.DefaultClient,
		},
		Store: store,
		Log:   slog.New(slog.NewTextHandler(os.Stderr, nil)),
	}
}

func TestSignature(t *testing.T) {
	body := []byte(`{"event":"test"}`)
	validSig := sign(body)

	if !ValidSignature(testSecret, body, validSig) {
		t.Fatal("valid signature was rejected")
	}
	if ValidSignature(testSecret, body, sign([]byte(`{"event":"different"}`))) {
		t.Fatal("signature for different body was accepted")
	}
	if ValidSignature("wrong-secret", body, validSig) {
		t.Fatal("wrong secret signature was accepted")
	}
	if ValidSignature(testSecret, body, "") {
		t.Fatal("empty signature header was accepted")
	}
}

func TestVerifyChallenge(t *testing.T) {
	h := newHandler(t, "", newFakeStore())

	// Successful verification
	req := httptest.NewRequest(http.MethodGet,
		"/webhooks/whatsapp?hub.mode=subscribe&hub.verify_token="+testVerify+"&hub.challenge=challenge12345", nil)
	rec := httptest.NewRecorder()
	h.Verify(rec, req)

	if rec.Code != http.StatusOK || rec.Body.String() != "challenge12345" {
		t.Fatalf("got status %d body %q, want 200 'challenge12345'", rec.Code, rec.Body.String())
	}

	// Token mismatch
	badReq := httptest.NewRequest(http.MethodGet,
		"/webhooks/whatsapp?hub.mode=subscribe&hub.verify_token=wrong-token&hub.challenge=challenge12345", nil)
	badRec := httptest.NewRecorder()
	h.Verify(badRec, badReq)

	if badRec.Code != http.StatusForbidden {
		t.Fatalf("got status %d, want 403 Forbidden", badRec.Code)
	}
}

func TestReceiveStoresFileAndIsIdempotent(t *testing.T) {
	fileContent := []byte("%PDF-1.4 Fake Test Invoice Document Content")
	graph := fakeGraph(t, fileContent)
	store := newFakeStore()
	h := newHandler(t, graph.URL, store)

	body := documentPayload("wamid.ABC123XYZ")
	sendWebhook := func() int {
		req := httptest.NewRequest(http.MethodPost, "/webhooks/whatsapp", strings.NewReader(string(body)))
		req.Header.Set("X-Hub-Signature-256", sign(body))
		rec := httptest.NewRecorder()
		h.Receive(rec, req)
		return rec.Code
	}

	// 1. Initial delivery
	if code := sendWebhook(); code != http.StatusOK {
		t.Fatalf("first delivery failed with status %d, want 200", code)
	}

	// 2. Duplicate delivery (Meta retry)
	if code := sendWebhook(); code != http.StatusOK {
		t.Fatalf("duplicate delivery failed with status %d, want 200", code)
	}

	// Check idempotency in storage
	if len(store.rows) != 1 {
		t.Fatalf("expected 1 row in store, found %d (idempotency failed)", len(store.rows))
	}

	var savedDoc InboundDocument
	for _, row := range store.rows {
		savedDoc = row
	}

	savedBytes, err := os.ReadFile(savedDoc.StoragePath)
	if err != nil {
		t.Fatalf("could not read saved document from disk: %v", err)
	}
	if string(savedBytes) != string(fileContent) {
		t.Fatal("saved file bytes do not match expected downloaded content")
	}
	if !strings.HasPrefix(filepath.Base(savedDoc.StoragePath), "wamid.ABC123XYZ_") {
		t.Fatalf("file name %q missing expected message ID prefix", savedDoc.StoragePath)
	}
}

func TestReceiveRejectsBadSignature(t *testing.T) {
	h := newHandler(t, "", newFakeStore())
	body := documentPayload("wamid.INVALID_SIG")

	req := httptest.NewRequest(http.MethodPost, "/webhooks/whatsapp", strings.NewReader(string(body)))
	req.Header.Set("X-Hub-Signature-256", "sha256=deadbeefcafebabe")
	rec := httptest.NewRecorder()
	h.Receive(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("got status %d, want 401 Unauthorized", rec.Code)
	}
}

