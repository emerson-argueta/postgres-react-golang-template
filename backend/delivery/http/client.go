package http

import (
	"emersonargueta/m/v1/communitygoaltracker"
	"net/url"
)

// Multiple service structure for client

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

// CommunitygoaltrackerJwtService returns an http implementation of
// communitygoaltracker processes using jwt authentication.
func (c *Client) CommunitygoaltrackerJwtService() communitygoaltracker.Processes {
	return &c.services.communitygoaltrackerjwt
}
