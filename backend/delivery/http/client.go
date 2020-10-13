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
	communitygoaltrackerjwt communitygoaltrackerjwtservice
}

// NewClient creates a connection to http services.
func NewClient() *Client {
	c := &Client{}
	c.services.communitygoaltrackerjwt.client = c
	return c
}

// CommunitygoaltrackerJwtService returns the service associated with the client.
func (c *Client) CommunitygoaltrackerJwtService() communitygoaltracker.Processes {
	return &c.services.communitygoaltrackerjwt
}
