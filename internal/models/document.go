package models

import (
	"time"
)

type Document struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	ExpiresAt   time.Time `json:"expires_at"`
	IsActive    bool      `json:"is_active"`
	Versions    []Version `json:"versions"`
}

type Version struct {
	ID         string    `json:"id"`
	DocumentID string    `json:"document_id"`
	Version    string    `json:"version"`
	FilePath   string    `json:"file_path"`
	CreatedAt  time.Time `json:"created_at"`
	IsLatest   bool      `json:"is_latest"`
}
