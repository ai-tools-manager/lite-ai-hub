// Package dto contains DTO objects for entities.
package dto

import (
	"net/url"

	"github.com/google/uuid"
)

// Library represents an installed library with metadata.
type Library struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	GitURL      url.URL   `json:"git_url"`
}

// LibraryListResponse represents the response containing a list of libraries.
type LibraryListResponse []Library

// LibraryRequest represents a request to install a library.
type LibraryRequest struct {
	GitURL url.URL `json:"git_url"`
}
