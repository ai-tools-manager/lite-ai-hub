package models

import "time"

// Lib представляет библиотеку инструментов ИИ.
// Манифест содержит независимые от языка определения инструментов.
type Lib struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Manifest    string    `json:"manifest"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
