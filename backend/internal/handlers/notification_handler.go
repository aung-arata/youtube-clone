package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/aung-arata/youtube-clone/backend/internal/models"
	"github.com/gorilla/mux"
)

type NotificationHandler struct {
	db *sql.DB
}

func NewNotificationHandler(db *sql.DB) *NotificationHandler {
	return &NotificationHandler{db: db}
}

// GetUserNotifications returns all notifications for a user
func (h *NotificationHandler) GetUserNotifications(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, err := strconv.Atoi(vars["userId"])
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	// Parse query parameters
	unreadOnly := r.URL.Query().Get("unread") == "true"
	limitStr := r.URL.Query().Get("limit")
	limit := 50
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}

	query := `
		SELECT id, user_id, type, title, message, link, is_read, created_at
		FROM notifications
		WHERE user_id = $1
	`
	
	args := []interface{}{userID}
	if unreadOnly {
		query += " AND is_read = FALSE"
	}
	
	query += " ORDER BY created_at DESC LIMIT $2"
	args = append(args, limit)

	rows, err := h.db.Query(query, args...)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var notifications []models.Notification
	for rows.Next() {
		var n models.Notification
		var link sql.NullString
		err := rows.Scan(&n.ID, &n.UserID, &n.Type, &n.Title, &n.Message, 
			&link, &n.IsRead, &n.CreatedAt)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if link.Valid {
			n.Link = link.String
		}
		notifications = append(notifications, n)
	}

	if notifications == nil {
		notifications = []models.Notification{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(notifications)
}

// CreateNotification creates a new notification
func (h *NotificationHandler) CreateNotification(w http.ResponseWriter, r *http.Request) {
	var req models.CreateNotificationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if req.UserID == 0 || req.Type == "" || req.Title == "" || req.Message == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	query := `
		INSERT INTO notifications (user_id, type, title, message, link)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at
	`

	var notification models.Notification
	err := h.db.QueryRow(query, req.UserID, req.Type, req.Title, req.Message, 
		sql.NullString{String: req.Link, Valid: req.Link != ""}).
		Scan(&notification.ID, &notification.CreatedAt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	notification.UserID = req.UserID
	notification.Type = req.Type
	notification.Title = req.Title
	notification.Message = req.Message
	notification.Link = req.Link
	notification.IsRead = false

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(notification)
}

// MarkAsRead marks a notification as read
func (h *NotificationHandler) MarkAsRead(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid notification ID", http.StatusBadRequest)
		return
	}

	query := `
		UPDATE notifications
		SET is_read = TRUE
		WHERE id = $1
		RETURNING id
	`

	var notificationID int
	err = h.db.QueryRow(query, id).Scan(&notificationID)
	if err == sql.ErrNoRows {
		http.Error(w, "Notification not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Notification marked as read",
		"id":      notificationID,
	})
}

// MarkAllAsRead marks all notifications for a user as read
func (h *NotificationHandler) MarkAllAsRead(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, err := strconv.Atoi(vars["userId"])
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	query := `
		UPDATE notifications
		SET is_read = TRUE
		WHERE user_id = $1 AND is_read = FALSE
	`

	result, err := h.db.Exec(query, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	rowsAffected, _ := result.RowsAffected()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":       "All notifications marked as read",
		"rows_affected": rowsAffected,
	})
}

// GetUnreadCount returns the count of unread notifications for a user
func (h *NotificationHandler) GetUnreadCount(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, err := strconv.Atoi(vars["userId"])
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	query := `
		SELECT COUNT(*)
		FROM notifications
		WHERE user_id = $1 AND is_read = FALSE
	`

	var count int
	err = h.db.QueryRow(query, userID).Scan(&count)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]int{"unread_count": count})
}
