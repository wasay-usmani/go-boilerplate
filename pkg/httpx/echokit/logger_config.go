package echokit

import (
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/wasay-usmani/go-boilerplate/pkg/httpx"
)

// LoggerConfig defines custom logging behavior for Echo middleware.
// It skips logging for health check endpoints and formats logs in JSON with precise timestamps.
var LoggerConfig = middleware.LoggerConfig{
	Skipper: func(e echo.Context) bool {
		return strings.Contains(e.Request().RequestURI, httpx.HealthCheckPath)
	},
	Format: `{"time":"${time_rfc3339_nano}","request_id":"${request_id}","remote_ip":"${remote_ip}",` +
		`"host":"${host}","method":"${method}","uri":"${uri}","user_agent":"${user_agent}",` +
		`"status":${status},"error":"${error}", "latency":"${latency_human}"` + `}` + "\n",
	CustomTimeFormat: "2006-01-02 15:04:05.00000",
}
