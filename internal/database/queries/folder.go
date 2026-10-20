package queries

// Database query functions for folder domain
// This layer contains ONLY database operations, no HTTP logic

import (
	"database/sql"

	"github.com/carloacutis1/FileShareAPI/internal/database"
	"github.com/carloacutis1/FileShareAPI/internal/models"
)

// GetAllFolders retrieves all folders from the database
func GetAllFolders() ([]models.Folder, error) {
	query := `SELECT id, parent_folder_id, name, created_at, user_id FROM folders ORDER BY created_at DESC`
	rows, err := database.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var folders []models.Folder
	for rows.Next() {
		var f models.Folder
		var parentID sql.NullString
		err := rows.Scan(&f.ID, &parentID, &f.Name, &f.CreatedAt, &f.UserID)
		if err != nil {
			return nil, err
		}
		f.ParentFolderID = parentID.String
		folders = append(folders, f)
	}
	return folders, rows.Err()
}

// GetFolderByID retrieves a single folder by its ID
func GetFolderByID(id string) (*models.Folder, error) {
	query := `SELECT id, parent_folder_id, name, created_at, user_id FROM folders WHERE id = $1`

	var f models.Folder
	var parentID sql.NullString

	err := database.DB.QueryRow(query, id).Scan(&f.ID, &parentID, &f.Name, &f.CreatedAt, &f.UserID)
	if err != nil {
		return nil, err
	}
	f.ParentFolderID = parentID.String
	return &f, nil
}

// CreateFolder inserts a new folder record into the database
func CreateFolder(folder models.Folder) error {
	query := `
	INSERT INTO folders (id, parent_folder_id, name, user_id)
	VALUES ($1, $2, $3, $4)
	`

	var parentID interface{}
	if folder.ParentFolderID == "" {
		parentID = nil
	} else {
		parentID = folder.ParentFolderID
	}

	_, err := database.DB.Exec(query, folder.ID, parentID, folder.Name, folder.UserID)
	return err
}

// UpdateFolder updates an existing folder record by ID
func UpdateFolder(id string, folder models.Folder) error {
	query := `
	UPDATE folders SET
		parent_folder_id = $1,
		name               = $2,
		user_id            = $3
	WHERE id = $4
	`

	var parentID interface{}
	if folder.ParentFolderID == "" {
		parentID = nil
	} else {
		parentID = folder.ParentFolderID
	}

	result, err := database.DB.Exec(query, parentID, folder.Name, folder.UserID, id)
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

// DeleteFolder removes a folder record from the database by ID
func DeleteFolder(id string) error {
	query := `DELETE FROM folders WHERE id = $1`
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
