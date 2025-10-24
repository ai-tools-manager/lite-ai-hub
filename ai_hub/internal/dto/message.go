// Package dto contains DTO objects for entities.
package dto

// Message represents a chat message with role and content.
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// GetMessagesResponse представляет ответ на запрос списка сообщений чата.
type GetMessagesResponse struct {
	Messages []Message `json:"messages"`
}

// MessageListResponse represents the response containing a list of messages.
type MessageListResponse struct {
	Messages []Message `json:"messages"`
}

// MessageRequest represents a request to send a message.
type MessageRequest struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}
