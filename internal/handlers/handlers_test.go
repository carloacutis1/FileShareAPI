package handlers

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"
	"time"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"gitlab.franciscan.edu/spr-26-backend/basement/internal/storage"
)

// init and setupDB are declared in testing_helpers_test.go and shared across this package.
// Due to a redeclaration of setupDB causing build errors, it was deleted here and kept in testing_helpers.go (now testing_helpers_test.go). - MH

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
	"last_accessed", "checksum", "storage_object_key", "custom_attributes",
}

func TestListFiles_OK(t *testing.T) {
	mock := setupDB(t)
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, folder_id`)).
		WillReturnRows(sqlmock.NewRows(fileCols))

	fh := NewFileHandler(storage.NewMemoryStore())
	router := gin.New()
	router.GET("/files", fh.ListFiles)

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

	fh := NewFileHandler(storage.NewMemoryStore())
	router := gin.New()
	router.GET("/files", fh.ListFiles)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/files", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", w.Code)
	}
}

func TestCreateFile_BadJSON(t *testing.T) {
	fh := NewFileHandler(storage.NewMemoryStore())
	router := gin.New()
	router.POST("/files", fh.CreateFile)

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
			"", nil, nil,
		).
		WillReturnResult(sqlmock.NewResult(1, 1))

	body, _ := json.Marshal(map[string]interface{}{
		"id":          "file-99",
		"name":        "test.txt",
		"size":        100,
		"uploaded_by": "user-1",
	})

	fh := NewFileHandler(storage.NewMemoryStore())
	router := gin.New()
	router.POST("/files", fh.CreateFile)

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
			"", nil, nil,
		).
		WillReturnError(sql.ErrConnDone)

	body, _ := json.Marshal(map[string]interface{}{
		"id":          "file-99",
		"name":        "test.txt",
		"size":        100,
		"uploaded_by": "user-1",
	})

	fh := NewFileHandler(storage.NewMemoryStore())
	router := gin.New()
	router.POST("/files", fh.CreateFile)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/files", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", w.Code)
	}
}

func TestUpdateFile_BadJSON(t *testing.T) {
	fh := NewFileHandler(storage.NewMemoryStore())
	router := gin.New()
	router.PUT("/files/:id", fh.UpdateFile)

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
			"", time.Time{}, "", nil, nil,
			"file-1",
		).
		WillReturnResult(sqlmock.NewResult(0, 1))

	body, _ := json.Marshal(map[string]interface{}{"name": "updated.txt", "size": 200})

	fh := NewFileHandler(storage.NewMemoryStore())
	router := gin.New()
	router.PUT("/files/:id", fh.UpdateFile)

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
			"", time.Time{}, "", nil, nil,
			"file-999",
		).
		WillReturnResult(sqlmock.NewResult(0, 0))

	body, _ := json.Marshal(map[string]interface{}{"name": "updated.txt"})

	fh := NewFileHandler(storage.NewMemoryStore())
	router := gin.New()
	router.PUT("/files/:id", fh.UpdateFile)

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
	rows := sqlmock.NewRows(fileCols).AddRow(
		"file-1", "", "t.txt", 10, "",
		time.Time{}, "", "user-1", "", time.Time{}, "", "", "",
	)
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, folder_id`)).
		WithArgs("file-1").
		WillReturnRows(rows)
	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM files WHERE id = $1`)).
		WithArgs("file-1").
		WillReturnResult(sqlmock.NewResult(0, 1))

	fh := NewFileHandler(storage.NewMemoryStore())
	router := gin.New()
	router.DELETE("/files/:id", fh.DeleteFile)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodDelete, "/files/file-1", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestDeleteFile_NotFound(t *testing.T) {
	mock := setupDB(t)
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, folder_id`)).
		WithArgs("file-999").
		WillReturnError(sql.ErrNoRows)

	fh := NewFileHandler(storage.NewMemoryStore())
	router := gin.New()
	router.DELETE("/files/:id", fh.DeleteFile)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodDelete, "/files/file-999", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}

func TestUploadFileContent_OK(t *testing.T) {
	mock := setupDB(t)
	rows := sqlmock.NewRows(fileCols).AddRow(
		"file-1", "", "doc.txt", 0, "",
		time.Time{}, "", "user-1", "", time.Time{}, "", "", "",
	)
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, folder_id`)).
		WithArgs("file-1").
		WillReturnRows(rows)
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE files SET`)).
		WithArgs(
			nil, "doc.txt", 5, "text/plain", "",
			"", time.Time{}, "", "files/file-1", nil,
			"file-1",
		).
		WillReturnResult(sqlmock.NewResult(0, 1))

	mem := storage.NewMemoryStore()
	fh := NewFileHandler(mem)
	router := gin.New()
	router.PUT("/files/:id/content", fh.UploadFileContent)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPut, "/files/file-1/content", strings.NewReader("hello"))
	req.Header.Set("Content-Type", "text/plain")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
	rc, err := mem.Get(context.Background(), "files/file-1")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = rc.Close() }()
	got, err := io.ReadAll(rc)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != "hello" {
		t.Fatalf("blob content: got %q", got)
	}
}

func TestDownloadFileContent_OK(t *testing.T) {
	mem := storage.NewMemoryStore()
	if err := mem.Put(context.Background(), "files/f1", strings.NewReader("abc"), 3, "text/plain"); err != nil {
		t.Fatal(err)
	}

	mock := setupDB(t)
	rows := sqlmock.NewRows(fileCols).AddRow(
		"f1", "", "note.txt", 3, "text/plain",
		time.Time{}, "", "user-1", "", time.Time{}, "", "files/f1", "",
	)
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, folder_id`)).
		WithArgs("f1").
		WillReturnRows(rows)

	fh := NewFileHandler(mem)
	router := gin.New()
	router.GET("/files/:id/content", fh.DownloadFileContent)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/files/f1/content", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
	if w.Body.String() != "abc" {
		t.Fatalf("body: got %q", w.Body.String())
	}
	if ct := w.Header().Get("Content-Type"); ct != "text/plain" {
		t.Fatalf("Content-Type: got %q", ct)
	}
}

func TestDownloadFileContent_NoStorageKey(t *testing.T) {
	mock := setupDB(t)
	rows := sqlmock.NewRows(fileCols).AddRow(
		"f1", "", "note.txt", 0, "",
		time.Time{}, "", "user-1", "", time.Time{}, "", "", "",
	)
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, folder_id`)).
		WithArgs("f1").
		WillReturnRows(rows)

	fh := NewFileHandler(storage.NewMemoryStore())
	router := gin.New()
	router.GET("/files/:id/content", fh.DownloadFileContent)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/files/f1/content", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}
