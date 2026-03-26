package models

import (
	"time"

	"gitlab.franciscan.edu/spr-26-backend/basement/internal/database"
)

// creates a struct with the same values as the sql table

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
    custom_attributes JSONB
);

*/

type File struct {
	id               string    `json:"id"`
	folder_id        string    `json:"folder_id"`
	name             string    `json:"name"`
	size             int       `json:"size"`
	file_type        string    `json:"file_type"`
	uploaded_at      time.Time `json:"uploaded_at"`
	description      string    `json:"description"`
	uploaded_by      string    `json:"uploaded_by"`
	tags             string    `json:"tags"`
	last_accessed    time.Time `json:"last_accessed"`
	checksum         string    `json:"checksum"`
	custom_attributes string   `json:"custom_attributes"` // JSONB as string for simplicity
}

// return the file slice from database
func GetFiles() ([]File, error) {
	query := `SELECT id, folder_id, name, size, file_type, uploaded_at, description, uploaded_by, tags, last_accessed, checksum, custom_attributes FROM files`
	rows, err := database.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var files []File
	for rows.Next() {
		var f File
		err := rows.Scan(&f.id, &f.folder_id, &f.name, &f.size, &f.file_type, &f.uploaded_at, &f.description, &f.uploaded_by, &f.tags, &f.last_accessed, &f.checksum, &f.custom_attributes)
		if err != nil {
			return nil, err
		}
		files = append(files, f)
	}
	return files, rows.Err()
}
