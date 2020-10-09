package http

import (
	"emersonargueta/m/v1/communitygoaltracker"
	"net/url"
)

// Client represents a client to connect to the HTTP server.
type Client struct {
	URL      *url.URL
	Services Services
}

// Services represents the services that jwt service provides
type Services struct {
	Communitygoaltracker Communitygoaltracker
}

// NewClient returns a new instance of Client.
func NewClient() *Client {
	c := &Client{}
	c.Services.Communitygoaltracker.client = c
	return c
}

// CommunitygoaltrackerService returns the service associated with the client.
func (c *Client) CommunitygoaltrackerService() communitygoaltracker.Service {
	return &c.Services.Communitygoaltracker
}
