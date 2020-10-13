package communitygoaltracker

import (
	"emersonargueta/m/v1/config"
)

// Client represents a client to connect to the communitygoaltracker services.
type Client struct {
	config               *config.Config
	Communitygoaltracker Service
}

// NewClient function
func NewClient() *Client {
	c := &Client{}
	c.Communitygoaltracker.client = c

	return c
}

// Service returns the communitygoaltracker service associated with the client.
func (c *Client) Service() Processes {
	return &c.Communitygoaltracker
}
