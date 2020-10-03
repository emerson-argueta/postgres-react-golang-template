package http

import (
	"net/url"
)

// Client represents a client to connect to the HTTP server.
type Client struct {
	URL      *url.URL
	Services Services
}

// Services represents the services that jwt service provides
type Services struct {
	Administrator Administrator
}

// NewClient returns a new instance of Client.
func NewClient() *Client {
	c := &Client{}
	c.Services.Administrator.client = c
	return c
}

// AdministratorService returns the admin service associated with the client.
func (c *Client) AdministratorService() AdministratorActions { return &c.Services.Administrator }
