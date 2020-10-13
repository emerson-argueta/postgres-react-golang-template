package authorization

import "emersonargueta/m/v1/config"

// Multiple service structure for client

// Client to the jwt authorization service.
type Client struct {
	config   *config.Config
	services services
}
type services struct {
	jwt jwtservice
}

// NewClient to the jwt authorization service.
func NewClient(config *config.Config) *Client {
	c := &Client{}

	c.services.jwt.client = c

	c.config = config

	return c
}

// JwtService returns the jwt authorization service associated with the client.
func (c *Client) JwtService() Processes { return &c.services.jwt }
