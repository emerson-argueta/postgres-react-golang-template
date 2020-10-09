package identity

import (
	"emersonargueta/m/v1/config"
)

// Client represents a client to the underlying stripe client.
type Client struct {
	config   *config.Config
	Services Services
}

// Services represents the services that jwt service provides
type Services struct {
	Identity Identity
}

// NewClient function
func NewClient() *Client {
	c := &Client{}
	c.Services.Identity.client = c

	// get configuration stucts via .env file
	config, err := config.NewConfig()
	if err != nil {
		panic(err)
	}
	c.config = config

	return c
}

// Initialize the stripe client.
func (c *Client) Initialize() {
}

// IdentityService returns the jwt service associated with the client.
func (c *Client) IdentityService() Identity { return c.Services.Identity }