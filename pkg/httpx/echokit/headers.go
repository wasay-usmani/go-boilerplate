package echokit

import "github.com/labstack/echo/v4"

// AddHeaders returns middleware that sets the provided headers on the HTTP response.
// It applies the headers before invoking the next handler in the chain.
func AddHeaders(headers map[string]string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			for k, v := range headers {
				c.Response().Header().Set(k, v)
			}

			if err := next(c); err != nil {
				c.Error(err)
			}

			return nil
		}
	}
}
