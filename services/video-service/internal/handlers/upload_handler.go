package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/aung-arata/youtube-clone/services/video-service/internal/middleware"
	"github.com/aung-arata/youtube-clone/services/video-service/internal/models"
	"github.com/aung-arata/youtube-clone/services/video-service/internal/storage"
	"github.com/aung-arata/youtube-clone/services/video-service/internal/transcoding"
	"github.com/gorilla/mux"
)

type UploadHandler struct {
	db            *sql.DB
	storage       *storage.FileStorage
	transcoding   *transcoding.TranscodingService
	uploadBaseDir string
}

func NewUploadHandler(db *sql.DB, fileStorage *storage.FileStorage, transcodingService *transcoding.TranscodingService, uploadBaseDir string) *UploadHandler {
	if uploadBaseDir == "" {
		uploadBaseDir = storage.UploadPath
	}

	return &UploadHandler{
		db:            db,
		storage:       fileStorage,
		transcoding:   transcodingService,
		uploadBaseDir: uploadBaseDir,
	}
}

// UploadVideo handles video upload with multipart form data
func (h *UploadHandler) UploadVideo(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(int)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	if err := r.ParseMultipartForm(510 * 1024 * 1024); err != nil {
		http.Error(w, "Error parsing form data", http.StatusBadRequest)
		return
	}

	title := r.FormValue("title")
	description := r.FormValue("description")
	category := r.FormValue("category")
	channelName := r.FormValue("channel_name")
	channelAvatar := r.FormValue("channel_avatar")
	visibility := strings.ToLower(strings.TrimSpace(r.FormValue("visibility")))
	if visibility == "" {
		visibility = "public"
	}
	if visibility != "public" && visibility != "unlisted" && visibility != "private" {
		http.Error(w, "Visibility must be public, unlisted, or private", http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(title) == "" {
		http.Error(w, "Title is required", http.StatusBadRequest)
		return
	}
	if strings.TrimSpace(channelName) == "" {
		http.Error(w, "Channel name is required", http.StatusBadRequest)
		return
	}

	videoFile, videoHeader, err := r.FormFile("video")
	if err != nil {
		http.Error(w, "Video file is required", http.StatusBadRequest)
		return
	}
	defer videoFile.Close()

	var videoID int
	if err := h.db.QueryRow(`SELECT nextval(pg_get_serial_sequence('videos', 'id'))`).Scan(&videoID); err != nil {
		http.Error(w, "Error reserving video ID", http.StatusInternalServerError)
		return
	}

	videoURL, err := h.storage.SaveVideo(videoFile, videoHeader)
	if err != nil {
		http.Error(w, "Error saving video: "+err.Error(), http.StatusBadRequest)
		return
	}

	thumbnailURL := ""
	thumbnailPath := ""
	thumbnailFile, thumbnailHeader, err := r.FormFile("thumbnail")
	if err == nil {
		defer thumbnailFile.Close()
		thumbnailURL, err = h.storage.SaveThumbnail(thumbnailFile, thumbnailHeader)
		if err != nil {
			h.storage.DeleteFile(videoURL)
			http.Error(w, "Error saving thumbnail: "+err.Error(), http.StatusBadRequest)
			return
		}
		thumbnailPath = filepath.Join(h.uploadBaseDir, strings.TrimPrefix(thumbnailURL, "/uploads/"))
	}

	duration := r.FormValue("duration")
	if duration == "" {
		duration = "00:00"
	}

	videoPath := filepath.Join(h.uploadBaseDir, strings.TrimPrefix(videoURL, "/uploads/"))
	thumbnailOutputDir := filepath.Join(h.uploadBaseDir, "thumbnails")

	if thumbnailURL == "" {
		if err := os.MkdirAll(thumbnailOutputDir, 0750); err != nil {
			h.storage.DeleteFile(videoURL)
			http.Error(w, "Error preparing thumbnail directory", http.StatusInternalServerError)
			return
		}

		extractedPath, extractErr := h.transcoding.ExtractThumbnail(videoPath, thumbnailOutputDir, videoID, 3.0)
		if extractErr != nil {
			h.storage.DeleteFile(videoURL)
			http.Error(w, "Error extracting thumbnail: "+extractErr.Error(), http.StatusInternalServerError)
			return
		}

		thumbnailPath = extractedPath
		thumbnailURL = "/uploads/thumbnails/" + filepath.Base(extractedPath)
	}

	query := `
INSERT INTO videos (id, user_id, title, description, url, thumbnail, channel_name, channel_avatar, category, duration, visibility, processing_status, views, likes, dislikes)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, 'pending', 0, 0, 0)
RETURNING id, user_id, title, description, url, thumbnail, channel_name, channel_avatar, visibility, processing_status, views, likes, dislikes, category, duration, uploaded_at, created_at, updated_at
`

	var video models.Video
	var userIDValue sql.NullInt64
	err = h.db.QueryRow(query, videoID, userID, title, description, videoURL, thumbnailURL, channelName, channelAvatar, category, duration, visibility).Scan(
		&video.ID, &userIDValue, &video.Title, &video.Description, &video.URL, &video.Thumbnail,
		&video.ChannelName, &video.ChannelAvatar, &video.Visibility, &video.ProcessingStatus, &video.Views, &video.Likes, &video.Dislikes,
		&video.Category, &video.Duration, &video.UploadedAt, &video.CreatedAt, &video.UpdatedAt,
	)
	if userIDValue.Valid {
		uid := int(userIDValue.Int64)
		video.UserID = &uid
	}

	if err != nil {
		h.storage.DeleteFile(videoURL)
		if thumbnailURL != "" {
			h.storage.DeleteFile(thumbnailURL)
		}
		http.Error(w, "Error creating video record: "+err.Error(), http.StatusInternalServerError)
		return
	}

	qualities := []string{"360p", "480p", "720p", "1080p"}
	if err := h.transcoding.QueueTranscoding(video.ID, videoPath, qualities); err != nil {
		h.storage.DeleteFile(videoURL)
		if thumbnailPath != "" {
			_ = os.Remove(thumbnailPath)
		}
		http.Error(w, "Error queueing transcoding jobs: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(video)
}

// GetTranscodingStatus returns transcoding status for a video.
func (h *UploadHandler) GetTranscodingStatus(w http.ResponseWriter, r *http.Request) {
	if _, ok := r.Context().Value(middleware.UserIDKey).(int); !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	videoID, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "Invalid video ID", http.StatusBadRequest)
		return
	}

	var processingStatus string
	if err := h.db.QueryRow(`SELECT processing_status FROM videos WHERE id = $1`, videoID).Scan(&processingStatus); err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Video not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	jobs, err := h.transcoding.GetTranscodingStatus(videoID)
	if err != nil {
		http.Error(w, "Error fetching transcoding jobs", http.StatusInternalServerError)
		return
	}

	type jobStatus struct {
		Quality  string `json:"quality"`
		Status   string `json:"status"`
		Progress int    `json:"progress"`
	}
	responseJobs := make([]jobStatus, 0, len(jobs))

	totalProgress := 0
	completed := 0
	failed := 0
	active := 0
	for _, job := range jobs {
		responseJobs = append(responseJobs, jobStatus{
			Quality:  job.TargetQuality,
			Status:   job.Status,
			Progress: job.Progress,
		})
		totalProgress += job.Progress
		switch job.Status {
		case "completed":
			completed++
		case "failed":
			failed++
		case "processing":
			active++
		}
	}

	overallProgress := 0
	if len(jobs) > 0 {
		overallProgress = totalProgress / len(jobs)
	}

	updatedStatus := processingStatus
	switch {
	case failed > 0:
		updatedStatus = "failed"
	case len(jobs) > 0 && completed == len(jobs):
		updatedStatus = "ready"
		overallProgress = 100
	case len(jobs) > 0 && (active > 0 || completed > 0):
		updatedStatus = "processing"
	}

	if updatedStatus != processingStatus {
		_, _ = h.db.Exec(`UPDATE videos SET processing_status = $1, updated_at = NOW() WHERE id = $2`, updatedStatus, videoID)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"video_id":          videoID,
		"processing_status": updatedStatus,
		"jobs":              responseJobs,
		"overall_progress":  overallProgress,
	})
}

type updateMetadataRequest struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
	Category    *string `json:"category"`
	Visibility  *string `json:"visibility"`
}

