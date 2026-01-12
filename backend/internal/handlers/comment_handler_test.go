package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/aung-arata/youtube-clone/backend/internal/models"
	"github.com/gorilla/mux"
)

func TestGetComments(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock: %v", err)
	}
	defer db.Close()

	handler := NewCommentHandler(db)

	// Mock data
	now := time.Now()
	rows := sqlmock.NewRows([]string{"id", "video_id", "user_id", "content", "created_at", "updated_at"}).
		AddRow(1, 1, 1, "Great video!", now, now).
		AddRow(2, 1, 2, "Thanks for sharing", now, now)

	mock.ExpectQuery("SELECT (.+) FROM comments WHERE video_id").
		WithArgs(1).
		WillReturnRows(rows)

	req := httptest.NewRequest("GET", "/api/videos/1/comments", nil)
	req = mux.SetURLVars(req, map[string]string{"videoId": "1"})
	w := httptest.NewRecorder()

	handler.GetComments(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var comments []models.Comment
	json.NewDecoder(w.Body).Decode(&comments)

	if len(comments) != 2 {
		t.Errorf("Expected 2 comments, got %d", len(comments))
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestGetComment(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock: %v", err)
	}
	defer db.Close()

	handler := NewCommentHandler(db)

	// Mock data
	now := time.Now()
	rows := sqlmock.NewRows([]string{"id", "video_id", "user_id", "content", "created_at", "updated_at"}).
		AddRow(1, 1, 1, "Great video!", now, now)

	mock.ExpectQuery("SELECT (.+) FROM comments WHERE id").
		WithArgs(1).
		WillReturnRows(rows)

	req := httptest.NewRequest("GET", "/api/comments/1", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "1"})
	w := httptest.NewRecorder()

	handler.GetComment(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var comment models.Comment
	json.NewDecoder(w.Body).Decode(&comment)

	if comment.ID != 1 {
		t.Errorf("Expected comment ID 1, got %d", comment.ID)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestCreateComment(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock: %v", err)
	}
	defer db.Close()

	handler := NewCommentHandler(db)

	// Mock data
	now := time.Now()
	rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at"}).
		AddRow(1, now, now)

	mock.ExpectQuery("INSERT INTO comments").
		WithArgs(1, 1, "Great video!").
		WillReturnRows(rows)

	comment := models.Comment{
		UserID:  1,
		Content: "Great video!",
	}

	body, _ := json.Marshal(comment)
	req := httptest.NewRequest("POST", "/api/videos/1/comments", bytes.NewBuffer(body))
	req = mux.SetURLVars(req, map[string]string{"videoId": "1"})
	w := httptest.NewRecorder()

	handler.CreateComment(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status 201, got %d", w.Code)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestCreateComment_InvalidData(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock: %v", err)
	}
	defer db.Close()

	handler := NewCommentHandler(db)

	tests := []struct {
		name       string
		comment    models.Comment
		statusCode int
	}{
		{
			name:       "Empty content",
			comment:    models.Comment{UserID: 1, Content: ""},
			statusCode: http.StatusBadRequest,
		},
		{
			name:       "Missing user ID",
			comment:    models.Comment{UserID: 0, Content: "Great video!"},
			statusCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.comment)
			req := httptest.NewRequest("POST", "/api/videos/1/comments", bytes.NewBuffer(body))
			req = mux.SetURLVars(req, map[string]string{"videoId": "1"})
			w := httptest.NewRecorder()

			handler.CreateComment(w, req)

			if w.Code != tt.statusCode {
				t.Errorf("Expected status %d, got %d", tt.statusCode, w.Code)
			}
		})
	}
}

func TestUpdateComment(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock: %v", err)
	}
	defer db.Close()

	handler := NewCommentHandler(db)

	// Mock data
	now := time.Now()
	rows := sqlmock.NewRows([]string{"id", "video_id", "user_id", "content", "created_at", "updated_at"}).
		AddRow(1, 1, 1, "Updated comment!", now, now)

	mock.ExpectQuery("UPDATE comments").
		WithArgs("Updated comment!", 1).
		WillReturnRows(rows)

	comment := models.Comment{
		Content: "Updated comment!",
	}

	body, _ := json.Marshal(comment)
	req := httptest.NewRequest("PUT", "/api/comments/1", bytes.NewBuffer(body))
	req = mux.SetURLVars(req, map[string]string{"id": "1"})
	w := httptest.NewRecorder()

	handler.UpdateComment(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestDeleteComment(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock: %v", err)
	}
	defer db.Close()

	handler := NewCommentHandler(db)

	mock.ExpectExec("DELETE FROM comments WHERE id").
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(0, 1))

	req := httptest.NewRequest("DELETE", "/api/comments/1", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "1"})
	w := httptest.NewRecorder()

	handler.DeleteComment(w, req)

	if w.Code != http.StatusNoContent {
		t.Errorf("Expected status 204, got %d", w.Code)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestDeleteComment_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock: %v", err)
	}
	defer db.Close()

	handler := NewCommentHandler(db)

	mock.ExpectExec("DELETE FROM comments WHERE id").
		WithArgs(999).
		WillReturnResult(sqlmock.NewResult(0, 0))

	req := httptest.NewRequest("DELETE", "/api/comments/999", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "999"})
	w := httptest.NewRecorder()

	handler.DeleteComment(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", w.Code)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}
