package models

import "time"

type Comment struct {
	ID           int        `json:"id"`
	VideoID      int        `json:"video_id"`
	ParentID     *int       `json:"parent_id,omitempty"`
	UserID       int        `json:"user_id"`
	Username     string     `json:"username"`
	Content      string     `json:"content"`
	ReplyCount   int        `json:"reply_count,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}
