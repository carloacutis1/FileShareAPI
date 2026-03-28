package queries

import (
	"database/sql"
	"regexp"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"gitlab.franciscan.edu/spr-26-backend/basement/internal/models"
)

func sampleUser() models.UserModel {
	return models.UserModel{
		ID:    "user_1",
		Name:  "John Doe",
		Email: "jdoe001@programming.com",
		DOB:   "01/01/1900",
	}
}

// ── CreateUser ────────────────────────────────────────────────────────────────

func TestCreateUser_Success(t *testing.T) {
	mock := setupMock(t)
	u := sampleUser()

	// Expect the INSERT with the four user fields in order
	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO users (id, username, email, dob)`)).
		WithArgs(u.ID, u.Name, u.Email, u.DOB).
		WillReturnResult(sqlmock.NewResult(1, 1))

	if err := CreateUser(u); err != nil {
		t.Fatalf("CreateUser returned unexpected error: %v", err)
	}
}

func TestCreateUser_DBError(t *testing.T) {
	mock := setupMock(t)
	u := sampleUser()

	// Simulate a connection failure
	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO users (id, username, email, dob)`)).
		WithArgs(u.ID, u.Name, u.Email, u.DOB).
		WillReturnError(sql.ErrConnDone)

	if err := CreateUser(u); err == nil {
		t.Fatal("CreateUser expected an error but got nil")
	}
}
