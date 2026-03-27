package models

import (
	"database/sql"
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
	CustomAttributes string    `json:"custom_attributes"` // JSONB as string for simplicity
}

// create a new file in the db
func CreateFile(f File) error {
	query := `
    INSERT INTO files (
        id, folder_id, name, size, file_type,
        description, uploaded_by, tags, last_accessed,
        checksum, custom_attributes
    )
    VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)
    `
	_, err := database.DB.Exec(query,
		f.ID,
		f.FolderID,
		f.Name,
		f.Size,
		f.FileType,
		f.Description,
		f.UploadedBy,
		f.Tags,
		f.LastAccessed,
		f.Checksum,
		f.CustomAttributes,
	)
	return err
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
		// Use sql.Null* for nullable columns to handle NULL from DB
		var folderID, fileType, description, tags, checksum, customAttributes sql.NullString
		var lastAccessed sql.NullTime
		err := rows.Scan(
			&f.ID, &folderID, &f.Name, &f.Size, &fileType,
			&f.UploadedAt, &description, &f.UploadedBy, &tags,
			&lastAccessed, &checksum, &customAttributes,
		)
		if err != nil {
			return nil, err
		}
		f.FolderID = folderID.String
		f.FileType = fileType.String
		f.Description = description.String
		f.Tags = tags.String
		f.LastAccessed = lastAccessed.Time
		f.Checksum = checksum.String
		f.CustomAttributes = customAttributes.String
		files = append(files, f)
	}
	return files, rows.Err()
}

// update an existing file in the db by id
func UpdateFile(f File) error {
	query := `
    UPDATE files SET
        folder_id         = $1,
        name              = $2,
        size              = $3,
        file_type         = $4,
        description       = $5,
        tags              = $6,
        last_accessed     = $7,
        checksum          = $8,
        custom_attributes = $9
    WHERE id = $10
    `
	result, err := database.DB.Exec(query,
		f.FolderID,
		f.Name,
		f.Size,
		f.FileType,
		f.Description,
		f.Tags,
		f.LastAccessed,
		f.Checksum,
		f.CustomAttributes,
		f.ID,
	)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}

// delete a file from the db by id
func DeleteFile(id string) error {
	query := `DELETE FROM files WHERE id = $1`
	result, err := database.DB.Exec(query, id)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}
