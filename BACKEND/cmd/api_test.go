package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLocalFrontendCORS(t *testing.T) {
	allowedOrigins, err := parseCORSAllowedOrigins("https://app.example.com")
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name            string
		method          string
		origin          string
		wantStatus      int
		wantAllowOrigin string
		wantNext        bool
	}{
		{name: "local origin", method: http.MethodGet, origin: "http://localhost:5173", wantStatus: http.StatusOK, wantAllowOrigin: "http://localhost:5173", wantNext: true},
		{name: "configured origin", method: http.MethodGet, origin: "https://app.example.com", wantStatus: http.StatusOK, wantAllowOrigin: "https://app.example.com", wantNext: true},
		{name: "unconfigured origin", method: http.MethodGet, origin: "https://other.example.com", wantStatus: http.StatusForbidden},
		{name: "configured preflight", method: http.MethodOptions, origin: "https://app.example.com", wantStatus: http.StatusNoContent, wantAllowOrigin: "https://app.example.com"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nextCalled := false
			handler := localFrontendCORS(allowedOrigins)(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
				nextCalled = true
			}))
			request := httptest.NewRequest(tt.method, "/", nil)
			request.Header.Set("Origin", tt.origin)
			response := httptest.NewRecorder()

			handler.ServeHTTP(response, request)

			if response.Code != tt.wantStatus {
				t.Fatalf("status = %d, want %d", response.Code, tt.wantStatus)
			}
			if got := response.Header().Get("Access-Control-Allow-Origin"); got != tt.wantAllowOrigin {
				t.Fatalf("Access-Control-Allow-Origin = %q, want %q", got, tt.wantAllowOrigin)
			}
			if nextCalled != tt.wantNext {
				t.Fatalf("next called = %t, want %t", nextCalled, tt.wantNext)
			}
			if tt.method == http.MethodOptions && tt.wantAllowOrigin != "" &&
				response.Header().Get("Access-Control-Allow-Credentials") != "true" {
				t.Fatal("preflight response does not allow credentials")
			}
		})
	}
}

func TestParseCORSAllowedOrigins(t *testing.T) {
	origins, err := parseCORSAllowedOrigins(" https://app.example.com, http://localhost:4173 ")
	if err != nil {
		t.Fatal(err)
	}
	if len(origins) != 2 {
		t.Fatalf("got %d origins, want 2", len(origins))
	}

	for _, origin := range []string{"https://app.example.com/", "https://user@app.example.com", "ftp://app.example.com"} {
		if _, err := parseCORSAllowedOrigins(origin); err == nil {
			t.Errorf("parseCORSAllowedOrigins(%q) succeeded, want error", origin)
		}
	}
}
