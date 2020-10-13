package jwt

import (
	"emersonargueta/m/v1/authorization"
	"emersonargueta/m/v1/config"
)

// Client to the jwt authorization service.
type Client struct {
	config  *config.Config
	service service
}

// NewClient to the jwt authorization service.
func NewClient(config *config.Config) *Client {
	c := &Client{}

	c.service.client = c

	c.config = config

	return c
}

// Service returns the jwt authorization service associated with the client.
func (c *Client) Service() authorization.Processes { return &c.service }
