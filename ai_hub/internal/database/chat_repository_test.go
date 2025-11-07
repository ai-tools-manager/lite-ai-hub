package database

import (
	"database/sql"
	"lite_ai_hub/ai_hub/internal/models"
	"log"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
)

func setupTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("failed to open in-memory database: %v", err)
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
		log.Fatal("failed to create chats table", err)
	}

	return db
}

func TestChatRepository_Create(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewChatRepository(db)

	chat := &models.Chat{
		ChatID: 123,
	}

	err := repo.Create(chat)
	assert.NoError(t, err)
	assert.NotZero(t, chat.ID)
	assert.WithinDuration(t, time.Now(), chat.CreatedAt, time.Second)
	assert.WithinDuration(t, time.Now(), chat.UpdatedAt, time.Second)
}

func TestChatRepository_GetByChatID(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewChatRepository(db)

	chatToCreate := &models.Chat{
		ChatID: 123,
	}
	err := repo.Create(chatToCreate)
	assert.NoError(t, err)

	foundChat, err := repo.GetByChatID(123)
	assert.NoError(t, err)
	assert.NotNil(t, foundChat)
	assert.Equal(t, chatToCreate.ID, foundChat.ID)
	assert.Equal(t, chatToCreate.ChatID, foundChat.ChatID)

	notFoundChat, err := repo.GetByChatID(456)
	assert.NoError(t, err)
	assert.Nil(t, notFoundChat)
}

func TestChatRepository_Delete(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewChatRepository(db)

	chat := &models.Chat{
		ChatID: 123,
	}
	err := repo.Create(chat)
	assert.NoError(t, err)

	err = repo.Delete(chat.ID)
	assert.NoError(t, err)

	foundChat, err := repo.GetByChatID(chat.ChatID)
	assert.NoError(t, err)
	assert.Nil(t, foundChat)
}
