package models

import "time"

// SearchResult represents a unified item found during a search
type SearchResult struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Type      string    `json:"type"` // "file" or "folder"
	CreatedAt time.Time `json:"created_at"`
	UserID    string    `json:"user_id"`
}
