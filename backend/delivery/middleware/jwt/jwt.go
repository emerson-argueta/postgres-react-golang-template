package jwt

import (
	"emersonargueta/m/v1/delivery/middleware"

	"github.com/labstack/echo"
	echo_middleware "github.com/labstack/echo/middleware"
)

var _ middleware.Processes = &service{}

type service struct {
	client *Client
}

// JwtMiddleware function
func (m *service) MiddlewareFunc() echo.MiddlewareFunc {
	return echo_middleware.JWTWithConfig(echo_middleware.JWTConfig{
		SigningKey: []byte(m.client.config.Authorization.Secret),
	})
}
