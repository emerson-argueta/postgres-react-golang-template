package jwt

import (
	"emersonargueta/m/v1/config"
	"emersonargueta/m/v1/delivery/middleware"
)

// Client represents a client to the underlying stripe client.
type Client struct {
	config *config.Config

	// Services
	Services Services
}

// Services represents the services that jwt service provides
type Services struct {
	Jwt Jwt
}

// NewClient function
func NewClient(config *config.Config) *Client {
	c := &Client{}

	c.Services.Jwt.client = c
	c.Services.Jwt.middleware = middleware.New(config)

	c.config = config

	return c
}

// JwtService returns the jwt service associated with the client.
func (c *Client) JwtService() Jwt { return c.Services.Jwt }
