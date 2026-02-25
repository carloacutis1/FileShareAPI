package models

// this will instantiate a shared_file mock model 
// later, we will implement this into our PostreSQL database
// joshua fouch

type SharedFile struct {
	ID						int 			`json:"id"`
	FileID				int				`json:"file_id"`
	OwnerID				int				`json:"owner_id"`
	SharedWith		string		`json:"shared_with"`		
	SharedWithID	int				`json:"shared_with_id"`		
}
