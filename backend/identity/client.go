package identity

import (
	"emersonargueta/m/v1/config"
)

// Client represents a client to the underlying stripe client.
type Client struct {
	config  *config.Config
	Service Service
}

// NewClient function
func NewClient(config *config.Config) *Client {
	c := &Client{}
	c.Service.client = c

	c.config = config

	return c
}

// IdentityService returns the jwt service associated with the client.
func (c *Client) IdentityService() Service { return c.Service }
