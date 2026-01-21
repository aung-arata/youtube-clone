package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/aung-arata/youtube-clone/services/user-service/internal/models"
	"github.com/gorilla/mux"
)

type PlanHandler struct {
	db *sql.DB
}

func NewPlanHandler(db *sql.DB) *PlanHandler {
	return &PlanHandler{db: db}
}

// GetPlans returns all available plans
func (h *PlanHandler) GetPlans(w http.ResponseWriter, r *http.Request) {
	query := `
		SELECT id, name, price, max_video_quality, max_uploads_per_month, ads_free, created_at, updated_at
		FROM plans
		ORDER BY price ASC
	`

	rows, err := h.db.Query(query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var plans []models.Plan
	for rows.Next() {
		var p models.Plan
		err := rows.Scan(&p.ID, &p.Name, &p.Price, &p.MaxVideoQuality, &p.MaxUploadsPerMonth, &p.AdsFree, &p.CreatedAt, &p.UpdatedAt)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		plans = append(plans, p)
	}

	if err = rows.Err(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(plans)
}

// GetPlan returns a single plan by ID
func (h *PlanHandler) GetPlan(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid plan ID", http.StatusBadRequest)
		return
	}

	query := `
		SELECT id, name, price, max_video_quality, max_uploads_per_month, ads_free, created_at, updated_at
		FROM plans
		WHERE id = $1
	`

	var p models.Plan
	err = h.db.QueryRow(query, id).Scan(&p.ID, &p.Name, &p.Price, &p.MaxVideoQuality, &p.MaxUploadsPerMonth, &p.AdsFree, &p.CreatedAt, &p.UpdatedAt)

	if err == sql.ErrNoRows {
		http.Error(w, "Plan not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(p)
}

// CreatePlan creates a new subscription plan
func (h *PlanHandler) CreatePlan(w http.ResponseWriter, r *http.Request) {
	var p models.Plan
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Validate required fields
	if strings.TrimSpace(p.Name) == "" {
		http.Error(w, "Plan name is required", http.StatusBadRequest)
		return
	}
	if p.Price < 0 {
		http.Error(w, "Price must be non-negative", http.StatusBadRequest)
		return
	}

	query := `
		INSERT INTO plans (name, price, max_video_quality, max_uploads_per_month, ads_free)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, name, price, max_video_quality, max_uploads_per_month, ads_free, created_at, updated_at
	`

	err := h.db.QueryRow(query, p.Name, p.Price, p.MaxVideoQuality, p.MaxUploadsPerMonth, p.AdsFree).Scan(
		&p.ID, &p.Name, &p.Price, &p.MaxVideoQuality, &p.MaxUploadsPerMonth, &p.AdsFree, &p.CreatedAt, &p.UpdatedAt)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(p)
}

// UpdatePlan updates an existing plan
func (h *PlanHandler) UpdatePlan(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid plan ID", http.StatusBadRequest)
		return
	}

	var p models.Plan
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Validate required fields
	if strings.TrimSpace(p.Name) == "" {
		http.Error(w, "Plan name is required", http.StatusBadRequest)
		return
	}
	if p.Price < 0 {
		http.Error(w, "Price must be non-negative", http.StatusBadRequest)
		return
	}

	query := `
		UPDATE plans 
		SET name = $1, price = $2, max_video_quality = $3, max_uploads_per_month = $4, ads_free = $5, updated_at = CURRENT_TIMESTAMP
		WHERE id = $6
		RETURNING id, name, price, max_video_quality, max_uploads_per_month, ads_free, created_at, updated_at
	`

	err = h.db.QueryRow(query, p.Name, p.Price, p.MaxVideoQuality, p.MaxUploadsPerMonth, p.AdsFree, id).Scan(
		&p.ID, &p.Name, &p.Price, &p.MaxVideoQuality, &p.MaxUploadsPerMonth, &p.AdsFree, &p.CreatedAt, &p.UpdatedAt)

	if err == sql.ErrNoRows {
		http.Error(w, "Plan not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(p)
}

// DeletePlan deletes a plan by ID
func (h *PlanHandler) DeletePlan(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid plan ID", http.StatusBadRequest)
		return
	}

	query := `DELETE FROM plans WHERE id = $1`

	result, err := h.db.Exec(query, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if rowsAffected == 0 {
		http.Error(w, "Plan not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// GetUserPlan returns the plan for a specific user
func (h *PlanHandler) GetUserPlan(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, err := strconv.Atoi(vars["userId"])
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	query := `
		SELECT p.id, p.name, p.price, p.max_video_quality, p.max_uploads_per_month, p.ads_free, p.created_at, p.updated_at
		FROM plans p
		INNER JOIN users u ON u.plan_id = p.id
		WHERE u.id = $1
	`

	var p models.Plan
	err = h.db.QueryRow(query, userID).Scan(&p.ID, &p.Name, &p.Price, &p.MaxVideoQuality, &p.MaxUploadsPerMonth, &p.AdsFree, &p.CreatedAt, &p.UpdatedAt)

	if err == sql.ErrNoRows {
		http.Error(w, "User has no plan assigned", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(p)
}

// UpdateUserPlan updates the plan for a specific user
func (h *PlanHandler) UpdateUserPlan(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, err := strconv.Atoi(vars["userId"])
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	var req struct {
		PlanID int `json:"plan_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Verify plan exists
	var planExists bool
	err = h.db.QueryRow("SELECT EXISTS(SELECT 1 FROM plans WHERE id = $1)", req.PlanID).Scan(&planExists)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if !planExists {
		http.Error(w, "Plan not found", http.StatusNotFound)
		return
	}

	// Update user's plan
	query := `
		UPDATE users 
		SET plan_id = $1, updated_at = CURRENT_TIMESTAMP
		WHERE id = $2
		RETURNING id
	`

	var updatedUserID int
	err = h.db.QueryRow(query, req.PlanID, userID).Scan(&updatedUserID)

	if err == sql.ErrNoRows {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "User plan updated successfully",
		"user_id": updatedUserID,
		"plan_id": req.PlanID,
	})
}
