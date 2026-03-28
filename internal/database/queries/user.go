package queries

import (
	"gitlab.franciscan.edu/spr-26-backend/basement/internal/database"
	"gitlab.franciscan.edu/spr-26-backend/basement/internal/models"
)

// CreateUser inserts a new user record into the database.
// created_at is omitted so the database default (CURRENT_TIMESTAMP) is used.
func CreateUser(u models.UserModel) error {
	query := `INSERT INTO users (id, username, email, dob) VALUES ($1, $2, $3, $4)`
	_, err := database.DB.Exec(query, u.ID, u.Name, u.Email, u.DOB)
	return err
}

func GetUsers() ([]models.UserModel, error) {
	// Run database query
	rows, err := database.DB.Query("SELECT id, username, email, CAST(dob AS TEXT), CAST(created_at AS TEXT) FROM users")
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var userModel []models.UserModel
	for rows.Next() {
		var us models.UserModel
		if err := rows.Scan(&us.ID, &us.Name, &us.Email, &us.DOB, &us.CreatedAt); err != nil {
			return nil, err
		}
		userModel = append(userModel, us)
	}

	// Check for errors from row iteration
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return userModel, nil
}
