package database

import (
	"database/sql"
	"lite_ai_hub/ai_hub/internal/models"
	"log"
	"time"
)

type ChatRepository struct {
	db *sql.DB
}

func NewChatRepository(db *sql.DB) *ChatRepository {
	return &ChatRepository{db: db}
}

func (r *ChatRepository) Create(chat *models.Chat) error {
	stmt, err := r.db.Prepare("INSERT INTO chats(chat_id, created_at, updated_at) VALUES(?, ?, ?)")
	if err != nil {
		log.Printf("Error preparing statement for creating chat: %v", err)
		return err
	}
	defer func() {
		if err := stmt.Close(); err != nil {
			log.Printf("Error closing statement: %v", err)
		}
	}()

	chat.CreatedAt = time.Now()
	chat.UpdatedAt = time.Now()

	res, err := stmt.Exec(chat.ChatID, chat.CreatedAt, chat.UpdatedAt)
	if err != nil {
		log.Printf("Error executing statement for creating chat: %v", err)
		return err
	}

	id, err := res.LastInsertId()
	if err != nil {
		log.Printf("Error getting last insert ID for chat: %v", err)
		return err
	}
	chat.ID = uint(id)
	log.Printf("Successfully created chat with ID: %d", chat.ID)
	return nil
}

func (r *ChatRepository) GetByChatID(chatID uint) (*models.Chat, error) {
	row := r.db.QueryRow("SELECT id, chat_id, created_at, updated_at FROM chats WHERE chat_id = ?", chatID)

	var chat models.Chat
	if err := row.Scan(&chat.ID, &chat.ChatID, &chat.CreatedAt, &chat.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			log.Printf("No chat found with chat_id: %d", chatID)
			return nil, nil // Or a custom not found error
		}
		log.Printf("Error scanning chat row: %v", err)
		return nil, err
	}
	return &chat, nil
}

func (r *ChatRepository) Delete(id uint) error {
	stmt, err := r.db.Prepare("DELETE FROM chats WHERE id = ?")
	if err != nil {
		log.Printf("Error preparing statement for deleting chat: %v", err)
		return err
	}
	defer func() {
		if err := stmt.Close(); err != nil {
			log.Printf("Error closing statement: %v", err)
		}
	}()

	_, err = stmt.Exec(id)
	if err != nil {
		log.Printf("Error executing statement for deleting chat: %v", err)
	} else {
		log.Printf("Successfully deleted chat with ID: %d", id)
	}
	return err
}
