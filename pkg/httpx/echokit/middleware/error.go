package echomiddlewarekit

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/wasay-usmani/go-boilerplate/pkg/httpx"
	"github.com/wasay-usmani/go-boilerplate/pkg/logkit"
)

// HTTPErrorHandler handles all HTTP errors returned by Echo routes and middleware.
func HTTPErrorHandler(logger logkit.Logger) echo.HTTPErrorHandler {
	return func(err error, c echo.Context) {
		// Retrieve the request ID from the response headers
		requestID := c.Response().Header().Get(httpx.RequestIdHeader)

		// Default to internal server error
		code := http.StatusInternalServerError
		message := http.StatusText(code)

		// Check if the error is an *echo.HTTPError
		if he, ok := err.(*echo.HTTPError); ok {
			code = he.Code
			if he.Message != nil {
				message = fmt.Sprintf("%v", he.Message)
			}
			// Log the error with context
			logger.Warn("HTTP error",
				"request_id", requestID,
				"code", code,
				"message", message,
				"internal", he.Internal)
		} else {
			// Log unexpected errors
			logger.Error("Unhandled error", err, "request_id", requestID)
		}

		// Construct the error response
		// TODO: Implement error response construction after finalizing the error response structure.
		resp := map[string]any{
			"error": map[string]any{
				"code":    code,
				"message": message,
			},
		}

		// Send the JSON response
		if err := c.JSON(code, resp); err != nil {
			logger.Error("Failed to send error respons", err, "request_id", requestID)
		}
	}
}
