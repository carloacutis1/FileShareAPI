package routes

import (
	"net/http"
	"net/http/httptest"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"gitlab.franciscan.edu/spr-26-backend/basement/internal/database"
)

func init() {
	gin.SetMode(gin.TestMode)
}

// setupDummyDB injects a sqlmock DB so handlers that call database.DB don't
// panic on a nil pointer. Unexpected DB calls return sqlmock errors (causing
// 500 responses), which is fine for route-registration checks.
func setupDummyDB(t *testing.T) {
	t.Helper()
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New: %v", err)
	}
	t.Cleanup(func() { _ = db.Close() })
	database.DB = db
}

// assertRegistered fails the test if the given method+path returns 404.
func assertRegistered(t *testing.T, router *gin.Engine, method, path string) {
	t.Helper()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(method, path, nil)
	router.ServeHTTP(w, req)
	if w.Code == http.StatusNotFound {
		t.Errorf("%s %s: route not registered (got 404)", method, path)
	}
}

// TestSetupUserRoutes verifies that GET /users is reachable and returns 200
// (uses in-memory data, no DB required).
func TestSetupUserRoutes(t *testing.T) {
	router := gin.New()
	SetupUserRoutes(router)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/users", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("GET /users: expected 200, got %d", w.Code)
	}
}

// TestSetupFileRoutes verifies all four file endpoints are registered.
// GET  /files        – hits DB; sqlmock returns an error → handler replies 500 (not 404).
// POST /files        – fails JSON parse on empty body → 400 (not 404).
// PUT  /files/:id    – fails JSON parse on empty body → 400 (not 404).
// DELETE /files/:id  – hits DB; sqlmock returns an error → 500 (not 404).
func TestSetupFileRoutes(t *testing.T) {
	setupDummyDB(t)

	router := gin.New()
	SetupFileRoutes(router)

	for _, tc := range []struct{ method, path string }{
		{http.MethodGet, "/files"},
		{http.MethodPost, "/files"},
		{http.MethodPut, "/files/some-id"},
		{http.MethodDelete, "/files/some-id"},
	} {
		assertRegistered(t, router, tc.method, tc.path)
	}
}

// TestSetupSharedFileRoutes verifies GET /shared-files is registered.
func TestSetupSharedFileRoutes(t *testing.T) {
	setupDummyDB(t)

	router := gin.New()
	SetupSharedFileRoutes(router)

	assertRegistered(t, router, http.MethodGet, "/shared-files")
}

// TestSetupRoutes verifies that the combined router wires all domains.
func TestSetupRoutes(t *testing.T) {
	setupDummyDB(t)

	router := gin.New()
	SetupRoutes(router)

	for _, tc := range []struct{ method, path string }{
		{http.MethodGet, "/users"},
		{http.MethodGet, "/files"},
		{http.MethodPost, "/files"},
		{http.MethodPut, "/files/some-id"},
		{http.MethodDelete, "/files/some-id"},
		{http.MethodGet, "/shared-files"},
	} {
		assertRegistered(t, router, tc.method, tc.path)
	}
}
