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

type SubscriptionHandler struct {
	db *sql.DB
}

func NewSubscriptionHandler(db *sql.DB) *SubscriptionHandler {
	return &SubscriptionHandler{db: db}
}

// Subscribe allows a user to subscribe to a channel
func (h *SubscriptionHandler) Subscribe(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, err := strconv.Atoi(vars["userId"])
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	var req struct {
		ChannelName string `json:"channel_name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if req.ChannelName == "" {
		http.Error(w, "Channel name is required", http.StatusBadRequest)
		return
	}

	query := `
		INSERT INTO subscriptions (user_id, channel_name)
		VALUES ($1, $2)
		RETURNING id, user_id, channel_name, created_at
	`

	var sub models.Subscription
	err = h.db.QueryRow(query, userID, req.ChannelName).Scan(&sub.ID, &sub.UserID, &sub.ChannelName, &sub.CreatedAt)
	if err != nil {
		// Check if already subscribed (duplicate key constraint violation)
		if strings.Contains(err.Error(), "duplicate key") || strings.Contains(err.Error(), "unique constraint") {
			http.Error(w, "Already subscribed to this channel", http.StatusConflict)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(sub)
}

// Unsubscribe allows a user to unsubscribe from a channel
func (h *SubscriptionHandler) Unsubscribe(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, err := strconv.Atoi(vars["userId"])
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	channelName := vars["channelName"]
	if channelName == "" {
		http.Error(w, "Channel name is required", http.StatusBadRequest)
		return
	}

	query := `DELETE FROM subscriptions WHERE user_id = $1 AND channel_name = $2`
	result, err := h.db.Exec(query, userID, channelName)
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
		http.Error(w, "Subscription not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// GetUserSubscriptions returns all subscriptions for a user
func (h *SubscriptionHandler) GetUserSubscriptions(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, err := strconv.Atoi(vars["userId"])
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	query := `
		SELECT id, user_id, channel_name, created_at
		FROM subscriptions
		WHERE user_id = $1
		ORDER BY created_at DESC
	`

	rows, err := h.db.Query(query, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var subscriptions []models.Subscription
	for rows.Next() {
		var sub models.Subscription
		if err := rows.Scan(&sub.ID, &sub.UserID, &sub.ChannelName, &sub.CreatedAt); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		subscriptions = append(subscriptions, sub)
	}

	if err = rows.Err(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(subscriptions)
}

// CheckSubscription checks if a user is subscribed to a channel
func (h *SubscriptionHandler) CheckSubscription(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, err := strconv.Atoi(vars["userId"])
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	channelName := vars["channelName"]
	if channelName == "" {
		http.Error(w, "Channel name is required", http.StatusBadRequest)
		return
	}

	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM subscriptions WHERE user_id = $1 AND channel_name = $2)`
	err = h.db.QueryRow(query, userID, channelName).Scan(&exists)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]bool{"subscribed": exists})
}
