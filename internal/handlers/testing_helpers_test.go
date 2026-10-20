package handlers

import (
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/carloacutis1/FileShareAPI/internal/database"
	"github.com/gin-gonic/gin"
)

func init() {
	gin.SetMode(gin.TestMode)
}

// setupDB replaces database.DB with a sqlmock instance and returns
// the mock so each test can set expectations on it.
func setupDB(t *testing.T) sqlmock.Sqlmock {
	t.Helper()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New: %v", err)
	}
	t.Cleanup(func() {
		_ = db.Close()
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unfulfilled mock expectations: %v", err)
		}
	})
	database.DB = db
	return mock
}
