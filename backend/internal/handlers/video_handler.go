package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/aung-arata/youtube-clone/backend/internal/models"
	"github.com/gorilla/mux"
)

type VideoHandler struct {
	db *sql.DB
}

func NewVideoHandler(db *sql.DB) *VideoHandler {
	return &VideoHandler{db: db}
}

// GetVideos returns all videos with optional search, category filter and pagination
func (h *VideoHandler) GetVideos(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	searchQuery := r.URL.Query().Get("q")
	category := r.URL.Query().Get("category")
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")
	sortBy := r.URL.Query().Get("sort_by")      // views, likes, date, title
	orderBy := r.URL.Query().Get("order")        // asc, desc
	uploadedAfter := r.URL.Query().Get("uploaded_after")  // date filter
	minDuration := r.URL.Query().Get("min_duration")      // minimum duration in seconds
	maxDuration := r.URL.Query().Get("max_duration")      // maximum duration in seconds

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

	// Build query with optional search and category filter
	query := `
		SELECT id, title, description, url, thumbnail, channel_name, 
		       channel_avatar, views, likes, dislikes, category, duration, uploaded_at, created_at, updated_at
		FROM videos
	`
	
	var args []interface{}
	var conditions []string
	argIndex := 1

	if searchQuery != "" {
		conditions = append(conditions, fmt.Sprintf("(title ILIKE $%d OR description ILIKE $%d OR channel_name ILIKE $%d)", argIndex, argIndex, argIndex))
		args = append(args, "%"+searchQuery+"%")
		argIndex++
	}

	if category != "" {
		conditions = append(conditions, fmt.Sprintf("category = $%d", argIndex))
		args = append(args, category)
		argIndex++
	}

	// Add uploaded_after filter with date validation
	if uploadedAfter != "" {
		// Validate date format (ISO 8601)
		if _, err := time.Parse(time.RFC3339, uploadedAfter); err != nil {
			http.Error(w, "Invalid date format for uploaded_after. Use ISO 8601 format (e.g., 2024-01-15T00:00:00Z)", http.StatusBadRequest)
			return
		}
		conditions = append(conditions, fmt.Sprintf("uploaded_at >= $%d", argIndex))
		args = append(args, uploadedAfter)
		argIndex++
	}

	// Add duration filters with error handling
	if minDuration != "" {
		minDur, err := strconv.Atoi(minDuration)
		if err != nil {
			http.Error(w, "Invalid min_duration value. Must be an integer representing seconds", http.StatusBadRequest)
			return
		}
		if minDur < 0 {
			http.Error(w, "min_duration must be a positive number", http.StatusBadRequest)
			return
		}
		conditions = append(conditions, fmt.Sprintf("EXTRACT(EPOCH FROM (duration::interval)) >= $%d", argIndex))
		args = append(args, minDur)
		argIndex++
	}

	if maxDuration != "" {
		maxDur, err := strconv.Atoi(maxDuration)
		if err != nil {
			http.Error(w, "Invalid max_duration value. Must be an integer representing seconds", http.StatusBadRequest)
			return
		}
		if maxDur < 0 {
			http.Error(w, "max_duration must be a positive number", http.StatusBadRequest)
			return
		}
		conditions = append(conditions, fmt.Sprintf("EXTRACT(EPOCH FROM (duration::interval)) <= $%d", argIndex))
		args = append(args, maxDur)
		argIndex++
	}

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	// Add sorting
	orderClause := "uploaded_at DESC" // default sorting
	if sortBy != "" {
		order := "DESC"
		if orderBy == "asc" {
			order = "ASC"
		}
		
		switch sortBy {
		case "views":
			orderClause = "views " + order
		case "likes":
			orderClause = "likes " + order
		case "date":
			orderClause = "uploaded_at " + order
		case "title":
			orderClause = "title " + order
		default:
			orderClause = "uploaded_at DESC"
		}
	}

	query += fmt.Sprintf(` ORDER BY %s LIMIT $%d OFFSET $%d`, orderClause, argIndex, argIndex+1)
	args = append(args, limit, offset)

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
			&v.ChannelName, &v.ChannelAvatar, &v.Views, &v.Likes, &v.Dislikes, &v.Category, &v.Duration,
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
		       channel_avatar, views, likes, dislikes, category, duration, uploaded_at, created_at, updated_at
		FROM videos
		WHERE id = $1
	`

	var v models.Video
	err = h.db.QueryRow(query, id).Scan(&v.ID, &v.Title, &v.Description, &v.URL,
		&v.Thumbnail, &v.ChannelName, &v.ChannelAvatar, &v.Views, &v.Likes, &v.Dislikes, &v.Category, &v.Duration,
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

// GetCategories returns all unique video categories
func (h *VideoHandler) GetCategories(w http.ResponseWriter, r *http.Request) {
	query := `
		SELECT DISTINCT category 
		FROM videos 
		WHERE category IS NOT NULL AND category != ''
		ORDER BY category
	`

	rows, err := h.db.Query(query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var categories []string
	for rows.Next() {
		var category string
		if err := rows.Scan(&category); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		categories = append(categories, category)
	}

	if categories == nil {
		categories = []string{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(categories)
}

// GetRecommendations returns recommended videos based on a given video
// Algorithm: Returns videos from the same category, sorted by views, excluding the current video
func (h *VideoHandler) GetRecommendations(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid video ID", http.StatusBadRequest)
		return
	}

	// Get limit from query parameter
	limitStr := r.URL.Query().Get("limit")
	limit := 10 // Default limit
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 50 {
			limit = l
		}
	}

	// First, get the category of the current video
	var category string
	err = h.db.QueryRow("SELECT category FROM videos WHERE id = $1", id).Scan(&category)
	if err == sql.ErrNoRows {
		http.Error(w, "Video not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Get recommended videos from the same category, sorted by views
	query := `
		SELECT id, title, description, url, thumbnail, channel_name, channel_avatar,
		       views, likes, dislikes, category, duration, uploaded_at, created_at, updated_at
		FROM videos
		WHERE category = $1 AND id != $2
		ORDER BY views DESC, uploaded_at DESC
		LIMIT $3
	`

	rows, err := h.db.Query(query, category, id, limit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var videos []models.Video
	for rows.Next() {
		var v models.Video
		err := rows.Scan(&v.ID, &v.Title, &v.Description, &v.URL, &v.Thumbnail,
			&v.ChannelName, &v.ChannelAvatar, &v.Views, &v.Likes, &v.Dislikes,
			&v.Category, &v.Duration, &v.UploadedAt, &v.CreatedAt, &v.UpdatedAt)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		videos = append(videos, v)
	}

	if err = rows.Err(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if videos == nil {
		videos = []models.Video{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(videos)
}

// GetTrendingVideos returns videos sorted by views in the last 7 days
func (h *VideoHandler) GetTrendingVideos(w http.ResponseWriter, r *http.Request) {
	limitStr := r.URL.Query().Get("limit")
	limit := 20
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}

	// Get videos uploaded or viewed recently, sorted by views
	query := `
		SELECT id, title, description, url, thumbnail, channel_name, 
		       channel_avatar, views, likes, dislikes, category, duration, uploaded_at, created_at, updated_at
		FROM videos
		WHERE uploaded_at >= NOW() - INTERVAL '7 days'
		ORDER BY views DESC, likes DESC
		LIMIT $1
	`

	rows, err := h.db.Query(query, limit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var videos []models.Video
	for rows.Next() {
		var v models.Video
		err := rows.Scan(&v.ID, &v.Title, &v.Description, &v.URL, &v.Thumbnail,
			&v.ChannelName, &v.ChannelAvatar, &v.Views, &v.Likes, &v.Dislikes,
			&v.Category, &v.Duration, &v.UploadedAt, &v.CreatedAt, &v.UpdatedAt)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		videos = append(videos, v)
	}

	if err = rows.Err(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if videos == nil {
		videos = []models.Video{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(videos)
}

// GetPopularVideos returns most viewed videos of all time
func (h *VideoHandler) GetPopularVideos(w http.ResponseWriter, r *http.Request) {
	limitStr := r.URL.Query().Get("limit")
	limit := 20
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}

	query := `
		SELECT id, title, description, url, thumbnail, channel_name, 
		       channel_avatar, views, likes, dislikes, category, duration, uploaded_at, created_at, updated_at
		FROM videos
		ORDER BY views DESC, likes DESC
		LIMIT $1
	`

	rows, err := h.db.Query(query, limit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var videos []models.Video
	for rows.Next() {
		var v models.Video
		err := rows.Scan(&v.ID, &v.Title, &v.Description, &v.URL, &v.Thumbnail,
			&v.ChannelName, &v.ChannelAvatar, &v.Views, &v.Likes, &v.Dislikes,
			&v.Category, &v.Duration, &v.UploadedAt, &v.CreatedAt, &v.UpdatedAt)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		videos = append(videos, v)
	}

	if err = rows.Err(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if videos == nil {
		videos = []models.Video{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(videos)
}

// GetVideoAnalytics returns detailed analytics for a specific video
func (h *VideoHandler) GetVideoAnalytics(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid video ID", http.StatusBadRequest)
		return
	}

	query := `
		SELECT id, title, views, likes, dislikes, category, uploaded_at
		FROM videos
		WHERE id = $1
	`

	var analytics struct {
		ID           int    `json:"id"`
		Title        string `json:"title"`
		Views        int    `json:"views"`
		Likes        int    `json:"likes"`
		Dislikes     int    `json:"dislikes"`
		Category     string `json:"category"`
		UploadedAt   string `json:"uploaded_at"`
		LikeRatio    float64 `json:"like_ratio"`
		Engagement   int    `json:"engagement"`
	}

	err = h.db.QueryRow(query, id).Scan(&analytics.ID, &analytics.Title, &analytics.Views,
		&analytics.Likes, &analytics.Dislikes, &analytics.Category, &analytics.UploadedAt)
	if err == sql.ErrNoRows {
		http.Error(w, "Video not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Calculate analytics metrics
	totalReactions := analytics.Likes + analytics.Dislikes
	if totalReactions > 0 {
		analytics.LikeRatio = float64(analytics.Likes) / float64(totalReactions) * 100
	}
	analytics.Engagement = analytics.Likes + analytics.Dislikes

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(analytics)
}

