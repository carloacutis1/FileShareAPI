package queries

// Database query functions for file domain
// This layer contains ONLY database operations, no HTTP logic

import (
	"database/sql"

	"gitlab.franciscan.edu/spr-26-backend/basement/internal/database"
	"gitlab.franciscan.edu/spr-26-backend/basement/internal/models"
)

// GetAllFiles retrieves all files from the database
func GetAllFiles() ([]models.File, error) {
	query := `SELECT id, folder_id, name, size, file_type, uploaded_at, description, uploaded_by, tags, last_accessed, checksum, storage_object_key, custom_attributes FROM files`
	rows, err := database.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var files []models.File
	for rows.Next() {
		var f models.File
		// Use sql.Null* for nullable columns to handle NULL from DB
		var folderID, fileType, description, tags, checksum, storageKey, customAttributes sql.NullString
		var lastAccessed sql.NullTime
		err := rows.Scan(
			&f.ID, &folderID, &f.Name, &f.Size, &fileType,
			&f.UploadedAt, &description, &f.UploadedBy, &tags,
			&lastAccessed, &checksum, &storageKey, &customAttributes,
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
		f.StorageObjectKey = storageKey.String
		f.CustomAttributes = customAttributes.String
		files = append(files, f)
	}
	return files, rows.Err()
}

// GetFileByID retrieves a single file by its ID
func GetFileByID(id string) (*models.File, error) {
	query := `SELECT id, folder_id, name, size, file_type, uploaded_at, description, uploaded_by, tags, last_accessed, checksum, storage_object_key, custom_attributes FROM files WHERE id = $1`

	var f models.File
	var folderID, fileType, description, tags, checksum, storageKey, customAttributes sql.NullString
	var lastAccessed sql.NullTime

	err := database.DB.QueryRow(query, id).Scan(
		&f.ID, &folderID, &f.Name, &f.Size, &fileType,
		&f.UploadedAt, &description, &f.UploadedBy, &tags,
		&lastAccessed, &checksum, &storageKey, &customAttributes,
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
	f.StorageObjectKey = storageKey.String
	f.CustomAttributes = customAttributes.String

	return &f, nil
}

// CreateFile inserts a new file record into the database
func CreateFile(file models.File) error {
	query := `
    INSERT INTO files (
        id, folder_id, name, size, file_type,
        description, uploaded_by, tags, last_accessed,
        checksum, storage_object_key, custom_attributes
    )
    VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)
    `

	// Handle nullable fields - convert empty strings to NULL for database
	var folderID interface{}
	if file.FolderID == "" {
		folderID = nil
	} else {
		folderID = file.FolderID
	}

	var storageKey interface{}
	if file.StorageObjectKey == "" {
		storageKey = nil
	} else {
		storageKey = file.StorageObjectKey
	}

	var customAttrs interface{}
	if file.CustomAttributes == "" {
		customAttrs = nil
	} else {
		customAttrs = file.CustomAttributes
	}

	_, err := database.DB.Exec(query,
		file.ID,
		folderID,
		file.Name,
		file.Size,
		file.FileType,
		file.Description,
		file.UploadedBy,
		file.Tags,
		file.LastAccessed,
		file.Checksum,
		storageKey,
		customAttrs,
	)
	return err
}

// UpdateFile updates an existing file record by ID
func UpdateFile(id string, file models.File) error {
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
        storage_object_key = $9,
        custom_attributes = $10
    WHERE id = $11
    `

	// Handle nullable fields - convert empty strings to NULL for database
	var folderID interface{}
	if file.FolderID == "" {
		folderID = nil
	} else {
		folderID = file.FolderID
	}

	var storageKey interface{}
	if file.StorageObjectKey == "" {
		storageKey = nil
	} else {
		storageKey = file.StorageObjectKey
	}

	var customAttrs interface{}
	if file.CustomAttributes == "" {
		customAttrs = nil
	} else {
		customAttrs = file.CustomAttributes
	}

	result, err := database.DB.Exec(query,
		folderID,
		file.Name,
		file.Size,
		file.FileType,
		file.Description,
		file.Tags,
		file.LastAccessed,
		file.Checksum,
		storageKey,
		customAttrs,
		id,
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

// DeleteFile removes a file record from the database by ID
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
