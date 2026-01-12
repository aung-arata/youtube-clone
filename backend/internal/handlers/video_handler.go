package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/aung-arata/youtube-clone/backend/internal/models"
	"github.com/gorilla/mux"
)

type VideoHandler struct {
	db *sql.DB
}

func NewVideoHandler(db *sql.DB) *VideoHandler {
	return &VideoHandler{db: db}
}

// GetVideos returns all videos
func (h *VideoHandler) GetVideos(w http.ResponseWriter, r *http.Request) {
	query := `
		SELECT id, title, description, url, thumbnail, channel_name, 
		       channel_avatar, views, duration, uploaded_at, created_at, updated_at
		FROM videos
		ORDER BY uploaded_at DESC
	`

	rows, err := h.db.Query(query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var videos []models.Video
	for rows.Next() {
		var v models.Video
		err := rows.Scan(&v.ID, &v.Title, &v.Description, &v.URL, &v.Thumbnail,
			&v.ChannelName, &v.ChannelAvatar, &v.Views, &v.Duration,
			&v.UploadedAt, &v.CreatedAt, &v.UpdatedAt)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		videos = append(videos, v)
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
		       channel_avatar, views, duration, uploaded_at, created_at, updated_at
		FROM videos
		WHERE id = $1
	`

	var v models.Video
	err = h.db.QueryRow(query, id).Scan(&v.ID, &v.Title, &v.Description, &v.URL,
		&v.Thumbnail, &v.ChannelName, &v.ChannelAvatar, &v.Views, &v.Duration,
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
