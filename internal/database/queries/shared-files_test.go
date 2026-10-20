package queries

import (
	"database/sql"
	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/carloacutis1/FileShareAPI/internal/models"
	"regexp"
	"testing"
)

// Helper function to create sample shared file
func sampleSharedFile() models.SharedFile {
	return models.SharedFile{
		ID:          "sf_1",
		FileID:      "file_1",
		OwnerID:     "user_1",
		SharedWith:  "user_2",
		Permissions: "read",
	}
}

// GET unit test success
// mock db used
func TestGetSharedFiles_Success(t *testing.T) {
	mock := setupMock(t)
	sf := sampleSharedFile()

	cols := []string{"id", "file_id", "owner_id", "shared_with", "permissions"}
	rows := sqlmock.NewRows(cols).
		AddRow(sf.ID, sf.FileID, sf.OwnerID, sf.SharedWith, sf.Permissions).
		AddRow("sf_2", "file_2", "user_1", "user_3", "write")

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, file_id, owner_id, shared_with, permissions FROM shared_files`)).
		WillReturnRows(rows)

	sharedFiles, err := GetSharedFiles()
	if err != nil {
		t.Fatalf("GetSharedFiles returned unexpected error: %v", err)
	}
	if len(sharedFiles) != 2 {
		t.Fatalf("expected 2 shared files, got %d", len(sharedFiles))
	}

	// Validate first shared file
	got := sharedFiles[0]
	if got.ID != sf.ID {
		t.Errorf("ID: got %q, want %q", got.ID, sf.ID)
	}
	if got.FileID != sf.FileID {
		t.Errorf("FileID: got %q, want %q", got.FileID, sf.FileID)
	}
	if got.OwnerID != sf.OwnerID {
		t.Errorf("OwnerID: got %q, want %q", got.OwnerID, sf.OwnerID)
	}
	if got.SharedWith != sf.SharedWith {
		t.Errorf("SharedWith: got %q, want %q", got.SharedWith, sf.SharedWith)
	}
	if got.Permissions != sf.Permissions {
		t.Errorf("Permissions: got %q, want %q", got.Permissions, sf.Permissions)
	}
}

// GET unit test empty
func TestGetSharedFiles_Empty(t *testing.T) {
	mock := setupMock(t)

	cols := []string{"id", "file_id", "owner_id", "shared_with", "permissions"}
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, file_id, owner_id, shared_with, permissions FROM shared_files`)).
		WillReturnRows(sqlmock.NewRows(cols))

	sharedFiles, err := GetSharedFiles()
	if err != nil {
		t.Fatalf("GetSharedFiles returned unexpected error: %v", err)
	}
	if len(sharedFiles) != 0 {
		t.Fatalf("expected 0 shared files, got %d", len(sharedFiles))
	}
}

// GET unit test failure
func TestGetSharedFiles_DBError(t *testing.T) {
	mock := setupMock(t)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, file_id, owner_id, shared_with, permissions FROM shared_files`)).
		WillReturnError(sql.ErrConnDone)

	sharedFiles, err := GetSharedFiles()
	if err == nil {
		t.Fatal("GetSharedFiles expected an error but got nil")
	}
	if sharedFiles != nil {
		t.Fatalf("expected nil slice on error, got %v", sharedFiles)
	}
}

// ── CreateSharedFile ──────────────────────────────────────────────────────────

func TestCreateSharedFile_Success(t *testing.T) {
	mock := setupMock(t)
	sf := sampleSharedFile()

	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO shared_files`)).
		WithArgs(sf.ID, sf.FileID, sf.OwnerID, sf.SharedWith, sf.Permissions).
		WillReturnResult(sqlmock.NewResult(1, 1))

	if err := CreateSharedFile(sf); err != nil {
		t.Fatalf("CreateSharedFile returned unexpected error: %v", err)
	}
}

func TestCreateSharedFile_DBError(t *testing.T) {
	mock := setupMock(t)
	sf := sampleSharedFile()

	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO shared_files`)).
		WithArgs(sf.ID, sf.FileID, sf.OwnerID, sf.SharedWith, sf.Permissions).
		WillReturnError(sql.ErrConnDone)

	if err := CreateSharedFile(sf); err == nil {
		t.Fatal("CreateSharedFile expected an error but got nil")
	}
}

// ── DeleteSharedFile ──────────────────────────────────────────────────────────

func TestDeleteSharedFile_Success(t *testing.T) {
	mock := setupMock(t)

	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM shared_files WHERE id = $1`)).
		WithArgs("sf_1").
		WillReturnResult(sqlmock.NewResult(0, 1))

	if err := DeleteSharedFile("sf_1"); err != nil {
		t.Fatalf("DeleteSharedFile returned unexpected error: %v", err)
	}
}

func TestDeleteSharedFile_NotFound(t *testing.T) {
	mock := setupMock(t)

	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM shared_files WHERE id = $1`)).
		WithArgs("nonexistent").
		WillReturnResult(sqlmock.NewResult(0, 0))

	if err := DeleteSharedFile("nonexistent"); err != sql.ErrNoRows {
		t.Fatalf("DeleteSharedFile: got %v, want sql.ErrNoRows", err)
	}
}

func TestDeleteSharedFile_DBError(t *testing.T) {
	mock := setupMock(t)

	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM shared_files WHERE id = $1`)).
		WithArgs("sf_1").
		WillReturnError(sql.ErrConnDone)

	if err := DeleteSharedFile("sf_1"); err == nil {
		t.Fatal("DeleteSharedFile expected an error but got nil")
	}
}

// ── UpdateSharedFile ──────────────────────────────────────────────────────────

func TestUpdateSharedFile_Success(t *testing.T) {
	mock := setupMock(t)
	sf := sampleSharedFile()

	mock.ExpectExec(regexp.QuoteMeta(`UPDATE shared_files SET`)).
		WithArgs(sf.FileID, sf.OwnerID, sf.SharedWith, sf.Permissions, "sf_1").
		WillReturnResult(sqlmock.NewResult(0, 1))

	if err := UpdateSharedFile("sf_1", sf); err != nil {
		t.Fatalf("UpdateSharedFile returned unexpected error: %v", err)
	}
}

func TestUpdateSharedFile_NotFound(t *testing.T) {
	mock := setupMock(t)
	sf := sampleSharedFile()

	mock.ExpectExec(regexp.QuoteMeta(`UPDATE shared_files SET`)).
		WithArgs(sf.FileID, sf.OwnerID, sf.SharedWith, sf.Permissions, "nonexistent").
		WillReturnResult(sqlmock.NewResult(0, 0))

	if err := UpdateSharedFile("nonexistent", sf); err != sql.ErrNoRows {
		t.Fatalf("UpdateSharedFile: got %v, want sql.ErrNoRows", err)
	}
}

func TestUpdateSharedFile_DBError(t *testing.T) {
	mock := setupMock(t)
	sf := sampleSharedFile()

	mock.ExpectExec(regexp.QuoteMeta(`UPDATE shared_files SET`)).
		WithArgs(sf.FileID, sf.OwnerID, sf.SharedWith, sf.Permissions, sf.ID).
		WillReturnError(sql.ErrConnDone)

	if err := UpdateSharedFile(sf.ID, sf); err == nil {
		t.Fatal("UpdateSharedFile expected an error but got nil")
	}
}
