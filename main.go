// Package main is the entry point for the Kite session management API
package main

import (
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/nsvirk/mbservices/handlers/session"
)

func main() {
	e := echo.New()
	e.HideBanner = true

	// Set up middleware
	setupMiddleware(e)

	// Set up routes
	setupRoutes(e)

	// Start the server
	startServer(e)
}

// setupMiddleware configures and adds middleware to the Echo instance
func setupMiddleware(e *echo.Echo) {
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "${time_rfc3339}: ip=${remote_ip}, req=${method}, uri=${uri}, status=${status}\n",
	}))
	e.Use(middleware.Recover())
}

// setupRoutes adds the API routes to the Echo instance
func setupRoutes(e *echo.Echo) {
	e.POST("/session/totp", session.GenerateTOTP)
	e.POST("/session/login", session.GenerateSession)
	e.POST("/session/valid", session.CheckEnctoken)
}

// startServer starts the Echo server on the specified port
func startServer(e *echo.Echo) {
	port := os.Getenv("MBSERVICES_PORT")
	if port == "" {
		port = "3008"
	}
	e.Logger.Infof("Starting mbservices server on port %s", port)
	e.Logger.Fatal(e.Start(":" + port))
}
