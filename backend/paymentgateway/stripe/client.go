package stripe

import (
	"emersonargueta/m/v1/config"
	"emersonargueta/m/v1/paymentgateway"

	"github.com/stripe/stripe-go/client"
)

// Client to the stripe service.
type Client struct {
	config  *config.Config
	service service
	stripe  *client.API
}

// NewClient creates connection to strinpe service.
func NewClient(config *config.Config) *Client {
	c := &Client{}

	c.service.client = c
	c.config = config

	return c
}

// Initialize the underlying stripe client.
func (c *Client) Initialize() {
	apiKey := c.config.PaymentGateway.APIKey
	stripe := client.New(apiKey, nil)

	c.stripe = stripe
}

// Service returns the stripe service associated with the client.
func (c *Client) Service() paymentgateway.Processes { return &c.service }
