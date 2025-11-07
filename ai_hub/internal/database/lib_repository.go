package database

import (
	"database/sql"
	"lite_ai_hub/ai_hub/internal/models"
	"log"
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
		log.Printf("Error preparing statement for creating lib: %v", err)
		return err
	}
	defer func() {
		if err := stmt.Close(); err != nil {
			log.Printf("Error closing statement: %v", err)
		}
	}()

	lib.CreatedAt = time.Now()
	lib.UpdatedAt = time.Now()

	res, err := stmt.Exec(lib.Name, lib.Description, lib.Manifest, lib.CreatedAt, lib.UpdatedAt)
	if err != nil {
		log.Printf("Error executing statement for creating lib: %v", err)
		return err
	}

	id, err := res.LastInsertId()
	if err != nil {
		log.Printf("Error getting last insert ID for lib: %v", err)
		return err
	}
	lib.ID = uint(id)
	log.Printf("Successfully created lib with ID: %d", lib.ID)
	return nil
}

func (r *LibRepository) GetAll() ([]models.Lib, error) {
	rows, err := r.db.Query("SELECT id, name, description, manifest, created_at, updated_at FROM libs")
	if err != nil {
		log.Printf("Error querying all libs: %v", err)
		return nil, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Printf("Error closing rows: %v", err)
		}
	}()

	var libs []models.Lib
	for rows.Next() {
		var lib models.Lib
		if err := rows.Scan(&lib.ID, &lib.Name, &lib.Description, &lib.Manifest, &lib.CreatedAt, &lib.UpdatedAt); err != nil {
			log.Printf("Error scanning lib row: %v", err)
			return nil, err
		}
		libs = append(libs, lib)
	}
	log.Printf("Successfully retrieved %d libs", len(libs))
	return libs, nil
}

func (r *LibRepository) GetByID(id uint) (*models.Lib, error) {
	row := r.db.QueryRow("SELECT id, name, description, manifest, created_at, updated_at FROM libs WHERE id = ?", id)

	var lib models.Lib
	if err := row.Scan(&lib.ID, &lib.Name, &lib.Description, &lib.Manifest, &lib.CreatedAt, &lib.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			log.Printf("No lib found with id: %d", id)
			return nil, nil // Or a custom not found error
		}
		log.Printf("Error scanning lib row: %v", err)
		return nil, err
	}
	log.Printf("Successfully retrieved lib with ID: %d", id)
	return &lib, nil
}

func (r *LibRepository) Update(lib *models.Lib) error {
	stmt, err := r.db.Prepare("UPDATE libs SET name = ?, description = ?, manifest = ?, updated_at = ? WHERE id = ?")
	if err != nil {
		log.Printf("Error preparing statement for updating lib: %v", err)
		return err
	}
	defer func() {
		if err := stmt.Close(); err != nil {
			log.Printf("Error closing statement: %v", err)
		}
	}()

	lib.UpdatedAt = time.Now()

	_, err = stmt.Exec(lib.Name, lib.Description, lib.Manifest, lib.UpdatedAt, lib.ID)
	if err != nil {
		log.Printf("Error executing statement for updating lib with ID %d: %v", lib.ID, err)
	} else {
		log.Printf("Successfully updated lib with ID: %d", lib.ID)
	}
	return err
}

func (r *LibRepository) Delete(id uint) error {
	stmt, err := r.db.Prepare("DELETE FROM libs WHERE id = ?")
	if err != nil {
		log.Printf("Error preparing statement for deleting lib: %v", err)
		return err
	}
	defer func() {
		if err := stmt.Close(); err != nil {
			log.Printf("Error closing statement: %v", err)
		}
	}()

	_, err = stmt.Exec(id)
	if err != nil {
		log.Printf("Error executing statement for deleting lib with ID %d: %v", id, err)
	} else {
		log.Printf("Successfully deleted lib with ID: %d", id)
	}
	return err
}
