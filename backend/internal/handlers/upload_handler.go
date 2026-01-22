package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/aung-arata/youtube-clone/backend/internal/middleware"
	"github.com/aung-arata/youtube-clone/backend/internal/models"
	"github.com/aung-arata/youtube-clone/backend/internal/storage"
)

type UploadHandler struct {
	db      *sql.DB
	storage *storage.FileStorage
}

func NewUploadHandler(db *sql.DB, fileStorage *storage.FileStorage) *UploadHandler {
	return &UploadHandler{
		db:      db,
		storage: fileStorage,
	}
}

// UploadVideo handles video upload with multipart form data
func (h *UploadHandler) UploadVideo(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context (set by AuthMiddleware)
	_, ok := r.Context().Value(middleware.UserIDKey).(int)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	// Note: userID can be used in future for tracking uploads per user

	// Parse multipart form (max 500MB + 5MB for thumbnail)
	if err := r.ParseMultipartForm(510 * 1024 * 1024); err != nil {
		http.Error(w, "Error parsing form data", http.StatusBadRequest)
		return
	}

	// Get form fields
	title := r.FormValue("title")
	description := r.FormValue("description")
	category := r.FormValue("category")
	channelName := r.FormValue("channel_name")
	channelAvatar := r.FormValue("channel_avatar")

	// Validate required fields
	if strings.TrimSpace(title) == "" {
		http.Error(w, "Title is required", http.StatusBadRequest)
		return
	}
	if strings.TrimSpace(channelName) == "" {
		http.Error(w, "Channel name is required", http.StatusBadRequest)
		return
	}

	// Get video file
	videoFile, videoHeader, err := r.FormFile("video")
	if err != nil {
		http.Error(w, "Video file is required", http.StatusBadRequest)
		return
	}
	defer videoFile.Close()

	// Save video file
	videoURL, err := h.storage.SaveVideo(videoFile, videoHeader)
	if err != nil {
		http.Error(w, "Error saving video: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Get thumbnail file (optional)
	thumbnailURL := ""
	thumbnailFile, thumbnailHeader, err := r.FormFile("thumbnail")
	if err == nil {
		defer thumbnailFile.Close()
		thumbnailURL, err = h.storage.SaveThumbnail(thumbnailFile, thumbnailHeader)
		if err != nil {
			// Clean up video file if thumbnail fails
			h.storage.DeleteFile(videoURL)
			http.Error(w, "Error saving thumbnail: "+err.Error(), http.StatusBadRequest)
			return
		}
	}

	// Get video duration (should be calculated from the video file in production)
	duration := r.FormValue("duration")
	if duration == "" {
		duration = "00:00"
	}

	// Insert video into database
	query := `
		INSERT INTO videos (title, description, url, thumbnail, channel_name, channel_avatar, category, duration, views, likes, dislikes)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, 0, 0, 0)
		RETURNING id, title, description, url, thumbnail, channel_name, channel_avatar, views, likes, dislikes, category, duration, uploaded_at, created_at, updated_at
	`

	var video models.Video
	err = h.db.QueryRow(query, title, description, videoURL, thumbnailURL, channelName, channelAvatar, category, duration).Scan(
		&video.ID, &video.Title, &video.Description, &video.URL, &video.Thumbnail,
		&video.ChannelName, &video.ChannelAvatar, &video.Views, &video.Likes, &video.Dislikes,
		&video.Category, &video.Duration, &video.UploadedAt, &video.CreatedAt, &video.UpdatedAt,
	)

	if err != nil {
		// Clean up uploaded files if database insert fails
		h.storage.DeleteFile(videoURL)
		if thumbnailURL != "" {
			h.storage.DeleteFile(thumbnailURL)
		}
		http.Error(w, "Error creating video record: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Create video record
	// Note: In a complete implementation, you would:
	// - Store the userID as the uploader in a videos.user_id column
	// - Add audit logging with userID, timestamp, and file paths
	// - Track upload statistics per user

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(video)
}

// DeleteVideo handles video deletion including file cleanup
func (h *UploadHandler) DeleteVideo(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context
	_, ok := r.Context().Value(middleware.UserIDKey).(int)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Get video ID from path
	videoID := r.URL.Query().Get("id")
	if videoID == "" {
		http.Error(w, "Video ID is required", http.StatusBadRequest)
		return
	}

	// Get video details first to delete files
	query := `SELECT url, thumbnail FROM videos WHERE id = $1`
	var videoURL, thumbnailURL string
	err := h.db.QueryRow(query, videoID).Scan(&videoURL, &thumbnailURL)
	if err == sql.ErrNoRows {
		http.Error(w, "Video not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	// Delete video record from database
	deleteQuery := `DELETE FROM videos WHERE id = $1`
	_, err = h.db.Exec(deleteQuery, videoID)
	if err != nil {
		http.Error(w, "Error deleting video", http.StatusInternalServerError)
		return
	}

	// Delete files from storage
	h.storage.DeleteFile(videoURL)
	if thumbnailURL != "" {
		h.storage.DeleteFile(thumbnailURL)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Video deleted successfully"})
}