// UpdateMetadata updates metadata for a video owned by the authenticated user.
func (h *UploadHandler) UpdateMetadata(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(int)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	videoID, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "Invalid video ID", http.StatusBadRequest)
		return
	}

	var ownerID sql.NullInt64
	if err := h.db.QueryRow(`SELECT user_id FROM videos WHERE id = $1`, videoID).Scan(&ownerID); err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Video not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	if !ownerID.Valid || int(ownerID.Int64) != userID {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	var req updateMetadataRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	if req.Title != nil && strings.TrimSpace(*req.Title) == "" {
		http.Error(w, "Title cannot be empty", http.StatusBadRequest)
		return
	}
	if req.Visibility != nil {
		v := strings.ToLower(strings.TrimSpace(*req.Visibility))
		if v != "public" && v != "unlisted" && v != "private" {
			http.Error(w, "Visibility must be public, unlisted, or private", http.StatusBadRequest)
			return
		}
		req.Visibility = &v
	}

	setClauses := []string{}
	args := []interface{}{}
	argIndex := 1
	if req.Title != nil {
		title := strings.TrimSpace(*req.Title)
		setClauses = append(setClauses, fmt.Sprintf("title = $%d", argIndex))
		args = append(args, title)
		argIndex++
	}
	if req.Description != nil {
		setClauses = append(setClauses, fmt.Sprintf("description = $%d", argIndex))
		args = append(args, strings.TrimSpace(*req.Description))
		argIndex++
	}
	if req.Category != nil {
		setClauses = append(setClauses, fmt.Sprintf("category = $%d", argIndex))
		args = append(args, strings.TrimSpace(*req.Category))
		argIndex++
	}
	if req.Visibility != nil {
		setClauses = append(setClauses, fmt.Sprintf("visibility = $%d", argIndex))
		args = append(args, *req.Visibility)
		argIndex++
	}

	if len(setClauses) == 0 {
		http.Error(w, "No fields to update", http.StatusBadRequest)
		return
	}

	setClauses = append(setClauses, "updated_at = NOW()")
	args = append(args, videoID)
	query := fmt.Sprintf(`UPDATE videos SET %s WHERE id = $%d`, strings.Join(setClauses, ", "), argIndex)
	if _, err := h.db.Exec(query, args...); err != nil {
		http.Error(w, "Failed to update metadata", http.StatusInternalServerError)
		return
	}

	var video models.Video
	var userIDValue sql.NullInt64
	getQuery := `
SELECT id, user_id, title, description, url, thumbnail, channel_name, channel_avatar, visibility, processing_status, views, likes, dislikes, category, duration, uploaded_at, created_at, updated_at
FROM videos WHERE id = $1
`
	if err := h.db.QueryRow(getQuery, videoID).Scan(
		&video.ID, &userIDValue, &video.Title, &video.Description, &video.URL, &video.Thumbnail,
		&video.ChannelName, &video.ChannelAvatar, &video.Visibility, &video.ProcessingStatus, &video.Views, &video.Likes, &video.Dislikes,
		&video.Category, &video.Duration, &video.UploadedAt, &video.CreatedAt, &video.UpdatedAt,
	); err != nil {
		http.Error(w, "Failed to fetch updated video", http.StatusInternalServerError)
		return
	}
	if userIDValue.Valid {
		uid := int(userIDValue.Int64)
		video.UserID = &uid
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(video)
}

// DeleteVideo handles video deletion including file cleanup
func (h *UploadHandler) DeleteVideo(w http.ResponseWriter, r *http.Request) {
	_, ok := r.Context().Value(middleware.UserIDKey).(int)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	videoID := r.URL.Query().Get("id")
	if videoID == "" {
		http.Error(w, "Video ID is required", http.StatusBadRequest)
		return
	}

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

	deleteQuery := `DELETE FROM videos WHERE id = $1`
	_, err = h.db.Exec(deleteQuery, videoID)
	if err != nil {
		http.Error(w, "Error deleting video", http.StatusInternalServerError)
		return
	}

	h.storage.DeleteFile(videoURL)
	if thumbnailURL != "" {
		h.storage.DeleteFile(thumbnailURL)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Video deleted successfully"})
}
