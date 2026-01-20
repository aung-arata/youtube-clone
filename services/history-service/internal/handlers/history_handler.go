package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/aung-arata/youtube-clone/services/history-service/internal/models"
	"github.com/gorilla/mux"
)

type HistoryHandler struct {
	db         *sql.DB
	httpClient *http.Client
}

func NewHistoryHandler(db *sql.DB) *HistoryHandler {
	return &HistoryHandler{
		db: db,
		httpClient: &http.Client{
			Timeout: 5 * time.Second, // 5 second timeout for video service calls
		},
	}
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
		SELECT id, user_id, video_id, watched_at
		FROM watch_history
		WHERE user_id = $1
		ORDER BY watched_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := h.db.Query(query, userID, limit, offset)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var history []models.WatchHistory
	for rows.Next() {
		var item models.WatchHistory
		err := rows.Scan(&item.ID, &item.UserID, &item.VideoID, &item.WatchedAt)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		history = append(history, item)
	}

	// Enrich history with video details from video service
	var enrichedHistory []models.VideoWithHistory
	videoServiceURL := os.Getenv("VIDEO_SERVICE_URL")
	if videoServiceURL == "" {
		videoServiceURL = "http://video-service:8081" // default for docker-compose
	}

	for _, item := range history {
		// Fetch video details from video service
		videoURL := fmt.Sprintf("%s/videos/%d", videoServiceURL, item.VideoID)
		resp, err := h.httpClient.Get(videoURL)
		if err != nil {
			// If we can't fetch video details, skip this item
			// Note: when Get returns an error, resp is nil so no need to close
			continue
		}

		if resp.StatusCode == http.StatusOK {
			var video models.Video
			if err := json.NewDecoder(resp.Body).Decode(&video); err == nil {
				enrichedHistory = append(enrichedHistory, models.VideoWithHistory{
					Video:     video,
					WatchedAt: item.WatchedAt,
				})
			}
		}
		resp.Body.Close() // Close immediately after use, not deferred
	}

	if enrichedHistory == nil {
		enrichedHistory = []models.VideoWithHistory{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(enrichedHistory)
}
