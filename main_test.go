package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestSetupRoutes(t *testing.T) {
	e := echo.New()
	setupRoutes(e)

	// Test TOTP route
	req := httptest.NewRequest(http.MethodPost, "/session/totp", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusBadRequest, rec.Code) // Assuming BadRequest for empty body

	// Test Login route
	req = httptest.NewRequest(http.MethodPost, "/session/login", nil)
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusBadRequest, rec.Code) // Assuming BadRequest for empty body

	// Test Enctoken route
	req = httptest.NewRequest(http.MethodPost, "/session/valid", nil)
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusBadRequest, rec.Code) // Assuming BadRequest for empty body
}

func TestStartServer(t *testing.T) {
	// Save current PORT env var
	oldPort := os.Getenv("PORT")
	defer os.Setenv("PORT", oldPort)

	// Test with PORT env var set
	os.Setenv("PORT", "9090")
	e := echo.New()
	go startServer(e) // Start server in a goroutine
	// Add a small delay to allow server to start
	// time.Sleep(100 * time.Millisecond)
	resp, err := http.Get("http://localhost:9090")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode) // Assuming no route set returns 404

	// Test with PORT env var unset
	os.Unsetenv("PORT")
	e = echo.New()
	go startServer(e) // Start server in a goroutine
	// Add a small delay to allow server to start
	// time.Sleep(100 * time.Millisecond)
	resp, err = http.Get("http://localhost:8080")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode) // Assuming no route set returns 404
}

func TestMain(t *testing.T) {
	// This is a basic test to ensure main() runs without panicking
	// For a more comprehensive test, you might want to refactor main() to return an echo.Echo instance
	assert.NotPanics(t, func() {
		go main()
	})
}
