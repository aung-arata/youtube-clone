package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/aung-arata/youtube-clone/services/comment-service/internal/models"
	"github.com/gorilla/mux"
)

type CommentHandler struct {
	db *sql.DB
}

func NewCommentHandler(db *sql.DB) *CommentHandler {
	return &CommentHandler{db: db}
}

// GetComments returns all top-level comments for a specific video (with reply counts)
func (h *CommentHandler) GetComments(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	videoID, err := strconv.Atoi(vars["videoId"])
	if err != nil {
		http.Error(w, "Invalid video ID", http.StatusBadRequest)
		return
	}

	query := `
		SELECT c.id, c.video_id, c.user_id, c.username, c.content, c.created_at, c.updated_at,
		       COUNT(r.id) AS reply_count
		FROM comments c
		LEFT JOIN comments r ON r.parent_id = c.id
		WHERE c.video_id = $1 AND c.parent_id IS NULL
		GROUP BY c.id
		ORDER BY c.created_at DESC
	`

	rows, err := h.db.Query(query, videoID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var comments []models.Comment
	for rows.Next() {
		var c models.Comment
		err := rows.Scan(&c.ID, &c.VideoID, &c.UserID, &c.Username, &c.Content, &c.CreatedAt, &c.UpdatedAt, &c.ReplyCount)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		comments = append(comments, c)
	}

	if comments == nil {
		comments = []models.Comment{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(comments)
}

// GetComment returns a single comment by ID
func (h *CommentHandler) GetComment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid comment ID", http.StatusBadRequest)
		return
	}

	query := `
		SELECT id, video_id, user_id, username, content, created_at, updated_at
		FROM comments
		WHERE id = $1
	`

	var c models.Comment
	err = h.db.QueryRow(query, id).Scan(&c.ID, &c.VideoID, &c.UserID, &c.Username, &c.Content, &c.CreatedAt, &c.UpdatedAt)

	if err == sql.ErrNoRows {
		http.Error(w, "Comment not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(c)
}

// CreateComment creates a new comment
func (h *CommentHandler) CreateComment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	videoID, err := strconv.Atoi(vars["videoId"])
	if err != nil {
		http.Error(w, "Invalid video ID", http.StatusBadRequest)
		return
	}

	var c models.Comment
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Set video ID from URL parameter
	c.VideoID = videoID

	// Validate required fields
	if strings.TrimSpace(c.Content) == "" {
		http.Error(w, "Content is required", http.StatusBadRequest)
		return
	}
	if c.UserID <= 0 {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}

	query := `
		INSERT INTO comments (video_id, user_id, username, content)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, updated_at
	`

	err = h.db.QueryRow(query, c.VideoID, c.UserID, c.Username, c.Content).Scan(&c.ID, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(c)
}

// UpdateComment updates an existing comment
func (h *CommentHandler) UpdateComment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid comment ID", http.StatusBadRequest)
		return
	}

	var c models.Comment
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Validate required fields
	if strings.TrimSpace(c.Content) == "" {
		http.Error(w, "Content is required", http.StatusBadRequest)
		return
	}

	query := `
		UPDATE comments 
		SET content = $1, updated_at = CURRENT_TIMESTAMP
		WHERE id = $2
		RETURNING id, video_id, user_id, username, content, created_at, updated_at
	`

	err = h.db.QueryRow(query, c.Content, id).Scan(&c.ID, &c.VideoID, &c.UserID, &c.Username, &c.Content, &c.CreatedAt, &c.UpdatedAt)

	if err == sql.ErrNoRows {
		http.Error(w, "Comment not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(c)
}

// DeleteComment deletes a comment
func (h *CommentHandler) DeleteComment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid comment ID", http.StatusBadRequest)
		return
	}

	query := `DELETE FROM comments WHERE id = $1`

	result, err := h.db.Exec(query, id)
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
		http.Error(w, "Comment not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// GetReplies returns all replies for a specific comment
func (h *CommentHandler) GetReplies(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	commentID, err := strconv.Atoi(vars["commentId"])
	if err != nil {
		http.Error(w, "Invalid comment ID", http.StatusBadRequest)
		return
	}

	query := `
		SELECT id, video_id, parent_id, user_id, username, content, created_at, updated_at
		FROM comments
		WHERE parent_id = $1
		ORDER BY created_at ASC
	`

	rows, err := h.db.Query(query, commentID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var replies []models.Comment
	for rows.Next() {
		var c models.Comment
		err := rows.Scan(&c.ID, &c.VideoID, &c.ParentID, &c.UserID, &c.Username, &c.Content, &c.CreatedAt, &c.UpdatedAt)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		replies = append(replies, c)
	}

	if replies == nil {
		replies = []models.Comment{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(replies)
}

// CreateReply creates a reply to a comment
func (h *CommentHandler) CreateReply(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	videoID, err := strconv.Atoi(vars["videoId"])
	if err != nil {
		http.Error(w, "Invalid video ID", http.StatusBadRequest)
		return
	}
	commentID, err := strconv.Atoi(vars["commentId"])
	if err != nil {
		http.Error(w, "Invalid comment ID", http.StatusBadRequest)
		return
	}

	var c models.Comment
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	c.VideoID = videoID
	c.ParentID = &commentID

	if strings.TrimSpace(c.Content) == "" {
		http.Error(w, "Content is required", http.StatusBadRequest)
		return
	}
	if c.UserID <= 0 {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}

	// Verify parent comment exists and belongs to this video
	var parentVideoID int
	err = h.db.QueryRow(`SELECT video_id FROM comments WHERE id = $1 AND parent_id IS NULL`, commentID).Scan(&parentVideoID)
	if err == sql.ErrNoRows {
		http.Error(w, "Parent comment not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	query := `
		INSERT INTO comments (video_id, parent_id, user_id, username, content)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, updated_at
	`

	err = h.db.QueryRow(query, c.VideoID, c.ParentID, c.UserID, c.Username, c.Content).Scan(&c.ID, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(c)
}
