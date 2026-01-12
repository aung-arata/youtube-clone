package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/aung-arata/youtube-clone/backend/internal/models"
	"github.com/gorilla/mux"
)

type VideoHandler struct {
	db *sql.DB
}

func NewVideoHandler(db *sql.DB) *VideoHandler {
	return &VideoHandler{db: db}
}

// GetVideos returns all videos with optional search and pagination
func (h *VideoHandler) GetVideos(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	searchQuery := r.URL.Query().Get("q")
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")

	// Default pagination values
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

	// Build query with optional search
	query := `
		SELECT id, title, description, url, thumbnail, channel_name, 
		       channel_avatar, views, likes, dislikes, duration, uploaded_at, created_at, updated_at
		FROM videos
	`
	
	var args []interface{}
	if searchQuery != "" {
		query += ` WHERE title ILIKE $1 OR description ILIKE $1 OR channel_name ILIKE $1`
		args = append(args, "%"+searchQuery+"%")
		query += fmt.Sprintf(` ORDER BY uploaded_at DESC LIMIT $2 OFFSET $3`)
		args = append(args, limit, offset)
	} else {
		query += fmt.Sprintf(` ORDER BY uploaded_at DESC LIMIT $1 OFFSET $2`)
		args = append(args, limit, offset)
	}

	rows, err := h.db.Query(query, args...)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var videos []models.Video
	for rows.Next() {
		var v models.Video
		err := rows.Scan(&v.ID, &v.Title, &v.Description, &v.URL, &v.Thumbnail,
			&v.ChannelName, &v.ChannelAvatar, &v.Views, &v.Likes, &v.Dislikes, &v.Duration,
			&v.UploadedAt, &v.CreatedAt, &v.UpdatedAt)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		videos = append(videos, v)
	}

	if videos == nil {
		videos = []models.Video{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(videos)
}

// GetVideo returns a single video by ID
func (h *VideoHandler) GetVideo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid video ID", http.StatusBadRequest)
		return
	}

	query := `
		SELECT id, title, description, url, thumbnail, channel_name, 
		       channel_avatar, views, likes, dislikes, duration, uploaded_at, created_at, updated_at
		FROM videos
		WHERE id = $1
	`

	var v models.Video
	err = h.db.QueryRow(query, id).Scan(&v.ID, &v.Title, &v.Description, &v.URL,
		&v.Thumbnail, &v.ChannelName, &v.ChannelAvatar, &v.Views, &v.Likes, &v.Dislikes, &v.Duration,
		&v.UploadedAt, &v.CreatedAt, &v.UpdatedAt)

	if err == sql.ErrNoRows {
		http.Error(w, "Video not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(v)
}

// CreateVideo creates a new video
func (h *VideoHandler) CreateVideo(w http.ResponseWriter, r *http.Request) {
	var v models.Video
	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Validate required fields
	if strings.TrimSpace(v.Title) == "" {
		http.Error(w, "Title is required", http.StatusBadRequest)
		return
	}
	if strings.TrimSpace(v.URL) == "" {
		http.Error(w, "URL is required", http.StatusBadRequest)
		return
	}
	if strings.TrimSpace(v.ChannelName) == "" {
		http.Error(w, "Channel name is required", http.StatusBadRequest)
		return
	}

	query := `
		INSERT INTO videos (title, description, url, thumbnail, channel_name, channel_avatar, duration)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, uploaded_at, created_at, updated_at
	`

	err := h.db.QueryRow(query, v.Title, v.Description, v.URL, v.Thumbnail,
		v.ChannelName, v.ChannelAvatar, v.Duration).Scan(&v.ID, &v.UploadedAt, &v.CreatedAt, &v.UpdatedAt)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(v)
}

// IncrementViews increments the view count for a video
func (h *VideoHandler) IncrementViews(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid video ID", http.StatusBadRequest)
		return
	}

	query := `
		UPDATE videos 
		SET views = views + 1, updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
		RETURNING views
	`

	var views int
	err = h.db.QueryRow(query, id).Scan(&views)
	if err == sql.ErrNoRows {
		http.Error(w, "Video not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]int{"views": views})
}

// LikeVideo increments the like count for a video
func (h *VideoHandler) LikeVideo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid video ID", http.StatusBadRequest)
		return
	}

	query := `
		UPDATE videos 
		SET likes = likes + 1, updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
		RETURNING likes
	`

	var likes int
	err = h.db.QueryRow(query, id).Scan(&likes)
	if err == sql.ErrNoRows {
		http.Error(w, "Video not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]int{"likes": likes})
}

// DislikeVideo increments the dislike count for a video
func (h *VideoHandler) DislikeVideo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid video ID", http.StatusBadRequest)
		return
	}

	query := `
		UPDATE videos 
		SET dislikes = dislikes + 1, updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
		RETURNING dislikes
	`

	var dislikes int
	err = h.db.QueryRow(query, id).Scan(&dislikes)
	if err == sql.ErrNoRows {
		http.Error(w, "Video not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]int{"dislikes": dislikes})
}
