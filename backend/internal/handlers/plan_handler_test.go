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

func TestGetPlans(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock: %v", err)
	}
	defer db.Close()

	handler := NewPlanHandler(db)

	now := time.Now()
	rows := sqlmock.NewRows([]string{"id", "name", "price", "max_video_quality", "max_uploads_per_month", "ads_free", "created_at", "updated_at"}).
		AddRow(1, "Free", 0.0, "480p", 5, false, now, now).
		AddRow(2, "Premium", 9.99, "1080p", 100, true, now, now)

	mock.ExpectQuery("SELECT (.+) FROM plans ORDER BY price").
		WillReturnRows(rows)

	req := httptest.NewRequest("GET", "/api/plans", nil)
	w := httptest.NewRecorder()

	handler.GetPlans(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var plans []models.Plan
	json.NewDecoder(w.Body).Decode(&plans)

	if len(plans) != 2 {
		t.Errorf("Expected 2 plans, got %d", len(plans))
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestGetPlan(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock: %v", err)
	}
	defer db.Close()

	handler := NewPlanHandler(db)

	now := time.Now()
	rows := sqlmock.NewRows([]string{"id", "name", "price", "max_video_quality", "max_uploads_per_month", "ads_free", "created_at", "updated_at"}).
		AddRow(1, "Premium", 9.99, "1080p", 100, true, now, now)

	mock.ExpectQuery("SELECT (.+) FROM plans WHERE id").
		WithArgs(1).
		WillReturnRows(rows)

	req := httptest.NewRequest("GET", "/api/plans/1", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "1"})
	w := httptest.NewRecorder()

	handler.GetPlan(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var plan models.Plan
	json.NewDecoder(w.Body).Decode(&plan)

	if plan.Name != "Premium" {
		t.Errorf("Expected plan name Premium, got %s", plan.Name)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestCreatePlan(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock: %v", err)
	}
	defer db.Close()

	handler := NewPlanHandler(db)

	now := time.Now()
	rows := sqlmock.NewRows([]string{"id", "name", "price", "max_video_quality", "max_uploads_per_month", "ads_free", "created_at", "updated_at"}).
		AddRow(1, "Enterprise", 19.99, "4K", -1, true, now, now)

	mock.ExpectQuery("INSERT INTO plans").
		WithArgs("Enterprise", 19.99, "4K", -1, true).
		WillReturnRows(rows)

	plan := models.Plan{
		Name:               "Enterprise",
		Price:              19.99,
		MaxVideoQuality:    "4K",
		MaxUploadsPerMonth: -1,
		AdsFree:            true,
	}

	body, _ := json.Marshal(plan)
	req := httptest.NewRequest("POST", "/api/plans", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	handler.CreatePlan(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status 201, got %d", w.Code)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestUpdatePlan(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock: %v", err)
	}
	defer db.Close()

	handler := NewPlanHandler(db)

	now := time.Now()
	rows := sqlmock.NewRows([]string{"id", "name", "price", "max_video_quality", "max_uploads_per_month", "ads_free", "created_at", "updated_at"}).
		AddRow(1, "Premium Plus", 14.99, "4K", 200, true, now, now)

	mock.ExpectQuery("UPDATE plans").
		WithArgs("Premium Plus", 14.99, "4K", 200, true, 1).
		WillReturnRows(rows)

	plan := models.Plan{
		Name:               "Premium Plus",
		Price:              14.99,
		MaxVideoQuality:    "4K",
		MaxUploadsPerMonth: 200,
		AdsFree:            true,
	}

	body, _ := json.Marshal(plan)
	req := httptest.NewRequest("PUT", "/api/plans/1", bytes.NewBuffer(body))
	req = mux.SetURLVars(req, map[string]string{"id": "1"})
	w := httptest.NewRecorder()

	handler.UpdatePlan(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestDeletePlan(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock: %v", err)
	}
	defer db.Close()

	handler := NewPlanHandler(db)

	mock.ExpectExec("DELETE FROM plans WHERE id").
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(0, 1))

	req := httptest.NewRequest("DELETE", "/api/plans/1", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "1"})
	w := httptest.NewRecorder()

	handler.DeletePlan(w, req)

	if w.Code != http.StatusNoContent {
		t.Errorf("Expected status 204, got %d", w.Code)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestGetUserPlan(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock: %v", err)
	}
	defer db.Close()

	handler := NewPlanHandler(db)

	now := time.Now()
	rows := sqlmock.NewRows([]string{"id", "name", "price", "max_video_quality", "max_uploads_per_month", "ads_free", "created_at", "updated_at"}).
		AddRow(2, "Premium", 9.99, "1080p", 100, true, now, now)

	mock.ExpectQuery("SELECT (.+) FROM plans p INNER JOIN users u").
		WithArgs(1).
		WillReturnRows(rows)

	req := httptest.NewRequest("GET", "/api/users/1/plan", nil)
	req = mux.SetURLVars(req, map[string]string{"userId": "1"})
	w := httptest.NewRecorder()

	handler.GetUserPlan(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var plan models.Plan
	json.NewDecoder(w.Body).Decode(&plan)

	if plan.Name != "Premium" {
		t.Errorf("Expected plan name Premium, got %s", plan.Name)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestUpdateUserPlan(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock: %v", err)
	}
	defer db.Close()

	handler := NewPlanHandler(db)

	// Mock plan existence check
	existsRows := sqlmock.NewRows([]string{"exists"}).AddRow(true)
	mock.ExpectQuery("SELECT EXISTS").
		WithArgs(2).
		WillReturnRows(existsRows)

	// Mock user update
	updateRows := sqlmock.NewRows([]string{"id"}).AddRow(1)
	mock.ExpectQuery("UPDATE users SET plan_id").
		WithArgs(2, 1).
		WillReturnRows(updateRows)

	reqBody := map[string]int{"plan_id": 2}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("PUT", "/api/users/1/plan", bytes.NewBuffer(body))
	req = mux.SetURLVars(req, map[string]string{"userId": "1"})
	w := httptest.NewRecorder()

	handler.UpdateUserPlan(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}
