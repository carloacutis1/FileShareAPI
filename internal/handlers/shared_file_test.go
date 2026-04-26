package handlers

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
)

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

// ── CreateSharedFile ──────────────────────────────────────────────────────────
func TestCreateSharedFile_BadJSON(t *testing.T) {
	router := gin.New()
	router.POST("/shared-files", CreateSharedFile)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/shared-files", bytes.NewBufferString("{bad}"))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}
func TestCreateSharedFile_OK(t *testing.T) {
	mock := setupDB(t)

	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO shared_files`)).
		WithArgs("sf_1", "file_1", "user_1", "user_2", "read").
		WillReturnResult(sqlmock.NewResult(1, 1))

	body, _ := json.Marshal(map[string]interface{}{
		"id":          "sf_1",
		"file_id":     "file_1",
		"owner_id":    "user_1",
		"shared_with": "user_2",
		"permissions": "read",
	})

	router := gin.New()
	router.POST("/shared-files", CreateSharedFile)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/shared-files", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", w.Code, w.Body.String())
	}
}
func TestCreateSharedFile_DBError(t *testing.T) {
	mock := setupDB(t)

	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO shared_files`)).
		WithArgs("sf_1", "file_1", "user_1", "user_2", "read").
		WillReturnError(sql.ErrConnDone)

	body, _ := json.Marshal(map[string]interface{}{
		"id":          "sf_1",
		"file_id":     "file_1",
		"owner_id":    "user_1",
		"shared_with": "user_2",
		"permissions": "read",
	})

	router := gin.New()
	router.POST("/shared-files", CreateSharedFile)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/shared-files", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", w.Code)
	}
}

// ── DeleteSharedFile ──────────────────────────────────────────────────────────

func authAs(userID string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("user_id", userID)
		c.Next()
	}
}

func TestDeleteSharedFile_OK(t *testing.T) {
	mock := setupDB(t)
	sfCols := []string{"id", "file_id", "owner_id", "shared_with", "permissions"}
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, file_id, owner_id, shared_with, permissions FROM shared_files WHERE id = $1`)).
		WithArgs("sf-1").
		WillReturnRows(sqlmock.NewRows(sfCols).AddRow("sf-1", "file_1", "user_1", "user_2", "read"))
	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM shared_files WHERE id = $1`)).
		WithArgs("sf-1").
		WillReturnResult(sqlmock.NewResult(0, 1))

	router := gin.New()
	router.DELETE("/shared-files/:id", authAs("user_1"), DeleteSharedFile)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodDelete, "/shared-files/sf-1", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestDeleteSharedFile_NotFound(t *testing.T) {
	mock := setupDB(t)
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, file_id, owner_id, shared_with, permissions FROM shared_files WHERE id = $1`)).
		WithArgs("sf-999").
		WillReturnError(sql.ErrNoRows)

	router := gin.New()
	router.DELETE("/shared-files/:id", authAs("user_1"), DeleteSharedFile)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodDelete, "/shared-files/sf-999", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}

func TestDeleteSharedFile_DBError(t *testing.T) {
	mock := setupDB(t)
	sfCols := []string{"id", "file_id", "owner_id", "shared_with", "permissions"}
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, file_id, owner_id, shared_with, permissions FROM shared_files WHERE id = $1`)).
		WithArgs("sf-1").
		WillReturnRows(sqlmock.NewRows(sfCols).AddRow("sf-1", "file_1", "user_1", "user_2", "read"))
	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM shared_files WHERE id = $1`)).
		WithArgs("sf-1").
		WillReturnError(sql.ErrConnDone)

	router := gin.New()
	router.DELETE("/shared-files/:id", authAs("user_1"), DeleteSharedFile)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodDelete, "/shared-files/sf-1", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", w.Code)
	}
}

// ── UpdateSharedFile ──────────────────────────────────────────────────────────

func TestUpdateSharedFile_BadJSON(t *testing.T) {
	router := gin.New()
	router.PUT("/shared-files/:id", UpdateSharedFile)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPut, "/shared-files/sf-1", bytes.NewBufferString("{bad}"))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestUpdateSharedFile_OK(t *testing.T) {
	mock := setupDB(t)
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE shared_files SET`)).
		WithArgs("file_1", "user_1", "user_2", "write", "sf-1").
		WillReturnResult(sqlmock.NewResult(0, 1))

	body, _ := json.Marshal(map[string]interface{}{
		"file_id":     "file_1",
		"owner_id":    "user_1",
		"shared_with": "user_2",
		"permissions": "write",
	})

	router := gin.New()
	router.PUT("/shared-files/:id", UpdateSharedFile)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPut, "/shared-files/sf-1", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestUpdateSharedFile_NotFound(t *testing.T) {
	mock := setupDB(t)
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE shared_files SET`)).
		WithArgs("file_1", "user_1", "user_2", "write", "sf-999").
		WillReturnResult(sqlmock.NewResult(0, 0))

	body, _ := json.Marshal(map[string]interface{}{
		"file_id":     "file_1",
		"owner_id":    "user_1",
		"shared_with": "user_2",
		"permissions": "write",
	})

	router := gin.New()
	router.PUT("/shared-files/:id", UpdateSharedFile)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPut, "/shared-files/sf-999", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}

func TestUpdateSharedFile_DBError(t *testing.T) {
	mock := setupDB(t)
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE shared_files SET`)).
		WithArgs("file_1", "user_1", "user_2", "write", "sf-1").
		WillReturnError(sql.ErrConnDone)

	body, _ := json.Marshal(map[string]interface{}{
		"file_id":     "file_1",
		"owner_id":    "user_1",
		"shared_with": "user_2",
		"permissions": "write",
	})

	router := gin.New()
	router.PUT("/shared-files/:id", UpdateSharedFile)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPut, "/shared-files/sf-1", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", w.Code)
	}
}
