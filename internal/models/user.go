package models

// UserModel to store user data below
type UserModel struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	DOB       string `json:"dob"`
	CreatedAt string `json:"createdat"`
	// Add more fields as needed
}

var userModels = []UserModel{
	// User Name
	{ID: "user_1", Name: "John Doe"},
	{ID: "user_2", Name: "Jane Doe"},
	{ID: "user_3", Name: "John Programming"},
	{ID: "user_4", Name: "Noname the Barbarian"},
	{ID: "user_5", Name: "Just Sayin"},
	// Have 3-5 sample items for testing
	// Include data for other variables in this scope - MH
}

func GetAllUserModels() []UserModel {
	return userModels
}
