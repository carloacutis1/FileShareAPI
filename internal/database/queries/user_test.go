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
		ID:        "user_1",
		Name:      "John Doe",
		Email:     "jdoe001@programming.com",
		DOB:       "01/01/1900",
		CreatedAt: "2026-01-01 00:00:00",
	}
}

// ── GetUsers ──────────────────────────────────────────────────────────────────

func TestGetUsers_Success(t *testing.T) {
	mock := setupMock(t)
	u := sampleUser()

	cols := []string{"id", "username", "email", "dob", "created_at"}
	rows := sqlmock.NewRows(cols).AddRow(u.ID, u.Name, u.Email, u.DOB, u.CreatedAt)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, username, email`)).
		WillReturnRows(rows)

	users, err := GetUsers()
	if err != nil {
		t.Fatalf("GetUsers returned unexpected error: %v", err)
	}
	if len(users) != 1 {
		t.Fatalf("expected 1 user, got %d", len(users))
	}
	got := users[0]
	if got.ID != u.ID {
		t.Errorf("ID: got %q, want %q", got.ID, u.ID)
	}
	if got.Name != u.Name {
		t.Errorf("Name: got %q, want %q", got.Name, u.Name)
	}
	if got.Email != u.Email {
		t.Errorf("Email: got %q, want %q", got.Email, u.Email)
	}
}

func TestGetUsers_Empty(t *testing.T) {
	mock := setupMock(t)

	cols := []string{"id", "username", "email", "dob", "created_at"}
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, username, email`)).
		WillReturnRows(sqlmock.NewRows(cols))

	users, err := GetUsers()
	if err != nil {
		t.Fatalf("GetUsers returned unexpected error: %v", err)
	}
	if len(users) != 0 {
		t.Fatalf("expected 0 users, got %d", len(users))
	}
}

func TestGetUsers_DBError(t *testing.T) {
	mock := setupMock(t)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, username, email`)).
		WillReturnError(sql.ErrConnDone)

	if _, err := GetUsers(); err == nil {
		t.Fatal("GetUsers expected an error but got nil")
	}
}

// ── UpdateUser ────────────────────────────────────────────────────────────────

func TestUpdateUser_Success(t *testing.T) {
	mock := setupMock(t)
	u := sampleUser()

	mock.ExpectExec(regexp.QuoteMeta(`UPDATE users SET`)).
		WithArgs(u.Name, u.Email, u.DOB, u.ID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	if err := UpdateUser(u.ID, u); err != nil {
		t.Fatalf("UpdateUser returned unexpected error: %v", err)
	}
}

func TestUpdateUser_NotFound(t *testing.T) {
	mock := setupMock(t)
	u := sampleUser()

	mock.ExpectExec(regexp.QuoteMeta(`UPDATE users SET`)).
		WithArgs(u.Name, u.Email, u.DOB, "nonexistent").
		WillReturnResult(sqlmock.NewResult(0, 0))

	if err := UpdateUser("nonexistent", u); err != sql.ErrNoRows {
		t.Fatalf("UpdateUser: got %v, want sql.ErrNoRows", err)
	}
}

func TestUpdateUser_DBError(t *testing.T) {
	mock := setupMock(t)
	u := sampleUser()

	mock.ExpectExec(regexp.QuoteMeta(`UPDATE users SET`)).
		WithArgs(u.Name, u.Email, u.DOB, u.ID).
		WillReturnError(sql.ErrConnDone)

	if err := UpdateUser(u.ID, u); err == nil {
		t.Fatal("UpdateUser expected an error but got nil")
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
