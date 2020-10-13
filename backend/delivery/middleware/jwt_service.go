package middleware

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

var _ Processes = &jwtservice{}

type jwtservice struct {
	client *Client
}

// JwtMiddleware function
func (m *jwtservice) MiddlewareFunc() echo.MiddlewareFunc {
	return middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey: []byte(m.client.config.Authorization.Secret),
	})
}
