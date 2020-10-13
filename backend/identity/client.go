package identity

import (
	"emersonargueta/m/v1/config"
)

// Single service structure for package

// Client to the identity service.
type Client struct {
	config  *config.Config
	service service
}

// NewClient creates a connection to an identity service.
func NewClient() *Client {
	c := &Client{}
	c.service.client = c

	return c
}

// Service returns the identity service associated with the client.
func (c *Client) Service() Processes { return &c.service }
