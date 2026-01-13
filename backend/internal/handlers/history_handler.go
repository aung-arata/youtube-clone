package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/aung-arata/youtube-clone/backend/internal/models"
	"github.com/gorilla/mux"
)

type HistoryHandler struct {
	db *sql.DB
}

func NewHistoryHandler(db *sql.DB) *HistoryHandler {
	return &HistoryHandler{db: db}
}

// AddToHistory adds a video to user's watch history
func (h *HistoryHandler) AddToHistory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, err := strconv.Atoi(vars["userId"])
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	var req struct {
		VideoID int `json:"video_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if req.VideoID <= 0 {
		http.Error(w, "Video ID is required", http.StatusBadRequest)
		return
	}

	// Insert or update watch history (update timestamp if already exists)
	query := `
		INSERT INTO watch_history (user_id, video_id, watched_at)
		VALUES ($1, $2, CURRENT_TIMESTAMP)
		ON CONFLICT (user_id, video_id) 
		DO UPDATE SET watched_at = CURRENT_TIMESTAMP
		RETURNING id, user_id, video_id, watched_at
	`

	var history models.WatchHistory
	err = h.db.QueryRow(query, userID, req.VideoID).Scan(
		&history.ID, &history.UserID, &history.VideoID, &history.WatchedAt)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(history)
}

// GetHistory returns user's watch history with video details
func (h *HistoryHandler) GetHistory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, err := strconv.Atoi(vars["userId"])
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	// Get pagination parameters
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")

	page := 1
	limit := 20

	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}

	offset := (page - 1) * limit

	query := `
		SELECT 
			v.id, v.title, v.description, v.url, v.thumbnail, 
			v.channel_name, v.channel_avatar, v.views, v.likes, v.dislikes, 
			v.category, v.duration, v.uploaded_at, v.created_at, v.updated_at,
			wh.watched_at
		FROM watch_history wh
		JOIN videos v ON wh.video_id = v.id
		WHERE wh.user_id = $1
		ORDER BY wh.watched_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := h.db.Query(query, userID, limit, offset)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	type VideoWithHistory struct {
		models.Video
		WatchedAt string `json:"watched_at"`
	}

	var history []VideoWithHistory
	for rows.Next() {
		var item VideoWithHistory
		err := rows.Scan(
			&item.ID, &item.Title, &item.Description, &item.URL, &item.Thumbnail,
			&item.ChannelName, &item.ChannelAvatar, &item.Views, &item.Likes, &item.Dislikes,
			&item.Category, &item.Duration, &item.UploadedAt, &item.CreatedAt, &item.UpdatedAt,
			&item.WatchedAt)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		history = append(history, item)
	}

	if history == nil {
		history = []VideoWithHistory{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(history)
}
