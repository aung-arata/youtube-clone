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

type UserHandler struct {
	db *sql.DB
}

func NewUserHandler(db *sql.DB) *UserHandler {
	return &UserHandler{db: db}
}

// GetUser returns a single user by ID
func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	query := `
		SELECT id, username, email, avatar, plan_id, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	var u models.User
	err = h.db.QueryRow(query, id).Scan(&u.ID, &u.Username, &u.Email, &u.Avatar, &u.PlanID, &u.CreatedAt, &u.UpdatedAt)

	if err == sql.ErrNoRows {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(u)
}

// UpdateUser updates user profile information
func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	var u models.User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Validate required fields
	if strings.TrimSpace(u.Username) == "" {
		http.Error(w, "Username is required", http.StatusBadRequest)
		return
	}
	if strings.TrimSpace(u.Email) == "" {
		http.Error(w, "Email is required", http.StatusBadRequest)
		return
	}

	query := `
		UPDATE users 
		SET username = $1, email = $2, avatar = $3, updated_at = CURRENT_TIMESTAMP
		WHERE id = $4
		RETURNING id, username, email, avatar, plan_id, created_at, updated_at
	`

	err = h.db.QueryRow(query, u.Username, u.Email, u.Avatar, id).Scan(
		&u.ID, &u.Username, &u.Email, &u.Avatar, &u.PlanID, &u.CreatedAt, &u.UpdatedAt)

	if err == sql.ErrNoRows {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(u)
}

// CreateUser creates a new user
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var u models.User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Validate required fields
	if strings.TrimSpace(u.Username) == "" {
		http.Error(w, "Username is required", http.StatusBadRequest)
		return
	}
	if strings.TrimSpace(u.Email) == "" {
		http.Error(w, "Email is required", http.StatusBadRequest)
		return
	}

	query := `
		INSERT INTO users (username, email, avatar)
		VALUES ($1, $2, $3)
		RETURNING id, username, email, avatar, plan_id, created_at, updated_at
	`

	err := h.db.QueryRow(query, u.Username, u.Email, u.Avatar).Scan(
		&u.ID, &u.Username, &u.Email, &u.Avatar, &u.PlanID, &u.CreatedAt, &u.UpdatedAt)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(u)
}
