package handlers

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
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
