package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gorilla/mux"
)

func TestAddToHistory(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock: %v", err)
	}
	defer db.Close()

	handler := NewHistoryHandler(db)

	now := time.Now()
	rows := sqlmock.NewRows([]string{"id", "user_id", "video_id", "watched_at"}).
		AddRow(1, 1, 1, now)

	mock.ExpectQuery("INSERT INTO watch_history").
		WithArgs(1, 1).
		WillReturnRows(rows)

	reqBody := map[string]int{"video_id": 1}
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("POST", "/api/users/1/history", bytes.NewBuffer(body))
	req = mux.SetURLVars(req, map[string]string{"userId": "1"})
	w := httptest.NewRecorder()

	handler.AddToHistory(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status 201, got %d", w.Code)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestGetHistory(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock: %v", err)
	}
	defer db.Close()

	handler := NewHistoryHandler(db)

	now := time.Now()
	rows := sqlmock.NewRows([]string{
		"id", "title", "description", "url", "thumbnail",
		"channel_name", "channel_avatar", "views", "likes", "dislikes",
		"category", "duration", "uploaded_at", "created_at", "updated_at", "watched_at",
	}).AddRow(
		1, "Test Video", "Description", "http://example.com/video.mp4",
		"http://example.com/thumb.jpg", "Test Channel", "http://example.com/avatar.jpg",
		100, 5, 1, "Education", "10:00", now, now, now, now,
	)

	mock.ExpectQuery("SELECT (.+) FROM watch_history (.+) JOIN videos").
		WithArgs(1, 20, 0).
		WillReturnRows(rows)

	req := httptest.NewRequest("GET", "/api/users/1/history", nil)
	req = mux.SetURLVars(req, map[string]string{"userId": "1"})
	w := httptest.NewRecorder()

	handler.GetHistory(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}
