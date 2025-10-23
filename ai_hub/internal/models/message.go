package models

import "time"

type Message struct {
	ID        uint      `json:"id"`
	SessionID uint      `json:"session_id"`
	Role      string    `json:"role"` // "user", "assistant", "tool"
	Content   string    `json:"content"`
	ToolCall  string    `json:"tool_call,omitempty"` // JSON string for tool call
	CreatedAt time.Time `json:"created_at"`
}
