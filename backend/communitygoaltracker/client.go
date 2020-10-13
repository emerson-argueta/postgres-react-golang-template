package communitygoaltracker

import (
	"emersonargueta/m/v1/config"
)

// Client represents a client to connect to the communitygoaltracker services.
type Client struct {
	config  *config.Config
	service Service
}

// NewClient function
func NewClient(config *config.Config) *Client {
	c := &Client{}
	c.service.client = c

	c.config = config

	return c
}

// Service returns the communitygoaltracker service associated with the client.
func (c *Client) Service() Processes {
	return &c.service
}
