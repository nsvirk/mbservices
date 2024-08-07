// Package session implements handlers for Kite session management
package session

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	kitesession "github.com/nsvirk/gokitesession"
	"github.com/nsvirk/mbservices/utils"
)

// TOTPRequest represents the request body for generating a TOTP value
type TOTPRequest struct {
	TOTPSecret string `json:"totp_secret"`
}

// LoginRequest represents the request body for user login
type LoginRequest struct {
	UserID    string `json:"user_id"`
	Password  string `json:"password"`
	TOTPValue string `json:"totp_value"`
}

// EnctokenRequest represents the request body for checking an enctoken
type EnctokenRequest struct {
	Enctoken string `json:"enctoken"`
}

// GenerateTOTP handles the generation of TOTP values
func GenerateTOTP(c echo.Context) error {
	req := new(TOTPRequest)

	if err := c.Bind(req); err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, "InputException", "Invalid request body")
	}

	if req.TOTPSecret == "" {
		return utils.ErrorResponse(c, http.StatusBadRequest, "InputException", "totp_secret is required")
	}

	totpValue, err := kitesession.GenerateTOTPValue(req.TOTPSecret)
	if err != nil {
		return utils.ErrorResponse(c, http.StatusInternalServerError, "ServerException", "Failed to generate TOTP value")
	}

	return utils.SuccessResponse(c, map[string]string{"totp_value": totpValue})
}

// GenerateSession handles the user login process and generates a session
func GenerateSession(c echo.Context) error {
	req := new(LoginRequest)
	if err := c.Bind(req); err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, "InputException", "Invalid request body")
	}

	if req.UserID == "" || req.Password == "" || req.TOTPValue == "" {
		return utils.ErrorResponse(c, http.StatusBadRequest, "InputException", "user_id, password, and totp_value are required")
	}

	ks := kitesession.New()
	session, err := ks.GenerateSession(req.UserID, req.Password, req.TOTPValue)
	if err != nil {
		return utils.ErrorResponse(c, http.StatusUnauthorized, "AuthenticationException", fmt.Sprintf("Login failed: %v", err))
	}

	return utils.SuccessResponse(c, session)
}

// CheckEnctoken validates the provided enctoken
func CheckEnctoken(c echo.Context) error {
	req := new(EnctokenRequest)
	if err := c.Bind(req); err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, "InputException", "Invalid request body")
	}

	if req.Enctoken == "" {
		return utils.ErrorResponse(c, http.StatusBadRequest, "InputException", "enctoken is required")
	}

	ks := kitesession.New()
	isValid, err := ks.CheckEnctokenValid(req.Enctoken)
	if err != nil {
		return utils.ErrorResponse(c, http.StatusInternalServerError, "ServerException", fmt.Sprintf("Failed to check enctoken: %v", err))
	}

	return utils.SuccessResponse(c, map[string]bool{"is_valid": isValid})
}
