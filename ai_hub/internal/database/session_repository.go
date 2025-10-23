package database

import (
	"database/sql"
	"lite_ai_hub/ai_hub/internal/models"
	"time"
)

type SessionRepository struct {
	db *sql.DB
}

func NewSessionRepository(db *sql.DB) *SessionRepository {
	return &SessionRepository{db: db}
}

func (r *SessionRepository) Create(session *models.Session) error {
	stmt, err := r.db.Prepare("INSERT INTO sessions(session_id, created_at, updated_at) VALUES(?, ?, ?)")
	if err != nil {
		return err
	}
	defer func() {
		if err := stmt.Close(); err != nil {
			// Log the error or handle it
		}
	}()

	session.CreatedAt = time.Now()
	session.UpdatedAt = time.Now()

	res, err := stmt.Exec(session.SessionID, session.CreatedAt, session.UpdatedAt)
	if err != nil {
		return err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	session.ID = uint(id)
	return nil
}

func (r *SessionRepository) GetBySessionID(sessionID uint) (*models.Session, error) {
	row := r.db.QueryRow("SELECT id, session_id, created_at, updated_at FROM sessions WHERE session_id = ?", sessionID)

	var session models.Session
	if err := row.Scan(&session.ID, &session.SessionID, &session.CreatedAt, &session.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Or a custom not found error
		}
		return nil, err
	}
	return &session, nil
}

func (r *SessionRepository) Delete(id uint) error {
	stmt, err := r.db.Prepare("DELETE FROM sessions WHERE id = ?")
	if err != nil {
		return err
	}
	defer func() {
		if err := stmt.Close(); err != nil {
			// Log the error or handle it
		}
	}()

	_, err = stmt.Exec(id)
	return err
}
