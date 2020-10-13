package communitygoaltracker

import (
	"emersonargueta/m/v1/config"
)

// Client to the communitygoaltracker service.
type Client struct {
	config  *config.Config
	service service
}

// NewClient creates a connection to a communitygoaltracker service.
func NewClient() *Client {
	c := &Client{}
	c.service.client = c

	return c
}

// Service returns the communitygoaltracker service associated with the client.
func (c *Client) Service() Processes { return &c.service }
