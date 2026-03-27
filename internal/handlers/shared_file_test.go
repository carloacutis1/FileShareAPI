package handlers

import (
	"database/sql"
	"bytes"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"
	"encoding/json"

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
