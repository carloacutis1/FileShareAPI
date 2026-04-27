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

// init and setupDB are declared in handlers_test.go and shared across this package.

var userCols = []string{"id", "username", "email", "dob", "created_at"}

// ── ListUserModels ────────────────────────────────────────────────────────────

func TestListUserModels_OK(t *testing.T) {
	mock := setupDB(t)
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, username, email`)).
		WillReturnRows(sqlmock.NewRows(userCols))

	router := gin.New()
	router.GET("/users", ListUserModels)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/users", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestListUserModels_DBError(t *testing.T) {
	mock := setupDB(t)
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, username, email`)).
		WillReturnError(sql.ErrConnDone)

	router := gin.New()
	router.GET("/users", ListUserModels)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/users", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", w.Code)
	}
}

// ── CreateUserModel ───────────────────────────────────────────────────────────

func TestCreateUserModel_BadJSON(t *testing.T) {
	router := gin.New()
	router.POST("/users", CreateUserModel)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/users", bytes.NewBufferString("{bad}"))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestCreateUserModel_OK(t *testing.T) {
	mock := setupDB(t)
	hash := "$2b$10$EnxcfsfGyZZPwkeWIkMLze2qRS1vhBsUeWaxCX6AYDhHbdQ9RKxY6"
	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO users (id, username, email, dob, password_hash)`)).
		WithArgs("user_6", "Test User", "test@example.com", "01/01/2000", hash).
		WillReturnResult(sqlmock.NewResult(1, 1))

	body, _ := json.Marshal(map[string]interface{}{
		"id":            "user_6",
		"name":          "Test User",
		"email":         "test@example.com",
		"dob":           "01/01/2000",
		"password_hash": hash,
	})

	router := gin.New()
	router.POST("/users", CreateUserModel)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/users", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", w.Code, w.Body.String())
	}
}

func TestCreateUserModel_DBError(t *testing.T) {
	mock := setupDB(t)
	hash := "$2b$10$EnxcfsfGyZZPwkeWIkMLze2qRS1vhBsUeWaxCX6AYDhHbdQ9RKxY6"
	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO users (id, username, email, dob, password_hash)`)).
		WithArgs("user_6", "Test User", "test@example.com", "01/01/2000", hash).
		WillReturnError(sql.ErrConnDone)

	body, _ := json.Marshal(map[string]interface{}{
		"id":            "user_6",
		"name":          "Test User",
		"email":         "test@example.com",
		"dob":           "01/01/2000",
		"password_hash": hash,
	})

	router := gin.New()
	router.POST("/users", CreateUserModel)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/users", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", w.Code)
	}
}

// ── UpdateUserModel ───────────────────────────────────────────────────────────

func TestUpdateUserModel_BadJSON(t *testing.T) {
	router := gin.New()
	router.PUT("/users/:id", UpdateUserModel)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPut, "/users/user_1", bytes.NewBufferString("{bad}"))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestUpdateUserModel_OK(t *testing.T) {
	mock := setupDB(t)
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE users SET`)).
		WithArgs("John Updated", "updated@example.com", "01/01/1900", "user_1").
		WillReturnResult(sqlmock.NewResult(0, 1))

	body, _ := json.Marshal(map[string]interface{}{
		"name":  "John Updated",
		"email": "updated@example.com",
		"dob":   "01/01/1900",
	})

	router := gin.New()
	router.PUT("/users/:id", UpdateUserModel)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPut, "/users/user_1", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestUpdateUserModel_NotFound(t *testing.T) {
	mock := setupDB(t)
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE users SET`)).
		WithArgs("John Updated", "updated@example.com", "01/01/1900", "user_999").
		WillReturnResult(sqlmock.NewResult(0, 0))

	body, _ := json.Marshal(map[string]interface{}{
		"name":  "John Updated",
		"email": "updated@example.com",
		"dob":   "01/01/1900",
	})

	router := gin.New()
	router.PUT("/users/:id", UpdateUserModel)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPut, "/users/user_999", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}

// ── DeleteUserModel ───────────────────────────────────────────────────────────

func TestDeleteUserModel_OK(t *testing.T) {
	mock := setupDB(t)
	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM users WHERE id = $1`)).
		WithArgs("user_1").
		WillReturnResult(sqlmock.NewResult(0, 1))

	router := gin.New()
	router.DELETE("/users/:id", DeleteUserModel)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodDelete, "/users/user_1", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", w.Code)
	}
}

func TestDeleteUserModel_NotFound(t *testing.T) {
	mock := setupDB(t)
	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM users WHERE id = $1`)).
		WithArgs("user_999").
		WillReturnResult(sqlmock.NewResult(0, 0))

	router := gin.New()
	router.DELETE("/users/:id", DeleteUserModel)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodDelete, "/users/user_999", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}
