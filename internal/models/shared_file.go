package models

// joshua fouch

type SharedFile struct {
	ID          string `json:"id"`
	FileID      string `json:"file_id"`
	OwnerID     string `json:"owner_id"`
	SharedWith  string `json:"shared_with"` // this is an id of the user it is shared with
	Permissions string `json:"permissions"`
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
