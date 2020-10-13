package http

import (
	"emersonargueta/m/v1/communitygoaltracker"
	"net/url"
)

// Client to http services.
type Client struct {
	URL      *url.URL
	services services
}

type services struct {
	communitygoaltracker communitygoaltrackerservice
}

// NewClient creates a connection to http services.
func NewClient() *Client {
	c := &Client{}
	c.services.communitygoaltracker.client = c
	return c
}

// CommunitygoaltrackerService returns the service associated with the client.
func (c *Client) CommunitygoaltrackerService() communitygoaltracker.Processes {
	return &c.services.communitygoaltracker
}
