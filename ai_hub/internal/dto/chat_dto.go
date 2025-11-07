// Package dto contains DTO objects for entities.
package dto

import "github.com/google/uuid"

// ChatResponse represents the response when creating a chat.
type ChatResponse struct {
	ChatID uuid.UUID `json:"chat_id"`
}

// ChatListResponse represents the response containing a list of chats.
type ChatListResponse struct {
	Chats []ChatResponse `json:"chats"`
}
