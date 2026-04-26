package models

import "time"

// File represents a file record in the database
// This struct contains ONLY the data model definition, no database operations

// SQL model:
/*

CREATE TABLE files (
    id TEXT PRIMARY KEY,
    folder_id TEXT REFERENCES folders(id) ON DELETE SET NULL,
    name TEXT NOT NULL,
    size INTEGER NOT NULL,
    file_type TEXT,
    uploaded_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    description TEXT,
    uploaded_by TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    tags TEXT,
    last_accessed TIMESTAMP,
    checksum TEXT,
    storage_object_key TEXT,
    custom_attributes JSONB
);

*/

type File struct {
	ID               string    `json:"id"`
	SourceURL        string    `json:"source_url"`
	FolderID         string    `json:"folder_id"`
	Name             string    `json:"name"`
	Size             int       `json:"size"`
	FileType         string    `json:"file_type"`
	UploadedAt       time.Time `json:"uploaded_at"`
	Description      string    `json:"description"`
	UploadedBy       string    `json:"uploaded_by"`
	Tags             string    `json:"tags"`
	LastAccessed     time.Time `json:"last_accessed"`
	Checksum         string    `json:"checksum"`
	StorageObjectKey string    `json:"storage_object_key"` // object path within the GCS bucket (empty = no blob yet)
	CustomAttributes string    `json:"custom_attributes"`  // JSONB as string for simplicity
}
