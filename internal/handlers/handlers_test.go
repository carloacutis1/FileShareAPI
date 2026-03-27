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
)


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
