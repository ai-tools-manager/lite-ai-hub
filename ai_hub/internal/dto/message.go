package dto

// Message представляет сообщение в чате.
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// GetMessagesResponse представляет ответ на запрос списка сообщений чата.
type GetMessagesResponse struct {
	Messages []Message `json:"messages"`
}
