package queries

import (
	"database/sql"
	"gitlab.franciscan.edu/spr-26-backend/basement/internal/models"
)

// SearchItems searches both files and folders for a name matching the query
func SearchItems(db *sql.DB, query string, userID string) ([]models.SearchResult, error) {
	sqlQuery := `
		SELECT id, name, 'file' as item_type, uploaded_at as created_at, uploaded_by as user_id 
		FROM files 
		WHERE name ILIKE $1 AND uploaded_by = $2
		UNION ALL
		SELECT id, name, 'folder' as item_type, created_at, user_id 
		FROM folders 
		WHERE name ILIKE $1 AND user_id = $2
	`
	rows, err := db.Query(sqlQuery, "%"+query+"%", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []models.SearchResult
	for rows.Next() {
		var res models.SearchResult
		if err := rows.Scan(&res.ID, &res.Name, &res.Type, &res.CreatedAt, &res.UserID); err != nil {
			return nil, err
		}
		results = append(results, res)
	}

	return results, nil
}
