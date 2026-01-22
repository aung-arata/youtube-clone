package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/aung-arata/youtube-clone/backend/internal/auth"
	"github.com/aung-arata/youtube-clone/backend/internal/middleware"
	"github.com/aung-arata/youtube-clone/backend/internal/models"
)

type AuthHandler struct {
	db *sql.DB
}

func NewAuthHandler(db *sql.DB) *AuthHandler {
	return &AuthHandler{db: db}
}

// SignupRequest represents the signup request body
type SignupRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Avatar   string `json:"avatar"`
}

// LoginRequest represents the login request body
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// AuthResponse represents the authentication response
type AuthResponse struct {
	Token string       `json:"token"`
	User  models.User `json:"user"`
}

// Signup handles user registration
func (h *AuthHandler) Signup(w http.ResponseWriter, r *http.Request) {
	var req SignupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if strings.TrimSpace(req.Username) == "" {
		http.Error(w, "Username is required", http.StatusBadRequest)
		return
	}
	if strings.TrimSpace(req.Email) == "" {
		http.Error(w, "Email is required", http.StatusBadRequest)
		return
	}
	if strings.TrimSpace(req.Password) == "" {
		http.Error(w, "Password is required", http.StatusBadRequest)
		return
	}
	if len(req.Password) < 6 {
		http.Error(w, "Password must be at least 6 characters", http.StatusBadRequest)
		return
	}

	// Check if user already exists
	var exists bool
	checkQuery := `SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)`
	err := h.db.QueryRow(checkQuery, req.Email).Scan(&exists)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	if exists {
		http.Error(w, "Email already registered", http.StatusConflict)
		return
	}

	// Hash password
	hashedPassword, err := auth.HashPassword(req.Password)
	if err != nil {
		http.Error(w, "Error processing password", http.StatusInternalServerError)
		return
	}

	// Create user with default role "user"
	query := `
		INSERT INTO users (username, email, password, avatar, role)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, username, email, avatar, role, plan_id, created_at, updated_at
	`

	var user models.User
	err = h.db.QueryRow(query, req.Username, req.Email, hashedPassword, req.Avatar, "user").Scan(
		&user.ID, &user.Username, &user.Email, &user.Avatar, &user.Role, &user.PlanID, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		http.Error(w, "Error creating user", http.StatusInternalServerError)
		return
	}

	// Generate JWT token
	token, err := auth.GenerateToken(user.ID, user.Username, user.Role)
	if err != nil {
		http.Error(w, "Error generating token", http.StatusInternalServerError)
		return
	}

	// Return response
	response := AuthResponse{
		Token: token,
		User:  user,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// Login handles user authentication
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if strings.TrimSpace(req.Email) == "" {
		http.Error(w, "Email is required", http.StatusBadRequest)
		return
	}
	if strings.TrimSpace(req.Password) == "" {
		http.Error(w, "Password is required", http.StatusBadRequest)
		return
	}

	// Get user by email
	query := `
		SELECT id, username, email, password, avatar, role, plan_id, created_at, updated_at
		FROM users
		WHERE email = $1
	`

	var user models.User
	var hashedPassword string
	err := h.db.QueryRow(query, req.Email).Scan(
		&user.ID, &user.Username, &user.Email, &hashedPassword, &user.Avatar, &user.Role, &user.PlanID, &user.CreatedAt, &user.UpdatedAt)

	if err == sql.ErrNoRows {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	} else if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	// Compare passwords
	if err := auth.ComparePasswords(hashedPassword, req.Password); err != nil {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	// Generate JWT token
	token, err := auth.GenerateToken(user.ID, user.Username, user.Role)
	if err != nil {
		http.Error(w, "Error generating token", http.StatusInternalServerError)
		return
	}

	// Return response
	response := AuthResponse{
		Token: token,
		User:  user,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// RefreshToken handles token refresh
func (h *AuthHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	// Get token from Authorization header
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, "Authorization header required", http.StatusUnauthorized)
		return
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		http.Error(w, "Invalid authorization header format", http.StatusUnauthorized)
		return
	}

	// Refresh token
	newToken, err := auth.RefreshToken(parts[1])
	if err != nil {
		http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
		return
	}

	response := map[string]string{
		"token": newToken,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetCurrentUser returns the currently authenticated user
func (h *AuthHandler) GetCurrentUser(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context (set by AuthMiddleware)
	userID, ok := r.Context().Value(middleware.UserIDKey).(int)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Get user from database
	query := `
		SELECT id, username, email, avatar, role, plan_id, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	var user models.User
	err := h.db.QueryRow(query, userID).Scan(
		&user.ID, &user.Username, &user.Email, &user.Avatar, &user.Role, &user.PlanID, &user.CreatedAt, &user.UpdatedAt)

	if err == sql.ErrNoRows {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}
