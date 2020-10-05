package jwt

import (
	"emersonargueta/m/v1/config"
)

// Client represents a client to the underlying stripe client.
type Client struct {
	config *config.Config

	// Services
	Services Services

	//authorization client not used this is just an example
	/** auth0 *client.API **/
}

// Services represents the services that jwt service provides
type Services struct {
	Administrator Administrator
}

// NewClient function
func NewClient() *Client {
	c := &Client{}
	c.Services.Administrator.client = c

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
	//authorization client not used this is just an example

	/**
	apiKey := c.config.Authorization.APIKey
	auth0 := client.New(apiKey, nil)

	c.auth0 = stripe
	**/
}

// AdministratorService returns the admin service associated with the client.
func (c *Client) AdministratorService() Administrator { return c.Services.Administrator }
