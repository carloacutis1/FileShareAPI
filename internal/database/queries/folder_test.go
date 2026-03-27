package queries

import (
	"database/sql"
	"regexp"
	"testing"
	"time"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"gitlab.franciscan.edu/spr-26-backend/basement/internal/models"
)

func sampleFolder() models.Folder {
	return models.Folder{
		ID:             "fld-1",
		ParentFolderID: "fld-parent",
		Name:           "Work",
		UserID:         "user-1",
	}
}

func sampleFolderTime() time.Time {
	return time.Date(2026, 3, 1, 12, 0, 0, 0, time.UTC)
}

var folderCols = []string{"id", "parent_folder_id", "name", "created_at", "user_id"}

// ── CreateFolder ──────────────────────────────────────────────────────────────

func TestCreateFolder_Success(t *testing.T) {
	mock := setupMock(t)
	f := sampleFolder()

	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO folders`)).
		WithArgs(f.ID, f.ParentFolderID, f.Name, f.UserID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	if err := CreateFolder(f); err != nil {
		t.Fatalf("CreateFolder returned unexpected error: %v", err)
	}
}

func TestCreateFolder_RootParentUsesNil(t *testing.T) {
	mock := setupMock(t)
	f := models.Folder{
		ID:             "fld-root",
		ParentFolderID: "",
		Name:           "Home",
		UserID:         "user-1",
	}

	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO folders`)).
		WithArgs(f.ID, nil, f.Name, f.UserID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	if err := CreateFolder(f); err != nil {
		t.Fatalf("CreateFolder returned unexpected error: %v", err)
	}
}

func TestCreateFolder_DBError(t *testing.T) {
	mock := setupMock(t)
	f := sampleFolder()

	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO folders`)).
		WithArgs(f.ID, f.ParentFolderID, f.Name, f.UserID).
		WillReturnError(sql.ErrConnDone)

	if err := CreateFolder(f); err == nil {
		t.Fatal("CreateFolder expected an error but got nil")
	}
}

// ── GetAllFolders ─────────────────────────────────────────────────────────────

func TestGetAllFolders_Success(t *testing.T) {
	mock := setupMock(t)
	f := sampleFolder()
	ts := sampleFolderTime()

	rows := sqlmock.NewRows(folderCols).AddRow(
		f.ID, f.ParentFolderID, f.Name, ts, f.UserID,
	)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, parent_folder_id, name, created_at, user_id FROM folders ORDER BY created_at DESC`)).
		WillReturnRows(rows)

	folders, err := GetAllFolders()
	if err != nil {
		t.Fatalf("GetAllFolders returned unexpected error: %v", err)
	}
	if len(folders) != 1 {
		t.Fatalf("expected 1 folder, got %d", len(folders))
	}
	got := folders[0]
	if got.ID != f.ID {
		t.Errorf("ID: got %q, want %q", got.ID, f.ID)
	}
	if got.ParentFolderID != f.ParentFolderID {
		t.Errorf("ParentFolderID: got %q, want %q", got.ParentFolderID, f.ParentFolderID)
	}
	if got.Name != f.Name {
		t.Errorf("Name: got %q, want %q", got.Name, f.Name)
	}
	if !got.CreatedAt.Equal(ts) {
		t.Errorf("CreatedAt: got %v, want %v", got.CreatedAt, ts)
	}
	if got.UserID != f.UserID {
		t.Errorf("UserID: got %q, want %q", got.UserID, f.UserID)
	}
}

func TestGetAllFolders_NullParentScansToEmptyString(t *testing.T) {
	mock := setupMock(t)
	ts := sampleFolderTime()

	rows := sqlmock.NewRows(folderCols).AddRow(
		"fld_home", nil, "Home", ts, "user_1",
	)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, parent_folder_id, name, created_at, user_id FROM folders ORDER BY created_at DESC`)).
		WillReturnRows(rows)

	folders, err := GetAllFolders()
	if err != nil {
		t.Fatalf("GetAllFolders returned unexpected error: %v", err)
	}
	if len(folders) != 1 {
		t.Fatalf("expected 1 folder, got %d", len(folders))
	}
	if folders[0].ParentFolderID != "" {
		t.Errorf("ParentFolderID: got %q, want empty string for NULL", folders[0].ParentFolderID)
	}
}

func TestGetAllFolders_Empty(t *testing.T) {
	mock := setupMock(t)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, parent_folder_id, name, created_at, user_id FROM folders ORDER BY created_at DESC`)).
		WillReturnRows(sqlmock.NewRows(folderCols))

	folders, err := GetAllFolders()
	if err != nil {
		t.Fatalf("GetAllFolders returned unexpected error: %v", err)
	}
	if len(folders) != 0 {
		t.Fatalf("expected 0 folders, got %d", len(folders))
	}
}

