package models

// UserModel to store user data below
type UserModel struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	DOB		  string 'json:"dob"'
	CreatedAt string `json:"createdat"`
	// Add more fields as needed
}

var userModels = []UserModel{
	// User Name
	{ID: 1, Name: "John Doe"},
	{ID: 2, Name: "Jane Doe"},
	{ID: 3, Name: "John Programming"},
	{ID: 4, Name: "Noname the Barbarian"},
	{ID: 5, Name: "Just Sayin"},
	// Have 3-5 sample items for testing
	// Include data for other variables in this scope - MH
}

func GetAllUserModels() []UserModel {
	return userModels
}
