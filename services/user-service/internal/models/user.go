package models

import "time"

type User struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  string    `json:"-"` // Never send password in JSON responses
	Avatar    string    `json:"avatar"`
	Role      string    `json:"role"` // "user" or "admin"
	PlanID    *int      `json:"plan_id,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
