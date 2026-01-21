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

func TestGetUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock: %v", err)
	}
	defer db.Close()

	handler := NewUserHandler(db)

	now := time.Now()
	rows := sqlmock.NewRows([]string{"id", "username", "email", "avatar", "plan_id", "created_at", "updated_at"}).
		AddRow(1, "testuser", "test@example.com", "https://example.com/avatar.jpg", nil, now, now)

	mock.ExpectQuery("SELECT (.+) FROM users WHERE id").
		WithArgs(1).
		WillReturnRows(rows)

	req := httptest.NewRequest("GET", "/api/users/1", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "1"})
	w := httptest.NewRecorder()

	handler.GetUser(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var user models.User
	json.NewDecoder(w.Body).Decode(&user)

	if user.Username != "testuser" {
		t.Errorf("Expected username testuser, got %s", user.Username)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestCreateUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock: %v", err)
	}
	defer db.Close()

	handler := NewUserHandler(db)

	now := time.Now()
	rows := sqlmock.NewRows([]string{"id", "username", "email", "avatar", "plan_id", "created_at", "updated_at"}).
		AddRow(1, "newuser", "new@example.com", "", nil, now, now)

	mock.ExpectQuery("INSERT INTO users").
		WithArgs("newuser", "new@example.com", "").
		WillReturnRows(rows)

	user := models.User{
		Username: "newuser",
		Email:    "new@example.com",
	}

	body, _ := json.Marshal(user)
	req := httptest.NewRequest("POST", "/api/users", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	handler.CreateUser(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status 201, got %d", w.Code)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestUpdateUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock: %v", err)
	}
	defer db.Close()

	handler := NewUserHandler(db)

	now := time.Now()
	rows := sqlmock.NewRows([]string{"id", "username", "email", "avatar", "plan_id", "created_at", "updated_at"}).
		AddRow(1, "updateduser", "updated@example.com", "https://example.com/new.jpg", nil, now, now)

	mock.ExpectQuery("UPDATE users").
		WithArgs("updateduser", "updated@example.com", "https://example.com/new.jpg", 1).
		WillReturnRows(rows)

	user := models.User{
		Username: "updateduser",
		Email:    "updated@example.com",
		Avatar:   "https://example.com/new.jpg",
	}

	body, _ := json.Marshal(user)
	req := httptest.NewRequest("PUT", "/api/users/1", bytes.NewBuffer(body))
	req = mux.SetURLVars(req, map[string]string{"id": "1"})
	w := httptest.NewRecorder()

	handler.UpdateUser(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}
