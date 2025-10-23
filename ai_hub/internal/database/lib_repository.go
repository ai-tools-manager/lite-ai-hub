package database

import (
	"database/sql"
	"lite_ai_hub/ai_hub/internal/models"
	"time"
)

type LibRepository struct {
	db *sql.DB
}

func NewLibRepository(db *sql.DB) *LibRepository {
	return &LibRepository{db: db}
}

func (r *LibRepository) Create(lib *models.Lib) error {
	stmt, err := r.db.Prepare("INSERT INTO libs(name, description, manifest, created_at, updated_at) VALUES(?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer func() {
		if err := stmt.Close(); err != nil {
			// Log the error or handle it
		}
	}()

	lib.CreatedAt = time.Now()
	lib.UpdatedAt = time.Now()

	res, err := stmt.Exec(lib.Name, lib.Description, lib.Manifest, lib.CreatedAt, lib.UpdatedAt)
	if err != nil {
		return err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	lib.ID = uint(id)
	return nil
}

func (r *LibRepository) GetAll() ([]models.Lib, error) {
	rows, err := r.db.Query("SELECT id, name, description, manifest, created_at, updated_at FROM libs")
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			// Log the error or handle it
		}
	}()

	var libs []models.Lib
	for rows.Next() {
		var lib models.Lib
		if err := rows.Scan(&lib.ID, &lib.Name, &lib.Description, &lib.Manifest, &lib.CreatedAt, &lib.UpdatedAt); err != nil {
			return nil, err
		}
		libs = append(libs, lib)
	}
	return libs, nil
}

func (r *LibRepository) GetByID(id uint) (*models.Lib, error) {
	row := r.db.QueryRow("SELECT id, name, description, manifest, created_at, updated_at FROM libs WHERE id = ?", id)

	var lib models.Lib
	if err := row.Scan(&lib.ID, &lib.Name, &lib.Description, &lib.Manifest, &lib.CreatedAt, &lib.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Or a custom not found error
		}
		return nil, err
	}
	return &lib, nil
}

func (r *LibRepository) Update(lib *models.Lib) error {
	stmt, err := r.db.Prepare("UPDATE libs SET name = ?, description = ?, manifest = ?, updated_at = ? WHERE id = ?")
	if err != nil {
		return err
	}
	defer func() {
		if err := stmt.Close(); err != nil {
			// Log the error or handle it
		}
	}()

	lib.UpdatedAt = time.Now()

	_, err = stmt.Exec(lib.Name, lib.Description, lib.Manifest, lib.UpdatedAt, lib.ID)
	return err
}

func (r *LibRepository) Delete(id uint) error {
	stmt, err := r.db.Prepare("DELETE FROM libs WHERE id = ?")
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
