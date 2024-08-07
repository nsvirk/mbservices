package session

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestGenerateTOTP(t *testing.T) {
	e := echo.New()

	// Test with valid input
	jsonStr := `{"totp_secret":"JBSWY3DPEHPK3PXP"}`
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(jsonStr))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if assert.NoError(t, GenerateTOTP(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Contains(t, rec.Body.String(), "totp_value")
	}

	// Test with invalid input
	jsonStr = `{"invalid_field":"value"}`
	req = httptest.NewRequest(http.MethodPost, "/", strings.NewReader(jsonStr))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)

	if assert.NoError(t, GenerateTOTP(c)) {
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Contains(t, rec.Body.String(), "totp_secret is required")
	}
}

func TestGenerateSession(t *testing.T) {
	e := echo.New()

	// Test with valid input
	jsonStr := `{"user_id":"testuser","password":"testpass","totp_value":"123456"}`
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(jsonStr))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := GenerateSession(c)
	assert.NoError(t, err)
	// Note: The actual response will depend on the kitesession.GenerateSession implementation
	// You might need to mock this function for a more accurate test

	// Test with invalid input
	jsonStr = `{"user_id":"testuser"}`
	req = httptest.NewRequest(http.MethodPost, "/", strings.NewReader(jsonStr))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)

	if assert.NoError(t, GenerateSession(c)) {
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Contains(t, rec.Body.String(), "user_id, password, and totp_value are required")
	}
}

func TestCheckEnctoken(t *testing.T) {
	e := echo.New()

	// Test with valid input
	jsonStr := `{"enctoken":"testenctoken"}`
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(jsonStr))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := CheckEnctoken(c)
	assert.NoError(t, err)
	// Note: The actual response will depend on the kitesession.CheckEnctokenValid implementation
	// You might need to mock this function for a more accurate test

	// Test with invalid input
	jsonStr = `{}`
	req = httptest.NewRequest(http.MethodPost, "/", strings.NewReader(jsonStr))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)

	if assert.NoError(t, CheckEnctoken(c)) {
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Contains(t, rec.Body.String(), "enctoken is required")
	}
}
