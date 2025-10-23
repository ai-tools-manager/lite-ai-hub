package database

import (
	"database/sql"
	"lite_ai_hub/ai_hub/internal/models"
	"time"
)

type MessageRepository struct {
	db *sql.DB
}

func NewMessageRepository(db *sql.DB) *MessageRepository {
	return &MessageRepository{db: db}
}

func (r *MessageRepository) Create(message *models.Message) error {
	stmt, err := r.db.Prepare("INSERT INTO messages(session_id, role, content, tool_call, created_at) VALUES(?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer func() {
		if err := stmt.Close(); err != nil {
			// Log the error or handle it
		}
	}()

	message.CreatedAt = time.Now()

	res, err := stmt.Exec(message.SessionID, message.Role, message.Content, message.ToolCall, message.CreatedAt)
	if err != nil {
		return err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	message.ID = uint(id)
	return nil
}

func (r *MessageRepository) GetBySessionID(sessionID uint) ([]models.Message, error) {
	rows, err := r.db.Query("SELECT id, session_id, role, content, tool_call, created_at FROM messages WHERE session_id = ?", sessionID)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			// Log the error or handle it
		}
	}()

	var messages []models.Message
	for rows.Next() {
		var msg models.Message
		if err := rows.Scan(&msg.ID, &msg.SessionID, &msg.Role, &msg.Content, &msg.ToolCall, &msg.CreatedAt); err != nil {
			return nil, err
		}
		messages = append(messages, msg)
	}
	return messages, nil
}
