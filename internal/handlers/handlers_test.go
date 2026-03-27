package handlers

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"
	"time"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"gitlab.franciscan.edu/spr-26-backend/basement/internal/database"
)

func init() {
	gin.SetMode(gin.TestMode)
}

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

// ── Users ─────────────────────────────────────────────────────────────────────

func TestListUserModels_OK(t *testing.T) {
	router := gin.New()
	router.GET("/users", ListUserModels)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/users", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var result []interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &result); err != nil {
		t.Fatalf("body is not a JSON array: %v", err)
	}
	if len(result) == 0 {
		t.Fatal("expected at least one user in the response")
	}
}

// ── Folders ───────────────────────────────────────────────────────────────────

var folderCols = []string{"id", "parent_folder_id", "name", "created_at", "user_id"}

func TestListFolders_OK(t *testing.T) {
	mock := setupDB(t)
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, parent_folder_id, name, created_at, user_id FROM folders`)).
		WillReturnRows(sqlmock.NewRows(folderCols))

	router := gin.New()
	router.GET("/folders", ListFolders)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/folders", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
}

// ── Files ─────────────────────────────────────────────────────────────────────

var fileCols = []string{
	"id", "folder_id", "name", "size", "file_type",
	"uploaded_at", "description", "uploaded_by", "tags",
	"last_accessed", "checksum", "custom_attributes",
}

func TestListFiles_OK(t *testing.T) {
	mock := setupDB(t)
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, folder_id`)).
		WillReturnRows(sqlmock.NewRows(fileCols))

	router := gin.New()
	router.GET("/files", ListFiles)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/files", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestListFiles_DBError(t *testing.T) {
	mock := setupDB(t)
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, folder_id`)).
		WillReturnError(sql.ErrConnDone)

	router := gin.New()
	router.GET("/files", ListFiles)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/files", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", w.Code)
	}
}

func TestCreateFile_BadJSON(t *testing.T) {
	router := gin.New()
	router.POST("/files", CreateFile)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/files", bytes.NewBufferString("{bad}"))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestCreateFile_OK(t *testing.T) {
	mock := setupDB(t)
	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO files`)).
		WithArgs(
			"file-99", nil, "test.txt", 100, "",
			"", "user-1", "", time.Time{},
			"", nil,
		).
		WillReturnResult(sqlmock.NewResult(1, 1))

	body, _ := json.Marshal(map[string]interface{}{
		"id":          "file-99",
		"name":        "test.txt",
		"size":        100,
		"uploaded_by": "user-1",
	})

	router := gin.New()
	router.POST("/files", CreateFile)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/files", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", w.Code, w.Body.String())
	}
}

func TestCreateFile_DBError(t *testing.T) {
	mock := setupDB(t)
	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO files`)).
		WithArgs(
			"file-99", nil, "test.txt", 100, "",
			"", "user-1", "", time.Time{},
			"", nil,
		).
		WillReturnError(sql.ErrConnDone)

	body, _ := json.Marshal(map[string]interface{}{
		"id":          "file-99",
		"name":        "test.txt",
		"size":        100,
		"uploaded_by": "user-1",
	})

	router := gin.New()
	router.POST("/files", CreateFile)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/files", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", w.Code)
	}
}

func TestUpdateFile_BadJSON(t *testing.T) {
	router := gin.New()
	router.PUT("/files/:id", UpdateFile)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPut, "/files/file-1", bytes.NewBufferString("{bad}"))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestUpdateFile_OK(t *testing.T) {
	mock := setupDB(t)
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE files SET`)).
		WithArgs(
			nil, "updated.txt", 200, "", "",
			"", time.Time{}, "", nil,
			"file-1",
		).
		WillReturnResult(sqlmock.NewResult(0, 1))

	body, _ := json.Marshal(map[string]interface{}{"name": "updated.txt", "size": 200})

	router := gin.New()
	router.PUT("/files/:id", UpdateFile)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPut, "/files/file-1", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestUpdateFile_NotFound(t *testing.T) {
	mock := setupDB(t)
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE files SET`)).
		WithArgs(
			nil, "updated.txt", 0, "", "",
			"", time.Time{}, "", nil,
			"file-999",
		).
		WillReturnResult(sqlmock.NewResult(0, 0))

	body, _ := json.Marshal(map[string]interface{}{"name": "updated.txt"})

	router := gin.New()
	router.PUT("/files/:id", UpdateFile)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPut, "/files/file-999", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}

func TestDeleteFile_OK(t *testing.T) {
	mock := setupDB(t)
	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM files WHERE id = $1`)).
		WithArgs("file-1").
		WillReturnResult(sqlmock.NewResult(0, 1))

	router := gin.New()
	router.DELETE("/files/:id", DeleteFile)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodDelete, "/files/file-1", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestDeleteFile_NotFound(t *testing.T) {
	mock := setupDB(t)
	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM files WHERE id = $1`)).
		WithArgs("file-999").
		WillReturnResult(sqlmock.NewResult(0, 0))

	router := gin.New()
	router.DELETE("/files/:id", DeleteFile)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodDelete, "/files/file-999", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}

// ── Shared Files ──────────────────────────────────────────────────────────────

var sharedFileCols = []string{"id", "file_id", "owner_id", "shared_with", "permissions"}

func TestListSharedFiles_OK(t *testing.T) {
	mock := setupDB(t)
	mock.ExpectQuery("SELECT id, file_id, owner_id, shared_with, permissions FROM shared_files").
		WillReturnRows(sqlmock.NewRows(sharedFileCols))

	router := gin.New()
	router.GET("/shared-files", ListSharedFiles)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/shared-files", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestListSharedFiles_DBError(t *testing.T) {
	mock := setupDB(t)
	mock.ExpectQuery("SELECT id, file_id, owner_id, shared_with, permissions FROM shared_files").
		WillReturnError(sql.ErrConnDone)

	router := gin.New()
	router.GET("/shared-files", ListSharedFiles)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/shared-files", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", w.Code)
	}
}
