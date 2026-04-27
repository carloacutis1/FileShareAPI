package queries

import (
	"database/sql"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestSearchItems(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	userID := "user_123"
	searchQuery := "test"
	now := time.Now()

	t.Run("Success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "name", "item_type", "created_at", "user_id"}).
			AddRow("f1", "test_file.txt", "file", now, userID).
			AddRow("d1", "test_folder", "folder", now, userID)

		// Match the UNION ALL query
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, name, 'file' as item_type, uploaded_at as created_at, uploaded_by as user_id`)).
			WithArgs("%"+searchQuery+"%", userID).
			WillReturnRows(rows)

		results, err := SearchItems(db, searchQuery, userID)

		if err != nil {
			t.Errorf("Expected no error, but got %v", err)
		}
		if len(results) != 2 {
			t.Errorf("Expected 2 results, but got %d", len(results))
		}

		if results[0].ID != "f1" {
			t.Errorf("Expected ID f1, but got %s", results[0].ID)
		}
		if results[0].Type != "file" {
			t.Errorf("Expected Type file, but got %s", results[0].Type)
		}

		if results[1].ID != "d1" {
			t.Errorf("Expected ID d1, but got %s", results[1].ID)
		}
		if results[1].Type != "folder" {
			t.Errorf("Expected Type folder, but got %s", results[1].Type)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	t.Run("NoResults", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, name, 'file' as item_type`)).
			WithArgs("%"+searchQuery+"%", userID).
			WillReturnRows(sqlmock.NewRows([]string{"id", "name", "item_type", "created_at", "user_id"}))

		results, err := SearchItems(db, searchQuery, userID)

		if err != nil {
			t.Errorf("Expected no error, but got %v", err)
		}
		if len(results) != 0 {
			t.Errorf("Expected 0 results, but got %d", len(results))
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	t.Run("DatabaseError", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, name, 'file' as item_type`)).
			WithArgs("%"+searchQuery+"%", userID).
			WillReturnError(sql.ErrConnDone)

		results, err := SearchItems(db, searchQuery, userID)

		if err == nil {
			t.Error("Expected error, but got nil")
		}
		if results != nil {
			t.Errorf("Expected nil results, but got %v", results)
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})
}
