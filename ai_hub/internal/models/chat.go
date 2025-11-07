package models

import "time"

type Chat struct {
	ID        uint      `json:"id"`
	ChatID    uint      `json:"chat_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
