package middleware

import (
	"emersonargueta/m/v1/config"
)

// Client to the postgres services.
type Client struct {
	config   *config.Config
	services services
}

type services struct {
	jwt jwtservice
}

// NewClient creates a connection to the postgres services.
func NewClient(config *config.Config) *Client {
	c := &Client{}

	c.services.jwt.client = c
	c.config = config

	return c
}

// JwtService returns the user service associated with the client.
func (c *Client) JwtService() Processes { return &c.services.jwt }
