package paymentgateway

import (
	"emersonargueta/m/v1/config"
)

// Multiple service structure for client

// Client to the stripe service.
type Client struct {
	config   *config.Config
	services services
}

type services struct {
	stripe stripeservice
}

// NewClient creates connection to strinpe service.
func NewClient(config *config.Config) *Client {
	c := &Client{}

	c.services.stripe.client = c
	c.config = config

	return c
}

// StripeService returns the stripe service associated with the client.
func (c *Client) StripeService() Processes {
	c.services.stripe.initialize()
	return &c.services.stripe
}