func TestGetAllFolders_DBError(t *testing.T) {
	mock := setupMock(t)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, parent_folder_id, name, created_at, user_id FROM folders ORDER BY created_at DESC`)).
		WillReturnError(sql.ErrConnDone)

	if _, err := GetAllFolders(); err == nil {
		t.Fatal("GetAllFolders expected an error but got nil")
	}
}

// ── GetFolderByID ─────────────────────────────────────────────────────────────

func TestGetFolderByID_Success(t *testing.T) {
	mock := setupMock(t)
	f := sampleFolder()
	ts := sampleFolderTime()

	rows := sqlmock.NewRows(folderCols).AddRow(
		f.ID, f.ParentFolderID, f.Name, ts, f.UserID,
	)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, parent_folder_id, name, created_at, user_id FROM folders WHERE id = $1`)).
		WithArgs("fld-1").
		WillReturnRows(rows)

	folder, err := GetFolderByID("fld-1")
	if err != nil {
		t.Fatalf("GetFolderByID returned unexpected error: %v", err)
	}
	if folder.ID != f.ID {
		t.Errorf("ID: got %q, want %q", folder.ID, f.ID)
	}
	if folder.ParentFolderID != f.ParentFolderID {
		t.Errorf("ParentFolderID: got %q, want %q", folder.ParentFolderID, f.ParentFolderID)
	}
}

func TestGetFolderByID_NotFound(t *testing.T) {
	mock := setupMock(t)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, parent_folder_id, name, created_at, user_id FROM folders WHERE id = $1`)).
		WithArgs("nonexistent").
		WillReturnError(sql.ErrNoRows)

	_, err := GetFolderByID("nonexistent")
	if err != sql.ErrNoRows {
		t.Fatalf("expected sql.ErrNoRows, got %v", err)
	}
}

// ── UpdateFolder ──────────────────────────────────────────────────────────────

func TestUpdateFolder_Success(t *testing.T) {
	mock := setupMock(t)
	f := sampleFolder()

	mock.ExpectExec(regexp.QuoteMeta(`UPDATE folders SET`)).
		WithArgs(f.ParentFolderID, f.Name, f.UserID, "fld-1").
		WillReturnResult(sqlmock.NewResult(0, 1))

	if err := UpdateFolder("fld-1", f); err != nil {
		t.Fatalf("UpdateFolder returned unexpected error: %v", err)
	}
}

func TestUpdateFolder_RootParentUsesNil(t *testing.T) {
	mock := setupMock(t)
	f := models.Folder{
		ParentFolderID: "",
		Name:           "Renamed",
		UserID:         "user-1",
	}

	mock.ExpectExec(regexp.QuoteMeta(`UPDATE folders SET`)).
		WithArgs(nil, f.Name, f.UserID, "fld-1").
		WillReturnResult(sqlmock.NewResult(0, 1))

	if err := UpdateFolder("fld-1", f); err != nil {
		t.Fatalf("UpdateFolder returned unexpected error: %v", err)
	}
}

func TestUpdateFolder_NotFound(t *testing.T) {
	mock := setupMock(t)
	f := sampleFolder()

	mock.ExpectExec(regexp.QuoteMeta(`UPDATE folders SET`)).
		WithArgs(f.ParentFolderID, f.Name, f.UserID, "nonexistent").
		WillReturnResult(sqlmock.NewResult(0, 0))

	if err := UpdateFolder("nonexistent", f); err != sql.ErrNoRows {
		t.Fatalf("UpdateFolder: got %v, want sql.ErrNoRows", err)
	}
}

func TestUpdateFolder_DBError(t *testing.T) {
	mock := setupMock(t)
	f := sampleFolder()

	mock.ExpectExec(regexp.QuoteMeta(`UPDATE folders SET`)).
		WithArgs(f.ParentFolderID, f.Name, f.UserID, f.ID).
		WillReturnError(sql.ErrConnDone)

	if err := UpdateFolder(f.ID, f); err == nil {
		t.Fatal("UpdateFolder expected an error but got nil")
	}
}

// ── DeleteFolder ──────────────────────────────────────────────────────────────

func TestDeleteFolder_Success(t *testing.T) {
	mock := setupMock(t)

	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM folders WHERE id = $1`)).
		WithArgs("fld-1").
		WillReturnResult(sqlmock.NewResult(0, 1))

	if err := DeleteFolder("fld-1"); err != nil {
		t.Fatalf("DeleteFolder returned unexpected error: %v", err)
	}
}

func TestDeleteFolder_NotFound(t *testing.T) {
	mock := setupMock(t)

	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM folders WHERE id = $1`)).
		WithArgs("nonexistent").
		WillReturnResult(sqlmock.NewResult(0, 0))

	if err := DeleteFolder("nonexistent"); err != sql.ErrNoRows {
		t.Fatalf("DeleteFolder: got %v, want sql.ErrNoRows", err)
	}
}

func TestDeleteFolder_DBError(t *testing.T) {
	mock := setupMock(t)

	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM folders WHERE id = $1`)).
		WithArgs("fld-1").
		WillReturnError(sql.ErrConnDone)

	if err := DeleteFolder("fld-1"); err == nil {
		t.Fatal("DeleteFolder expected an error but got nil")
	}
}
