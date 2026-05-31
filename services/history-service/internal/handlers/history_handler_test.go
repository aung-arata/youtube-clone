package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

func newTestRouter(method, path string, handler http.HandlerFunc) *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc(path, handler).Methods(method)
	return r
}

func newHistoryHandler() *HistoryHandler {
	return &HistoryHandler{db: nil, httpClient: nil}
}

func TestAddToHistory_InvalidUserID(t *testing.T) {
	h := newHistoryHandler()
	router := newTestRouter(http.MethodPost, "/users/{userId}/history", h.AddToHistory)

	body, _ := json.Marshal(map[string]interface{}{"video_id": 1})
	req := httptest.NewRequest(http.MethodPost, "/users/notanumber/history", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for invalid userId, got %d", rr.Code)
	}
}

func TestAddToHistory_BadJSON(t *testing.T) {
	h := newHistoryHandler()
	router := newTestRouter(http.MethodPost, "/users/{userId}/history", h.AddToHistory)

	req := httptest.NewRequest(http.MethodPost, "/users/1/history", bytes.NewBufferString("{invalid"))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for bad JSON, got %d", rr.Code)
	}
}

func TestAddToHistory_ZeroVideoID(t *testing.T) {
	h := newHistoryHandler()
	router := newTestRouter(http.MethodPost, "/users/{userId}/history", h.AddToHistory)

	body, _ := json.Marshal(map[string]interface{}{"video_id": 0})
	req := httptest.NewRequest(http.MethodPost, "/users/1/history", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for video_id=0, got %d", rr.Code)
	}
}

func TestGetHistory_InvalidUserID(t *testing.T) {
	h := newHistoryHandler()
	router := newTestRouter(http.MethodGet, "/users/{userId}/history", h.GetHistory)

	req := httptest.NewRequest(http.MethodGet, "/users/notanumber/history", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for invalid userId, got %d", rr.Code)
	}
}
