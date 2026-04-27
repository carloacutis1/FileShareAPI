package queries

import (
	"database/sql"

	"gitlab.franciscan.edu/spr-26-backend/basement/internal/database"
	"gitlab.franciscan.edu/spr-26-backend/basement/internal/models"
)

// CreateUser inserts a new user record into the database.
// created_at is omitted so the database default (CURRENT_TIMESTAMP) is used.
func CreateUser(u models.UserModel) error {
	query := `INSERT INTO users (id, username, email, dob, password_hash) VALUES ($1, $2, $3, $4, $5)`
	_, err := database.DB.Exec(query, u.ID, u.Name, u.Email, u.DOB, u.PasswordHash)
	return err
}

// DeleteUser removes a user record from the database by ID.
func DeleteUser(id string) error {
	query := `DELETE FROM users WHERE id = $1`
	result, err := database.DB.Exec(query, id)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}

// UpdateUser updates an existing user's username, email, and dob by ID.
func UpdateUser(id string, u models.UserModel) error {
	query := `UPDATE users SET username = $1, email = $2, dob = $3 WHERE id = $4`
	result, err := database.DB.Exec(query, u.Name, u.Email, u.DOB, id)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return sql.ErrNoRows
	}
	return nil
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
