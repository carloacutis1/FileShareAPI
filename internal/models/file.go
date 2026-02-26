package models

// this instantiates a mock file model
// later, will make into actual databse
// HC

type File struct {
	ID      int `json:"id"`
	FileID  int `json:"file_id"`
	OwnerID int `json:"owner_id"`
}

// mock data to test functionality
var files = []File{
	{ID: 1, FileID: 1001, OwnerID: 42},
	{ID: 2, FileID: 1002, OwnerID: 42},
	{ID: 3, FileID: 1003, OwnerID: 99},
	{ID: 4, FileID: 1004, OwnerID: 100},
	{ID: 5, FileID: 1005, OwnerID: 101},
}

// return the file slice
func GetFiles() []File {
	return files
}
