package models

import "time"

type WatchHistory struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	VideoID   int       `json:"video_id"`
	WatchedAt time.Time `json:"watched_at"`
}
