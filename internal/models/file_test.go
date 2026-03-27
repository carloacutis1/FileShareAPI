package models

import (
	"database/sql"
	"regexp"
	"testing"
	"time"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"gitlab.franciscan.edu/spr-26-backend/basement/internal/database"
)

// setupMock replaces database.DB with a sqlmock instance and returns
// the mock so each test can set expectations on it.
func setupMock(t *testing.T) sqlmock.Sqlmock {
	t.Helper()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New: %v", err)
	}
	t.Cleanup(func() {
		db.Close()
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unfulfilled mock expectations: %v", err)
		}
	})
	database.DB = db
	return mock
}

func sampleFile() File {
	return File{
		ID:               "file-1",
		FolderID:         "folder-1",
		Name:             "report.pdf",
		Size:             2048,
		FileType:         "application/pdf",
		UploadedAt:       time.Time{},
		Description:      "Annual report",
		UploadedBy:       "user-1",
		Tags:             "report,annual",
		LastAccessed:     time.Time{},
		Checksum:         "abc123",
		CustomAttributes: `{"confidential":true}`,
	}
}

// ── CreateFile ────────────────────────────────────────────────────────────────

func TestCreateFile_Success(t *testing.T) {
	mock := setupMock(t)
	f := sampleFile()

	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO files`)).
		WithArgs(
			f.ID, f.FolderID, f.Name, f.Size, f.FileType,
			f.Description, f.UploadedBy, f.Tags, f.LastAccessed,
			f.Checksum, f.CustomAttributes,
		).
		WillReturnResult(sqlmock.NewResult(1, 1))

	if err := CreateFile(f); err != nil {
		t.Fatalf("CreateFile returned unexpected error: %v", err)
	}
}

func TestCreateFile_DBError(t *testing.T) {
	mock := setupMock(t)
	f := sampleFile()

	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO files`)).
		WithArgs(
			f.ID, f.FolderID, f.Name, f.Size, f.FileType,
			f.Description, f.UploadedBy, f.Tags, f.LastAccessed,
			f.Checksum, f.CustomAttributes,
		).
		WillReturnError(sql.ErrConnDone)

	if err := CreateFile(f); err == nil {
		t.Fatal("CreateFile expected an error but got nil")
	}
}

// ── GetFiles ──────────────────────────────────────────────────────────────────

func TestGetFiles_Success(t *testing.T) {
	mock := setupMock(t)
	f := sampleFile()

	cols := []string{
		"id", "folder_id", "name", "size", "file_type",
		"uploaded_at", "description", "uploaded_by", "tags",
		"last_accessed", "checksum", "custom_attributes",
	}
	rows := sqlmock.NewRows(cols).AddRow(
		f.ID, f.FolderID, f.Name, f.Size, f.FileType,
		f.UploadedAt, f.Description, f.UploadedBy, f.Tags,
		f.LastAccessed, f.Checksum, f.CustomAttributes,
	)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, folder_id`)).
		WillReturnRows(rows)

	files, err := GetFiles()
	if err != nil {
		t.Fatalf("GetFiles returned unexpected error: %v", err)
	}
	if len(files) != 1 {
		t.Fatalf("expected 1 file, got %d", len(files))
	}
	got := files[0]
	if got.ID != f.ID {
		t.Errorf("ID: got %q, want %q", got.ID, f.ID)
	}
	if got.Name != f.Name {
		t.Errorf("Name: got %q, want %q", got.Name, f.Name)
	}
	if got.Size != f.Size {
		t.Errorf("Size: got %d, want %d", got.Size, f.Size)
	}
}

func TestGetFiles_Empty(t *testing.T) {
	mock := setupMock(t)

	cols := []string{
		"id", "folder_id", "name", "size", "file_type",
		"uploaded_at", "description", "uploaded_by", "tags",
		"last_accessed", "checksum", "custom_attributes",
	}
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, folder_id`)).
		WillReturnRows(sqlmock.NewRows(cols))

	files, err := GetFiles()
	if err != nil {
		t.Fatalf("GetFiles returned unexpected error: %v", err)
	}
	if len(files) != 0 {
		t.Fatalf("expected 0 files, got %d", len(files))
	}
}

func TestGetFiles_DBError(t *testing.T) {
	mock := setupMock(t)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, folder_id`)).
		WillReturnError(sql.ErrConnDone)

	if _, err := GetFiles(); err == nil {
		t.Fatal("GetFiles expected an error but got nil")
	}
}

// ── UpdateFile ────────────────────────────────────────────────────────────────

func TestUpdateFile_Success(t *testing.T) {
	mock := setupMock(t)
	f := sampleFile()

	mock.ExpectExec(regexp.QuoteMeta(`UPDATE files SET`)).
		WithArgs(
			f.FolderID, f.Name, f.Size, f.FileType, f.Description,
			f.Tags, f.LastAccessed, f.Checksum, f.CustomAttributes,
			f.ID,
		).
		WillReturnResult(sqlmock.NewResult(0, 1))

	if err := UpdateFile(f); err != nil {
		t.Fatalf("UpdateFile returned unexpected error: %v", err)
	}
}

func TestUpdateFile_NotFound(t *testing.T) {
	mock := setupMock(t)
	f := sampleFile()

	mock.ExpectExec(regexp.QuoteMeta(`UPDATE files SET`)).
		WithArgs(
			f.FolderID, f.Name, f.Size, f.FileType, f.Description,
			f.Tags, f.LastAccessed, f.Checksum, f.CustomAttributes,
			f.ID,
		).
		WillReturnResult(sqlmock.NewResult(0, 0))

	if err := UpdateFile(f); err != sql.ErrNoRows {
		t.Fatalf("UpdateFile: got %v, want sql.ErrNoRows", err)
	}
}

func TestUpdateFile_DBError(t *testing.T) {
	mock := setupMock(t)
	f := sampleFile()

	mock.ExpectExec(regexp.QuoteMeta(`UPDATE files SET`)).
		WithArgs(
			f.FolderID, f.Name, f.Size, f.FileType, f.Description,
			f.Tags, f.LastAccessed, f.Checksum, f.CustomAttributes,
			f.ID,
		).
		WillReturnError(sql.ErrConnDone)

	if err := UpdateFile(f); err == nil {
		t.Fatal("UpdateFile expected an error but got nil")
	}
}

// ── DeleteFile ────────────────────────────────────────────────────────────────

func TestDeleteFile_Success(t *testing.T) {
	mock := setupMock(t)

	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM files WHERE id = $1`)).
		WithArgs("file-1").
		WillReturnResult(sqlmock.NewResult(0, 1))

	if err := DeleteFile("file-1"); err != nil {
		t.Fatalf("DeleteFile returned unexpected error: %v", err)
	}
}

func TestDeleteFile_NotFound(t *testing.T) {
	mock := setupMock(t)

	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM files WHERE id = $1`)).
		WithArgs("nonexistent").
		WillReturnResult(sqlmock.NewResult(0, 0))

	if err := DeleteFile("nonexistent"); err != sql.ErrNoRows {
		t.Fatalf("DeleteFile: got %v, want sql.ErrNoRows", err)
	}
}

func TestDeleteFile_DBError(t *testing.T) {
	mock := setupMock(t)

	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM files WHERE id = $1`)).
		WithArgs("file-1").
		WillReturnError(sql.ErrConnDone)

	if err := DeleteFile("file-1"); err == nil {
		t.Fatal("DeleteFile expected an error but got nil")
	}
}
