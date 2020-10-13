package stripe

import (
	"emersonargueta/m/v1/config"

	"github.com/stripe/stripe-go/client"
)

// Client represents a client to the underlying stripe client.
type Client struct {
	config *config.Config

	// Service
	Service Service

	//stripe client
	stripe *client.API
}

// NewClient function
func NewClient(config *config.Config) *Client {
	c := &Client{}
	c.Service.client = c

	c.config = config

	return c
}

// Initialize the stripe client.
func (c *Client) Initialize() {
	apiKey := c.config.PaymentGateway.APIKey
	stripe := client.New(apiKey, nil)

	c.stripe = stripe
}

// StripeService returns the stripe service associated with the client.
func (c *Client) StripeService() Service {
	return c.Service
}
