package identity

import (
	"emersonargueta/m/v1/config"
)

// Client represents a client to the underlying stripe client.
type Client struct {
	config   *config.Config
	Identity Service
}

// NewClient function
func NewClient() *Client {
	c := &Client{}
	c.Identity.client = c

	return c
}

// Service returns the jwt service associated with the client.
func (c *Client) Service() Processes { return &c.Identity }
