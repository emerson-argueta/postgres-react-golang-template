package stripe

import (
	"github.com/stripe/stripe-go/client"
	"trustdonations.org/m/v2/config"
)

// Client represents a client to the underlying stripe client.
type Client struct {
	config *config.Config

	// Services
	Services Services

	//stripe client
	stripe *client.API
}

// Services represents the services that http service provides
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
	apiKey := c.config.PaymentGateway.APIKey
	stripe := client.New(apiKey, nil)

	c.stripe = stripe
}

// AdministratorService returns the admin service associated with the client.
// func (c *Client) AdministratorService() administrator.SubscriptionActions {
// 	return &c.Services.Administrator
// }
func (c *Client) AdministratorService() Administrator {
	return c.Services.Administrator
}
