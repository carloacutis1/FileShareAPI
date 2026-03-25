package queries

import (
    "gitlab.franciscan.edu/spr-26-backend/basement/internal/models"
    "gitlab.franciscan.edu/spr-26-backend/basement/internal/database"
)

func GetSharedFiles() ([]models.SharedFile, error) {
    // run the db query
    rows, err := database.DB.Query("SELECT id, file_id, owner_id, shared_with, shared_with_id FROM shared_files")
    if err != nil {
        return nil, err
    }
    defer func() { _ = rows.Close() }()

    var sharedFiles []models.SharedFile
    for rows.Next() {
        var sf models.SharedFile
        if err := rows.Scan(&sf.ID, &sf.FileID, &sf.OwnerID, &sf.SharedWith, &sf.SharedWithID); err != nil {
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

