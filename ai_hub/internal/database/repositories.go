package database

import "lite_ai_hub/ai_hub/internal/models"

// ChatRepositoryInterface defines the interface for chat database operations.
type ChatRepositoryInterface interface {
	CreateChat(chat *models.Chat) error
	GetChatByID(id uint) (*models.Chat, error)
	GetAllChats() ([]models.Chat, error)
}

// LibRepositoryInterface defines the interface for library database operations.
type LibRepositoryInterface interface {
	CreateLib(lib *models.Lib) error
	GetLibByID(id uint) (*models.Lib, error)
	GetAllLibs() ([]models.Lib, error)
	DeleteLib(id uint) error
}

// MessageRepositoryInterface defines the interface for message database operations.
type MessageRepositoryInterface interface {
	CreateMessage(message *models.Message) error
	GetMessagesByChatID(chatID uint) ([]models.Message, error)
}
