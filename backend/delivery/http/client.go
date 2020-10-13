package http

import (
	"emersonargueta/m/v1/communitygoaltracker"
	"net/url"
)

// Client represents a client to connect to the HTTP server.
type Client struct {
	URL      *url.URL
	services services
}

// Services represents the services that jwt service provides
type services struct {
	Communitygoaltracker CommunitygoaltrackerService
}

// NewClient returns a new instance of Client.
func NewClient() *Client {
	c := &Client{}
	c.services.Communitygoaltracker.client = c
	return c
}

// CommunitygoaltrackerService returns the service associated with the client.
func (c *Client) CommunitygoaltrackerService() communitygoaltracker.Processes {
	return &c.services.Communitygoaltracker
}
