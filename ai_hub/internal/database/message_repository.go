package database

import (
	"database/sql"
	"lite_ai_hub/ai_hub/internal/models"
	"log"
	"time"
)

type MessageRepository struct {
	db *sql.DB
}

func NewMessageRepository(db *sql.DB) *MessageRepository {
	return &MessageRepository{db: db}
}

func (r *MessageRepository) Create(message *models.Message) error {
	stmt, err := r.db.Prepare("INSERT INTO messages(chat_id, role, content, created_at) VALUES(?, ?, ?, ?)")
	if err != nil {
		log.Printf("Error preparing statement for creating message: %v", err)
		return err
	}
	defer func() {
		if err := stmt.Close(); err != nil {
			log.Printf("Error closing statement: %v", err)
		}
	}()

	message.CreatedAt = time.Now()

	res, err := stmt.Exec(message.ChatID, message.Role, message.Content, message.CreatedAt)
	if err != nil {
		log.Printf("Error executing statement for creating message: %v", err)
		return err
	}

	id, err := res.LastInsertId()
	if err != nil {
		log.Printf("Error getting last insert ID for message: %v", err)
		return err
	}
	message.ID = uint(id)
	log.Printf("Successfully created message for chat ID: %d", message.ChatID)
	return nil
}

func (r *MessageRepository) GetByChatID(chatID uint) ([]models.Message, error) {
	rows, err := r.db.Query("SELECT id, chat_id, role, content, created_at FROM messages WHERE chat_id = ?", chatID)
	if err != nil {
		log.Printf("Error querying messages by chat ID %d: %v", chatID, err)
		return nil, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Printf("Error closing rows: %v", err)
		}
	}()

	var messages []models.Message
	for rows.Next() {
		var msg models.Message
		if err := rows.Scan(&msg.ID, &msg.ChatID, &msg.Role, &msg.Content, &msg.CreatedAt); err != nil {
			log.Printf("Error scanning message row: %v", err)
			return nil, err
		}
		messages = append(messages, msg)
	}
	log.Printf("Successfully retrieved %d messages for chat ID: %d", len(messages), chatID)
	return messages, nil
}
