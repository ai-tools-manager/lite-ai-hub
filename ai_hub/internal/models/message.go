package models

import "time"

type Message struct {
	ID        uint      `json:"id"`
	ChatID    uint      `json:"chat_id"`
	Role      string    `json:"role"` // "user", "assistant", "tool"
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}
