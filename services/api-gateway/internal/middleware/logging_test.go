package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLoggingMiddleware_PassesThrough(t *testing.T) {
	called := false
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusCreated)
	})

	handler := LoggingMiddleware(next)
	req := httptest.NewRequest(http.MethodPost, "/api/test", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if !called {
		t.Error("expected next handler to be called")
	}
	if rr.Code != http.StatusCreated {
		t.Errorf("expected 201, got %d", rr.Code)
	}
}

func TestLoggingMiddleware_CapturesStatusCode(t *testing.T) {
	tests := []struct {
		name   string
		status int
	}{
		{"200 OK", http.StatusOK},
		{"404 Not Found", http.StatusNotFound},
		{"500 Internal Server Error", http.StatusInternalServerError},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := LoggingMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.status)
			}))
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)
			if rr.Code != tt.status {
				t.Errorf("expected %d, got %d", tt.status, rr.Code)
			}
		})
	}
}

func TestLoggingResponseWriter_DefaultStatus(t *testing.T) {
	// If WriteHeader is never called, default should be 200
	handler := LoggingMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Write body without calling WriteHeader explicitly
		w.Write([]byte("hello"))
	}))
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Errorf("expected default 200, got %d", rr.Code)
	}
}
