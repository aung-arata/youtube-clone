package models

import "time"

type Video struct {
	ID            int       `json:"id"`
	Title         string    `json:"title"`
	Description   string    `json:"description"`
	URL           string    `json:"url"`
	Thumbnail     string    `json:"thumbnail"`
	ChannelName   string    `json:"channel_name"`
	ChannelAvatar string    `json:"channel_avatar"`
	Views         int       `json:"views"`
	Likes         int       `json:"likes"`
	Dislikes      int       `json:"dislikes"`
	Category      string    `json:"category"`
	Duration      string    `json:"duration"`
	UploadedAt    time.Time `json:"uploaded_at"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type User struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Avatar    string    `json:"avatar"`
	PlanID    *int      `json:"plan_id,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Comment struct {
	ID        int       `json:"id"`
	VideoID   int       `json:"video_id"`
	UserID    int       `json:"user_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type WatchHistory struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	VideoID   int       `json:"video_id"`
	WatchedAt time.Time `json:"watched_at"`
}

type Plan struct {
	ID                 int       `json:"id"`
	Name               string    `json:"name"`
	Price              float64   `json:"price"`
	MaxVideoQuality    string    `json:"max_video_quality"`
	MaxUploadsPerMonth int       `json:"max_uploads_per_month"`
	AdsFree            bool      `json:"ads_free"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}
