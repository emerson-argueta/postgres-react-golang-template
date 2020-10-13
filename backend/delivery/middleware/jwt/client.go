package jwt

import (
	"emersonargueta/m/v1/config"
	"emersonargueta/m/v1/delivery/middleware"
)

// Client to the postgres services.
type Client struct {
	config  *config.Config
	service service
}

// NewClient creates a connection to the postgres services.
func NewClient(config *config.Config) *Client {
	c := &Client{}

	c.service.client = c
	c.config = config

	return c
}

// JwtService returns the user service associated with the client.
func (c *Client) Service() middleware.Processes { return &c.service }
