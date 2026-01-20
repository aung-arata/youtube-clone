package models

import "time"

type WatchHistory struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	VideoID   int       `json:"video_id"`
	WatchedAt time.Time `json:"watched_at"`
}

// Video represents video data from video service
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

// VideoWithHistory represents watch history with full video details
type VideoWithHistory struct {
	Video
	WatchedAt time.Time `json:"watched_at"`
}
