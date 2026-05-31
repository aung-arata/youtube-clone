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

func newNotificationHandler() *NotificationHandler {
	return &NotificationHandler{db: nil}
}

func TestGetUserNotifications_InvalidUserID(t *testing.T) {
	h := newNotificationHandler()
	router := newTestRouter(http.MethodGet, "/users/{userId}/notifications", h.GetUserNotifications)

	req := httptest.NewRequest(http.MethodGet, "/users/notanumber/notifications", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for invalid userId, got %d", rr.Code)
	}
}

func TestCreateNotification_BadJSON(t *testing.T) {
	h := newNotificationHandler()
	router := newTestRouter(http.MethodPost, "/notifications", h.CreateNotification)

	req := httptest.NewRequest(http.MethodPost, "/notifications", bytes.NewBufferString("{invalid"))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for bad JSON, got %d", rr.Code)
	}
}

func TestCreateNotification_MissingFields(t *testing.T) {
	h := newNotificationHandler()
	router := newTestRouter(http.MethodPost, "/notifications", h.CreateNotification)

	// Missing type, title, message
	body, _ := json.Marshal(map[string]interface{}{"user_id": 1})
	req := httptest.NewRequest(http.MethodPost, "/notifications", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for missing required fields, got %d", rr.Code)
	}
}

func TestCreateNotification_ZeroUserID(t *testing.T) {
	h := newNotificationHandler()
	router := newTestRouter(http.MethodPost, "/notifications", h.CreateNotification)

	body, _ := json.Marshal(map[string]interface{}{
		"user_id": 0,
		"type":    "info",
		"title":   "Test",
		"message": "Hello",
	})
	req := httptest.NewRequest(http.MethodPost, "/notifications", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for zero user_id, got %d", rr.Code)
	}
}
