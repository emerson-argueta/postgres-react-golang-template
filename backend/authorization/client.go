package authorization

import "emersonargueta/m/v1/shared/infrastructure"

// AuthorizationService has a multiple service structure for client
var AuthorizationService = newClient(infrastructure.GlobalConfig)

// Client to the jwt authorization service.
type Client struct {
	config   *infrastructure.Config
	services services
}
type services struct {
	jwt jwtservice
}

// NewClient to the jwt authorization service.
func newClient(config *infrastructure.Config) *Client {
	c := &Client{}

	c.services.jwt.client = c

	c.config = config

	return c
}

// JwtService returns the jwt authorization service associated with the client.
func (c *Client) JwtService() Processes { return &c.services.jwt }
