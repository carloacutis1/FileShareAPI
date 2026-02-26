package models

import "time"

// this is a mock model for the folders table
// in the future, we will implement this into our PostgreSQL database
// Matthew Suess

// Folder matches the `folders` table schema.
type Folder struct {
	ID             string    `json:"id"`               // id (pk)
	ParentFolderID string    `json:"parent_folder_id"` // parent_folder_id
	Name           string    `json:"name"`             // name
	CreatedAt      time.Time `json:"created_at"`       // created_at (timestamp)
	UserID         string    `json:"user_id"`          // user_id
}

// mock folders slice
var folders = []Folder{
	{
		ID:             "fld_1",
		ParentFolderID: "",
		Name:           "Home",
		CreatedAt:      time.Now().Add(-72 * time.Hour),
		UserID:         "user_1",
	},
	{
		ID:             "fld_2",
		ParentFolderID: "fld_1",
		Name:           "Documents",
		CreatedAt:      time.Now().Add(-48 * time.Hour),
		UserID:         "user_1",
	},
	{
		ID:             "fld_3",
		ParentFolderID: "fld_1",
		Name:           "Pictures",
		CreatedAt:      time.Now().Add(-24 * time.Hour),
		UserID:         "user_2",
	},
}

// GetFolders returns the mock folders slice.
func GetFolders() []Folder {
	return folders
}