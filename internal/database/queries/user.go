package queries

import (
	"gitlab.franciscan.edu/spr-26-backend/basement/internal/database"
	"gitlab.franciscan.edu/spr-26-backend/basement/internal/models"
)

func GetUsers() ([]models.UserModel, error) {
	// Run database query
	rows, err := database.DB.Query("SELECT")
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var UserModel []models.UserModel
	for rows.Next() {
		var us models.UserModel
		if err := rows.Scan(&us.ID, &us.Name, &us.Email, &us.DOB, &us.CreatedAt); err != nil {
			return nil, err
		}
		UserModel = append(UserModel, us)
	}

	// Check for errors from row iteration
	if err := rows.Err() != nil {
		return nil, err
	}
	return UserModel, nil
}
