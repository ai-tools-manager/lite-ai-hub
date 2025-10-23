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

	sessionsTable := `
	CREATE TABLE IF NOT EXISTS sessions (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		session_id INTEGER NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	messagesTable := `
	CREATE TABLE IF NOT EXISTS messages (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		session_id INTEGER NOT NULL,
		role TEXT NOT NULL,
		content TEXT,
		tool_call TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY(session_id) REFERENCES sessions(id)
	);`

	_, err := DB.Exec(libsTable)
	if err != nil {
		log.Fatalf("could not create libs table: %v", err)
	}

	_, err = DB.Exec(sessionsTable)
	if err != nil {
		log.Fatalf("could not create sessions table: %v", err)
	}

	_, err = DB.Exec(messagesTable)
	if err != nil {
		log.Fatalf("could not create messages table: %v", err)
	}
}
