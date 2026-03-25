package models

import (
    "gitlab.franciscan.edu/spr-26-backend/basement/internal/database"
)

// joshua fouch

type SharedFile struct {
    ID           int    `json:"id"`
    FileID       int    `json:"file_id"`
    OwnerID      int    `json:"owner_id"`
    SharedWith   string `json:"shared_with"`
    SharedWithID int    `json:"shared_with_id"`
}

// instantiate some SharedFiles in a slice
// go has a built in vector called a slice declared by []name
//var sharedFiles = []SharedFile{
//	{ID: 1, FileID: 101, OwnerID: 1, SharedWith: "alice", SharedWithID: 5},
//	{ID: 2, FileID: 102, OwnerID: 2, SharedWith: "bob", SharedWithID: 6},
//	{ID: 3, FileID: 103, OwnerID: 1, SharedWith: "chris", SharedWithID: 7},
//	{ID: 4, FileID: 104, OwnerID: 4, SharedWith: "charlie", SharedWithID: 8},
//	{ID: 5, FileID: 105, OwnerID: 3, SharedWith: "bob", SharedWithID: 9},
//}

// return the sharedfile slice
//func GetSharedFiles() []SharedFile {
//	return sharedFiles
//}

func GetSharedFiles() ([]SharedFile, error) {
    // run the db query
    rows, err := database.DB.Query("SELECT id, file_id, owner_id, shared_with, shared_with_id FROM shared_files")
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var sharedFiles []SharedFile
    for rows.Next() {
        var sf SharedFile
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
