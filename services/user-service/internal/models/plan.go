package models

import "time"

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

type Subscription struct {
	ID          int       `json:"id"`
	UserID      int       `json:"user_id"`
	ChannelName string    `json:"channel_name"`
	CreatedAt   time.Time `json:"created_at"`
}
