package http

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

const (
	v1BasePath  = "/api/v1"
	healthPath  = "/health"
	addUserPath = "/users"
)

// LoadRoutes loads the REST API routes
func (h *H) LoadRoutes() http.Handler {
	// Init router
	e := echo.New()

	// Use custom error handler

	// Load Middlewares
	e.Use(middleware.Recover())

	// Base API Group
	v1Base := e.Group(v1BasePath)

	// v1 Health Check Endpoint
	v1Base.GET(healthPath, h.getHealth)
	return e
}
