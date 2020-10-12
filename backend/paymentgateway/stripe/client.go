package stripe

import (
	"emersonargueta/m/v1/config"

	"github.com/stripe/stripe-go/client"
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
	Stripe Stripe
}

// NewClient function
func NewClient(config *config.Config) *Client {
	c := &Client{}
	c.Services.Stripe.client = c

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
func (c *Client) StripeService() Stripe {
	return c.Services.Stripe
}
