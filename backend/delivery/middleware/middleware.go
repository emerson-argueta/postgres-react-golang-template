package middleware

import "github.com/labstack/echo"

// Processes middleware
type Processes interface {
	// MiddlewareFunc used in echo router
	MiddlewareFunc() echo.MiddlewareFunc
}
