package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/aung-arata/youtube-clone/services/video-service/internal/models"
	"github.com/gorilla/mux"
)

type PlaylistHandler struct {
	db *sql.DB
}

func NewPlaylistHandler(db *sql.DB) *PlaylistHandler {
	return &PlaylistHandler{db: db}
}

// CreatePlaylist creates a new playlist for a user
func (h *PlaylistHandler) CreatePlaylist(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, err := strconv.Atoi(vars["userId"])
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	var req struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(req.Name) == "" {
		http.Error(w, "Playlist name is required", http.StatusBadRequest)
		return
	}

	query := `
		INSERT INTO playlists (user_id, name, description)
		VALUES ($1, $2, $3)
		RETURNING id, user_id, name, description, created_at, updated_at
	`

	var playlist models.Playlist
	err = h.db.QueryRow(query, userID, req.Name, req.Description).Scan(
		&playlist.ID, &playlist.UserID, &playlist.Name, &playlist.Description, &playlist.CreatedAt, &playlist.UpdatedAt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(playlist)
}

// GetUserPlaylists returns all playlists for a user
func (h *PlaylistHandler) GetUserPlaylists(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, err := strconv.Atoi(vars["userId"])
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	query := `
		SELECT id, user_id, name, description, created_at, updated_at
		FROM playlists
		WHERE user_id = $1
		ORDER BY created_at DESC
	`

	rows, err := h.db.Query(query, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var playlists []models.Playlist
	for rows.Next() {
		var p models.Playlist
		if err := rows.Scan(&p.ID, &p.UserID, &p.Name, &p.Description, &p.CreatedAt, &p.UpdatedAt); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		playlists = append(playlists, p)
	}

	if err = rows.Err(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(playlists)
}

// GetPlaylist returns a specific playlist with its videos
func (h *PlaylistHandler) GetPlaylist(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	playlistID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid playlist ID", http.StatusBadRequest)
		return
	}

	// Get playlist details
	var playlist models.PlaylistWithVideos
	query := `
		SELECT id, user_id, name, description, created_at, updated_at
		FROM playlists
		WHERE id = $1
	`
	err = h.db.QueryRow(query, playlistID).Scan(
		&playlist.ID, &playlist.UserID, &playlist.Name, &playlist.Description, &playlist.CreatedAt, &playlist.UpdatedAt)
	if err == sql.ErrNoRows {
		http.Error(w, "Playlist not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Get videos in playlist
	videoQuery := `
		SELECT v.id, v.title, v.description, v.url, v.thumbnail, v.channel_name, v.channel_avatar, 
		       v.views, v.likes, v.dislikes, v.category, v.duration, v.uploaded_at, v.created_at, v.updated_at
		FROM videos v
		INNER JOIN playlist_videos pv ON pv.video_id = v.id
		WHERE pv.playlist_id = $1
		ORDER BY pv.position ASC
	`

	rows, err := h.db.Query(videoQuery, playlistID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var videos []models.Video
	for rows.Next() {
		var v models.Video
		if err := rows.Scan(&v.ID, &v.Title, &v.Description, &v.URL, &v.Thumbnail, &v.ChannelName,
			&v.ChannelAvatar, &v.Views, &v.Likes, &v.Dislikes, &v.Category, &v.Duration,
			&v.UploadedAt, &v.CreatedAt, &v.UpdatedAt); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		videos = append(videos, v)
	}

	if err = rows.Err(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	playlist.Videos = videos

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(playlist)
}

// UpdatePlaylist updates a playlist
func (h *PlaylistHandler) UpdatePlaylist(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	playlistID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid playlist ID", http.StatusBadRequest)
		return
	}

	var req struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(req.Name) == "" {
		http.Error(w, "Playlist name is required", http.StatusBadRequest)
		return
	}

	query := `
		UPDATE playlists
		SET name = $1, description = $2, updated_at = CURRENT_TIMESTAMP
		WHERE id = $3
		RETURNING id, user_id, name, description, created_at, updated_at
	`

	var playlist models.Playlist
	err = h.db.QueryRow(query, req.Name, req.Description, playlistID).Scan(
		&playlist.ID, &playlist.UserID, &playlist.Name, &playlist.Description, &playlist.CreatedAt, &playlist.UpdatedAt)
	if err == sql.ErrNoRows {
		http.Error(w, "Playlist not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(playlist)
}

// DeletePlaylist deletes a playlist
func (h *PlaylistHandler) DeletePlaylist(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	playlistID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid playlist ID", http.StatusBadRequest)
		return
	}

	query := `DELETE FROM playlists WHERE id = $1`
	result, err := h.db.Exec(query, playlistID)
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
		http.Error(w, "Playlist not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// AddVideoToPlaylist adds a video to a playlist
func (h *PlaylistHandler) AddVideoToPlaylist(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	playlistID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid playlist ID", http.StatusBadRequest)
		return
	}

	var req struct {
		VideoID int `json:"video_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Get next position
	var maxPosition int
	h.db.QueryRow("SELECT COALESCE(MAX(position), -1) FROM playlist_videos WHERE playlist_id = $1", playlistID).Scan(&maxPosition)

	query := `
		INSERT INTO playlist_videos (playlist_id, video_id, position)
		VALUES ($1, $2, $3)
		RETURNING id, playlist_id, video_id, position, added_at
	`

	var pv models.PlaylistVideo
	err = h.db.QueryRow(query, playlistID, req.VideoID, maxPosition+1).Scan(
		&pv.ID, &pv.PlaylistID, &pv.VideoID, &pv.Position, &pv.AddedAt)
	if err != nil {
		// Check if already in playlist
		if strings.Contains(err.Error(), "duplicate key") {
			http.Error(w, "Video already in playlist", http.StatusConflict)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(pv)
}

// RemoveVideoFromPlaylist removes a video from a playlist
func (h *PlaylistHandler) RemoveVideoFromPlaylist(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	playlistID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid playlist ID", http.StatusBadRequest)
		return
	}

	videoID, err := strconv.Atoi(vars["videoId"])
	if err != nil {
		http.Error(w, "Invalid video ID", http.StatusBadRequest)
		return
	}

	query := `DELETE FROM playlist_videos WHERE playlist_id = $1 AND video_id = $2`
	result, err := h.db.Exec(query, playlistID, videoID)
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
		http.Error(w, "Video not in playlist", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
