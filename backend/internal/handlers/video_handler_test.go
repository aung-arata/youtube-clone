package handlers

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/aung-arata/youtube-clone/backend/internal/models"
	"github.com/gorilla/mux"
)

func TestGetVideos(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock database: %v", err)
	}
	defer db.Close()

	handler := NewVideoHandler(db)

	now := time.Now()
	rows := sqlmock.NewRows([]string{
		"id", "title", "description", "url", "thumbnail",
		"channel_name", "channel_avatar", "views", "duration",
		"uploaded_at", "created_at", "updated_at",
	}).AddRow(
		1, "Test Video", "Test Description", "http://example.com/video.mp4",
		"http://example.com/thumb.jpg", "Test Channel", "http://example.com/avatar.jpg",
		100, "10:00", now, now, now,
	)

	mock.ExpectQuery("SELECT (.+) FROM videos ORDER BY uploaded_at DESC LIMIT (.+) OFFSET (.+)").
		WithArgs(20, 0).
		WillReturnRows(rows)

	req, err := http.NewRequest("GET", "/api/videos", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler.GetVideos(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v, body: %s", status, http.StatusOK, rr.Body.String())
		return
	}

	var videos []models.Video
	if err := json.NewDecoder(rr.Body).Decode(&videos); err != nil {
		t.Errorf("Failed to decode response: %v", err)
	}

	if len(videos) != 1 {
		t.Errorf("Expected 1 video, got %d", len(videos))
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestCreateVideo_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock database: %v", err)
	}
	defer db.Close()

	handler := NewVideoHandler(db)

	video := models.Video{
		Title:         "New Video",
		Description:   "New Description",
		URL:           "http://example.com/video.mp4",
		Thumbnail:     "http://example.com/thumb.jpg",
		ChannelName:   "New Channel",
		ChannelAvatar: "http://example.com/avatar.jpg",
		Duration:      "5:00",
	}

	now := time.Now()
	rows := sqlmock.NewRows([]string{"id", "uploaded_at", "created_at", "updated_at"}).
		AddRow(1, now, now, now)

	mock.ExpectQuery("INSERT INTO videos (.+) VALUES (.+) RETURNING (.+)").
		WithArgs(video.Title, video.Description, video.URL, video.Thumbnail,
			video.ChannelName, video.ChannelAvatar, video.Duration).
		WillReturnRows(rows)

	body, _ := json.Marshal(video)
	req, err := http.NewRequest("POST", "/api/videos", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler.CreateVideo(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusCreated)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestCreateVideo_MissingTitle(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock database: %v", err)
	}
	defer db.Close()

	handler := NewVideoHandler(db)

	video := models.Video{
		URL:         "http://example.com/video.mp4",
		ChannelName: "Channel",
	}

	body, _ := json.Marshal(video)
	req, err := http.NewRequest("POST", "/api/videos", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler.CreateVideo(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}
}

func TestGetVideo_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock database: %v", err)
	}
	defer db.Close()

	handler := NewVideoHandler(db)

	now := time.Now()
	rows := sqlmock.NewRows([]string{
		"id", "title", "description", "url", "thumbnail",
		"channel_name", "channel_avatar", "views", "duration",
		"uploaded_at", "created_at", "updated_at",
	}).AddRow(
		1, "Test Video", "Test Description", "http://example.com/video.mp4",
		"http://example.com/thumb.jpg", "Test Channel", "http://example.com/avatar.jpg",
		100, "10:00", now, now, now,
	)

	mock.ExpectQuery("SELECT (.+) FROM videos WHERE id = (.+)").WithArgs(1).WillReturnRows(rows)

	req, err := http.NewRequest("GET", "/api/videos/1", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/api/videos/{id}", handler.GetVideo)
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestGetVideo_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock database: %v", err)
	}
	defer db.Close()

	handler := NewVideoHandler(db)

	mock.ExpectQuery("SELECT (.+) FROM videos WHERE id = (.+)").
		WithArgs(999).
		WillReturnError(sql.ErrNoRows)

	req, err := http.NewRequest("GET", "/api/videos/999", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/api/videos/{id}", handler.GetVideo)
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusNotFound)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestIncrementViews_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock database: %v", err)
	}
	defer db.Close()

	handler := NewVideoHandler(db)

	rows := sqlmock.NewRows([]string{"views"}).AddRow(101)

	mock.ExpectQuery("UPDATE videos SET views").WithArgs(1).WillReturnRows(rows)

	req, err := http.NewRequest("POST", "/api/videos/1/views", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/api/videos/{id}/views", handler.IncrementViews)
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}
