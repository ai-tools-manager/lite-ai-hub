package database

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func Connect() {
	var err error
	DB, err = sql.Open("sqlite3", "./sqlite.db")
	if err != nil {
		log.Fatal("failed to connect database", err)
	}

	createTables()
}

func createTables() {
	libsTable := `
	CREATE TABLE IF NOT EXISTS libs (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL UNIQUE,
		description TEXT,
		manifest TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	chatsTable := `
	CREATE TABLE IF NOT EXISTS chats (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		chat_id INTEGER NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	messagesTable := `
	CREATE TABLE IF NOT EXISTS messages (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		chat_id INTEGER NOT NULL,
		role TEXT NOT NULL,
		content TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY(chat_id) REFERENCES chats(id)
	);`

	_, err := DB.Exec(libsTable)
	if err != nil {
		log.Fatal("failed to create libs table", err)
	}

	_, err = DB.Exec(chatsTable)
	if err != nil {
		log.Fatal("failed to create chats table", err)
	}

	_, err = DB.Exec(messagesTable)
	if err != nil {
		log.Fatal("failed to create messages table", err)
	}
}
