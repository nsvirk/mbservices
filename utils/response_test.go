package utils

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestSuccessResponse(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Test data
	testData := map[string]string{"key": "value"}

	// Perform the test
	err := SuccessResponse(c, testData)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	// Parse the response body
	var response Response
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err)

	// Check the response structure
	assert.Equal(t, "ok", response.Status)

	// Check the data field
	responseData, ok := response.Data.(map[string]interface{})
	assert.True(t, ok, "Data should be a map[string]interface{}")
	assert.Equal(t, "value", responseData["key"])

	assert.Empty(t, response.ErrorType)
	assert.Empty(t, response.Message)
}

func TestErrorResponse(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Test data
	testStatus := http.StatusBadRequest
	testErrorType := "validation_error"
	testMessage := "Invalid input"

	// Perform the test
	err := ErrorResponse(c, testStatus, testErrorType, testMessage)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, testStatus, rec.Code)

	// Parse the response body
	var response Response
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err)

	// Check the response structure
	assert.Equal(t, "error", response.Status)
	assert.Nil(t, response.Data)
	assert.Equal(t, testErrorType, response.ErrorType)
	assert.Equal(t, testMessage, response.Message)
}

func TestResponseStructure(t *testing.T) {
	// Test Success Response Structure
	successResp := Response{
		Status: "ok",
		Data:   "test data",
	}

	successJSON, err := json.Marshal(successResp)
	assert.NoError(t, err)
	assert.Contains(t, string(successJSON), "status")
	assert.Contains(t, string(successJSON), "data")
	assert.NotContains(t, string(successJSON), "error_type")
	assert.NotContains(t, string(successJSON), "message")

	// Test Error Response Structure
	errorResp := Response{
		Status:    "error",
		ErrorType: "test_error",
		Message:   "test message",
	}

	errorJSON, err := json.Marshal(errorResp)
	assert.NoError(t, err)
	assert.Contains(t, string(errorJSON), "status")
	assert.Contains(t, string(errorJSON), "error_type")
	assert.Contains(t, string(errorJSON), "message")
	assert.NotContains(t, string(errorJSON), "data")
}
