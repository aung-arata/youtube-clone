package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/aung-arata/youtube-clone/services/video-service/internal/auth"
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
		t.Errorf("expected 401, got %d", rr.Code)
	}
}

func TestAuthMiddleware_InvalidFormat(t *testing.T) {
	handler := AuthMiddleware(http.HandlerFunc(okHandler))
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Token abc123")
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", rr.Code)
	}
}

func TestAuthMiddleware_InvalidToken(t *testing.T) {
	handler := AuthMiddleware(http.HandlerFunc(okHandler))
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer notavalidtoken")
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", rr.Code)
	}
}

func TestAuthMiddleware_ValidToken_ContextValues(t *testing.T) {
	token, err := auth.GenerateToken(42, "alice", "user")
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}

	var capturedUserID int
	var capturedUsername, capturedRole string
	inner := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedUserID = r.Context().Value(UserIDKey).(int)
		capturedUsername = r.Context().Value(UsernameKey).(string)
		capturedRole = r.Context().Value(UserRoleKey).(string)
		w.WriteHeader(http.StatusOK)
	})

	handler := AuthMiddleware(inner)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
	if capturedUserID != 42 {
		t.Errorf("expected userID 42, got %d", capturedUserID)
	}
	if capturedUsername != "alice" {
		t.Errorf("expected username alice, got %q", capturedUsername)
	}
	if capturedRole != "user" {
		t.Errorf("expected role user, got %q", capturedRole)
	}
}

func TestOptionalAuthMiddleware_NoHeader(t *testing.T) {
	handler := OptionalAuthMiddleware(http.HandlerFunc(okHandler))
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Errorf("expected 200 (pass-through), got %d", rr.Code)
	}
}

func TestOptionalAuthMiddleware_InvalidToken_PassesThrough(t *testing.T) {
	handler := OptionalAuthMiddleware(http.HandlerFunc(okHandler))
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer invalid.token.here")
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Errorf("expected 200 (pass-through on bad token), got %d", rr.Code)
	}
}

func TestOptionalAuthMiddleware_ValidToken_SetsContext(t *testing.T) {
	token, err := auth.GenerateToken(7, "bob", "admin")
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}

	var capturedID int
	inner := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedID = r.Context().Value(UserIDKey).(int)
		w.WriteHeader(http.StatusOK)
	})

	handler := OptionalAuthMiddleware(inner)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rr.Code)
	}
	if capturedID != 7 {
		t.Errorf("expected userID 7, got %d", capturedID)
	}
}

func TestAdminOnlyMiddleware_NoRole(t *testing.T) {
	handler := AdminOnlyMiddleware(http.HandlerFunc(okHandler))
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusForbidden {
		t.Errorf("expected 403, got %d", rr.Code)
	}
}

func TestAdminOnlyMiddleware_NonAdminRole(t *testing.T) {
	handler := AdminOnlyMiddleware(http.HandlerFunc(okHandler))
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	ctx := context.WithValue(req.Context(), UserRoleKey, "user")
	req = req.WithContext(ctx)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusForbidden {
		t.Errorf("expected 403 for non-admin, got %d", rr.Code)
	}
}

func TestAdminOnlyMiddleware_AdminRole(t *testing.T) {
	handler := AdminOnlyMiddleware(http.HandlerFunc(okHandler))
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	ctx := context.WithValue(req.Context(), UserRoleKey, "admin")
	req = req.WithContext(ctx)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Errorf("expected 200 for admin, got %d", rr.Code)
	}
}
