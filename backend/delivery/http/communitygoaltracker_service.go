package http

import "emersonargueta/m/v1/communitygoaltracker"

var _ communitygoaltracker.Services = &Communitygoaltracker{}

// Communitygoaltracker represents an HTTP implementation of communitygoaltracker.Service.
type Communitygoaltracker struct {
	client *Client
}
