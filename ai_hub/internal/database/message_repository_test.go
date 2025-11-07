package database

import (
	"database/sql"
	"lite_ai_hub/ai_hub/internal/models"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
)

func setupMessageTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("failed to open in-memory database: %v", err)
	}

	_, err = db.Exec("PRAGMA foreign_keys = ON;")
	if err != nil {
		t.Fatalf("failed to enable foreign keys: %v", err)
	}

	chatsTable := `
	CREATE TABLE IF NOT EXISTS chats (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		chat_id INTEGER NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`
	_, err = db.Exec(chatsTable)
	if err != nil {
		t.Fatalf("failed to create chats table: %v", err)
	}

	messagesTable := `
	CREATE TABLE IF NOT EXISTS messages (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		chat_id INTEGER NOT NULL,
		role TEXT NOT NULL,
		content TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY(chat_id) REFERENCES chats(id)
	);`
	_, err = db.Exec(messagesTable)
	if err != nil {
		t.Fatalf("failed to create messages table: %v", err)
	}

	return db
}

func TestMessageRepository_Create(t *testing.T) {
	db := setupMessageTestDB(t)
	defer db.Close()

	chatRepo := NewChatRepository(db)
	chat := &models.Chat{ChatID: 1}
	err := chatRepo.Create(chat)
	assert.NoError(t, err)

	repo := NewMessageRepository(db)

	message := &models.Message{
		ChatID:  chat.ID,
		Role:    "user",
		Content: "Hello",
	}

	err = repo.Create(message)
	assert.NoError(t, err)
	assert.NotZero(t, message.ID)
	assert.WithinDuration(t, time.Now(), message.CreatedAt, time.Second)
}

func TestMessageRepository_GetByChatID(t *testing.T) {
	db := setupMessageTestDB(t)
	defer db.Close()

	chatRepo := NewChatRepository(db)
	chat1 := &models.Chat{ChatID: 1}
	err := chatRepo.Create(chat1)
	assert.NoError(t, err)

	chat2 := &models.Chat{ChatID: 2}
	err = chatRepo.Create(chat2)
	assert.NoError(t, err)

	repo := NewMessageRepository(db)

	msg1 := &models.Message{ChatID: chat1.ID, Role: "user", Content: "Message 1"}
	msg2 := &models.Message{ChatID: chat1.ID, Role: "assistant", Content: "Message 2"}
	msg3 := &models.Message{ChatID: chat2.ID, Role: "user", Content: "Message 3"}

	err = repo.Create(msg1)
	assert.NoError(t, err)
	err = repo.Create(msg2)
	assert.NoError(t, err)
	err = repo.Create(msg3)
	assert.NoError(t, err)

	messages, err := repo.GetByChatID(chat1.ID)
	assert.NoError(t, err)
	assert.Len(t, messages, 2)
	assert.Equal(t, msg1.Content, messages[0].Content)
	assert.Equal(t, msg2.Content, messages[1].Content)

	messages, err = repo.GetByChatID(chat2.ID)
	assert.NoError(t, err)
	assert.Len(t, messages, 1)
	assert.Equal(t, msg3.Content, messages[0].Content)

	messages, err = repo.GetByChatID(999)
	assert.NoError(t, err)
	assert.Len(t, messages, 0)
}
