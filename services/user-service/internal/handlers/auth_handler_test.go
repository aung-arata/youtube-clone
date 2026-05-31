package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func newAuthHandler() *AuthHandler {
	return &AuthHandler{db: nil}
}

func TestSignup_BadJSON(t *testing.T) {
	h := newAuthHandler()
	req := httptest.NewRequest(http.MethodPost, "/auth/signup", bytes.NewBufferString("{invalid"))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	h.Signup(rr, req)
	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for bad JSON, got %d", rr.Code)
	}
}

func TestSignup_EmptyUsername(t *testing.T) {
	h := newAuthHandler()
	body, _ := json.Marshal(SignupRequest{Username: "", Email: "a@b.com", Password: "password"})
	req := httptest.NewRequest(http.MethodPost, "/auth/signup", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	h.Signup(rr, req)
	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for empty username, got %d", rr.Code)
	}
}

func TestSignup_EmptyEmail(t *testing.T) {
	h := newAuthHandler()
	body, _ := json.Marshal(SignupRequest{Username: "alice", Email: "", Password: "password"})
	req := httptest.NewRequest(http.MethodPost, "/auth/signup", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	h.Signup(rr, req)
	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for empty email, got %d", rr.Code)
	}
}

func TestSignup_EmptyPassword(t *testing.T) {
	h := newAuthHandler()
	body, _ := json.Marshal(SignupRequest{Username: "alice", Email: "a@b.com", Password: ""})
	req := httptest.NewRequest(http.MethodPost, "/auth/signup", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	h.Signup(rr, req)
	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for empty password, got %d", rr.Code)
	}
}

func TestSignup_ShortPassword(t *testing.T) {
	h := newAuthHandler()
	body, _ := json.Marshal(SignupRequest{Username: "alice", Email: "a@b.com", Password: "abc"})
	req := httptest.NewRequest(http.MethodPost, "/auth/signup", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	h.Signup(rr, req)
	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for short password (< 6 chars), got %d", rr.Code)
	}
}

func TestLogin_BadJSON(t *testing.T) {
	h := newAuthHandler()
	req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewBufferString("{invalid"))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	h.Login(rr, req)
	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for bad JSON, got %d", rr.Code)
	}
}

func TestLogin_EmptyEmail(t *testing.T) {
	h := newAuthHandler()
	body, _ := json.Marshal(LoginRequest{Email: "", Password: "password"})
	req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	h.Login(rr, req)
	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for empty email, got %d", rr.Code)
	}
}

func TestLogin_EmptyPassword(t *testing.T) {
	h := newAuthHandler()
	body, _ := json.Marshal(LoginRequest{Email: "a@b.com", Password: ""})
	req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	h.Login(rr, req)
	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for empty password, got %d", rr.Code)
	}
}
