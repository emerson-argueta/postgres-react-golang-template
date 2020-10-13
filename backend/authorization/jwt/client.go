package jwt

import (
	"emersonargueta/m/v1/authorization"
	"emersonargueta/m/v1/config"
	"emersonargueta/m/v1/delivery/middleware"
)

// Client represents a client to the jwt authorization serrvice.
type Client struct {
	config  *config.Config
	service service
}

// NewClient function
func NewClient(config *config.Config) *Client {
	c := &Client{}

	c.service.client = c
	c.service.middleware = middleware.New(config)

	c.config = config

	return c
}

// JwtService returns the jwt service associated with the client.
func (c *Client) Service() authorization.Processes { return &c.service }
