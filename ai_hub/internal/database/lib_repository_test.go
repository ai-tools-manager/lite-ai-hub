package database

import (
	"database/sql"
	"lite_ai_hub/ai_hub/internal/models"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
)

func setupLibTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("failed to open in-memory database: %v", err)
	}

	libsTable := `
	CREATE TABLE IF NOT EXISTS libs (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL UNIQUE,
		description TEXT,
		manifest TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	_, err = db.Exec(libsTable)
	if err != nil {
		t.Fatalf("failed to create libs table: %v", err)
	}

	return db
}

func TestLibRepository_Create(t *testing.T) {
	db := setupLibTestDB(t)
	defer db.Close()

	repo := NewLibRepository(db)

	lib := &models.Lib{
		Name:        "Test Lib",
		Description: "A test library",
		Manifest:    "{}",
	}

	err := repo.Create(lib)
	assert.NoError(t, err)
	assert.NotZero(t, lib.ID)
	assert.WithinDuration(t, time.Now(), lib.CreatedAt, time.Second)
	assert.WithinDuration(t, time.Now(), lib.UpdatedAt, time.Second)
}

func TestLibRepository_GetAll(t *testing.T) {
	db := setupLibTestDB(t)
	defer db.Close()

	repo := NewLibRepository(db)

	lib1 := &models.Lib{Name: "Lib 1"}
	lib2 := &models.Lib{Name: "Lib 2"}
	repo.Create(lib1)
	repo.Create(lib2)

	libs, err := repo.GetAll()
	assert.NoError(t, err)
	assert.Len(t, libs, 2)
}

func TestLibRepository_GetByID(t *testing.T) {
	db := setupLibTestDB(t)
	defer db.Close()

	repo := NewLibRepository(db)

	lib := &models.Lib{Name: "Test Lib"}
	repo.Create(lib)

	foundLib, err := repo.GetByID(lib.ID)
	assert.NoError(t, err)
	assert.NotNil(t, foundLib)
	assert.Equal(t, lib.ID, foundLib.ID)
	assert.Equal(t, lib.Name, foundLib.Name)

	notFoundLib, err := repo.GetByID(999)
	assert.NoError(t, err)
	assert.Nil(t, notFoundLib)
}

func TestLibRepository_Update(t *testing.T) {
	db := setupLibTestDB(t)
	defer db.Close()

	repo := NewLibRepository(db)

	lib := &models.Lib{Name: "Original Name"}
	repo.Create(lib)

	lib.Name = "Updated Name"
	lib.Description = "Updated Description"
	err := repo.Update(lib)
	assert.NoError(t, err)

	updatedLib, err := repo.GetByID(lib.ID)
	assert.NoError(t, err)
	assert.Equal(t, "Updated Name", updatedLib.Name)
	assert.Equal(t, "Updated Description", updatedLib.Description)
	assert.True(t, updatedLib.UpdatedAt.After(updatedLib.CreatedAt))
}

func TestLibRepository_Delete(t *testing.T) {
	db := setupLibTestDB(t)
	defer db.Close()

	repo := NewLibRepository(db)

	lib := &models.Lib{Name: "Test Lib"}
	repo.Create(lib)

	err := repo.Delete(lib.ID)
	assert.NoError(t, err)

	foundLib, err := repo.GetByID(lib.ID)
	assert.NoError(t, err)
	assert.Nil(t, foundLib)
}
