package queries

import (
"database/sql"
"gitlab.franciscan.edu/spr-26-backend/basement/internal/database"
"gitlab.franciscan.edu/spr-26-backend/basement/internal/models"
)

func GetSharedFiles() ([]models.SharedFile, error) {
// run the db query
rows, err := database.DB.Query("SELECT id, file_id, owner_id, shared_with, permissions FROM shared_files")
if err != nil {
return nil, err
}
defer func() { _ = rows.Close() }()

var sharedFiles []models.SharedFile
for rows.Next() {
var sf models.SharedFile
if err := rows.Scan(&sf.ID, &sf.FileID, &sf.OwnerID, &sf.SharedWith, &sf.Permissions); err != nil {
return nil, err
}
sharedFiles = append(sharedFiles, sf)
}
// Check for errors from iterating over rows
if err := rows.Err(); err != nil {
return nil, err
}
return sharedFiles, nil
}

// CreateSharedFile inserts a new shared_file record into the database
func CreateSharedFile(sf models.SharedFile) error {
query := `
        INSERT INTO shared_files (
            id, file_id, owner_id, shared_with, permissions
        )
        VALUES ($1, $2, $3, $4, $5)
    `

_, err := database.DB.Exec(query,
sf.ID,
sf.FileID,
sf.OwnerID,
sf.SharedWith,
sf.Permissions,
)
return err
}

// UpdateSharedFile updates an existing shared_file record by ID
func UpdateSharedFile(id string, sf models.SharedFile) error {
query := `
UPDATE shared_files SET
file_id     = $1,
owner_id    = $2,
shared_with = $3,
permissions = $4
WHERE id = $5
`

result, err := database.DB.Exec(query,
sf.FileID,
sf.OwnerID,
sf.SharedWith,
sf.Permissions,
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

// DeleteSharedFile removes a shared_file record from the database by ID
func DeleteSharedFile(id string) error {
query := `DELETE FROM shared_files WHERE id = $1`
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

func GetSharedFile(id string) (models.SharedFile, error) {
query := `SELECT id, file_id, owner_id, shared_with, permissions FROM shared_files WHERE id = $1`
var sf models.SharedFile
err := database.DB.QueryRow(query, id).Scan(&sf.ID, &sf.FileID, &sf.OwnerID, &sf.SharedWith, &sf.Permissions)
if err != nil {
return sf, err
}
return sf, nil
}
