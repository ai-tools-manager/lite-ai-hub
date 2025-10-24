package dto

type CreateChatResponse struct {
	ChatID uint `json:"chat_id"`
}

type ChatListItem struct {
	ChatID uint `json:"chat_id"`
}

type GetChatsResponse struct {
	Chats []ChatListItem `json:"chats"`
}
