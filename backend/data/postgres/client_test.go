package postgres_test

import (
	"time"

	"trustdonations.org/m/v2/data/postgres"
)

// Now is the mocked current time for testing.
var Now = time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC)

type Client struct {
	*postgres.Client
}

func NewClient() *Client {
	c := &Client{Client: postgres.NewClient()}
	c.Now = func() time.Time { return Now }

	return c
}

// MustOpenClient returns an new, open instance of Client.
func MustOpenClient() *Client {
	c := NewClient()
	if err := c.Open(); err != nil {
		panic(err)
	}
	return c
}

// Close closes the client and removes the underlying database.
func (c *Client) Close() error {
	return c.Client.Close()

}
