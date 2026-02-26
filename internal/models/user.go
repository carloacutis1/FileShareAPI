package models

// UserModel represents a [description of your domain object]

type UserModel struct {
	ID        int    `json:"id"`   // Unique identifier
	Name      string `json:"name"` // Add your fields here
	Email     string `json:"email"`
	CreatedAt string `json:"createdat"`
	// Add more fields as needed
}

var userModels = []UserModel{
	{ID: 1, Name: "John Doe"},
	{ID: 2, Name: "Jane Doe"},
	{ID: 3, Name: "John Programming"},
	{ID: 4, Name: "Noname the Barbarian"},
	{ID: 5, Name: "Just Sayin"},
	// Have 3-5 sample items for testing
}

func GetAllUserModels() []UserModel {
	return userModels
}
