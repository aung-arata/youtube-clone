package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/aung-arata/youtube-clone/services/api-gateway/internal/auth"
)

func okHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func TestAuthMiddleware_MissingHeader(t *testing.T) {
	handler := AuthMiddleware(http.HandlerFunc(okHandler))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Errorf("expected 401 Unauthorized, got %d", rr.Code)
	}
}

func TestAuthMiddleware_InvalidFormat(t *testing.T) {
	handler := AuthMiddleware(http.HandlerFunc(okHandler))

	tests := []struct {
		name   string
		header string
	}{
		{"no bearer prefix", "sometoken"},
		{"token only no space", "Bearer"},
		{"basic auth", "Basic dXNlcjpwYXNz"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			req.Header.Set("Authorization", tt.header)
			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)
			if rr.Code != http.StatusUnauthorized {
				t.Errorf("expected 401, got %d", rr.Code)
			}
		})
	}
}

func TestAuthMiddleware_InvalidToken(t *testing.T) {
	handler := AuthMiddleware(http.HandlerFunc(okHandler))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer invalid.token.here")
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Errorf("expected 401 Unauthorized, got %d", rr.Code)
	}
}

func TestAuthMiddleware_ValidToken(t *testing.T) {
	token, err := auth.GenerateToken(99, "gopher", "user")
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}

	var capturedUserID interface{}
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedUserID = r.Context().Value(UserIDKey)
		w.WriteHeader(http.StatusOK)
	})

	handler := AuthMiddleware(next)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rr.Code)
	}
	if capturedUserID != 99 {
		t.Errorf("UserID in context: got %v, want 99", capturedUserID)
	}
}

func TestOptionalAuthMiddleware_NoHeader_PassesThrough(t *testing.T) {
	called := false
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
	})

	handler := OptionalAuthMiddleware(next)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if !called {
		t.Error("expected next handler to be called")
	}
	if rr.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rr.Code)
	}
}

func TestOptionalAuthMiddleware_ValidToken_SetsContext(t *testing.T) {
	token, err := auth.GenerateToken(55, "optuser", "admin")
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}

	var capturedRole interface{}
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedRole = r.Context().Value(UserRoleKey)
		w.WriteHeader(http.StatusOK)
	})

	handler := OptionalAuthMiddleware(next)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rr.Code)
	}
	if capturedRole != "admin" {
		t.Errorf("Role in context: got %v, want admin", capturedRole)
	}
}

func TestOptionalAuthMiddleware_InvalidToken_PassesThrough(t *testing.T) {
	called := false
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
	})

	handler := OptionalAuthMiddleware(next)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer bad.token.here")
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if !called {
		t.Error("expected next handler to be called even with invalid token")
	}
}
