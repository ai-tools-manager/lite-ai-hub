package models

import "time"

type Session struct {
	ID        uint      `json:"id"`
	SessionID uint      `json:"session_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
