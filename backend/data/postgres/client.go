package postgres

import (
	"fmt"
	"log"
	"os"
	"time"

	"emersonargueta/m/v1/config"
	"emersonargueta/m/v1/identity/user"

	"github.com/jmoiron/sqlx"
	// using postgres implementation of sqlx
	_ "github.com/lib/pq"
)

// Client represents a client to the underlying PostgreSQL database.
type Client struct {

	// Returns the current time.
	Now func() time.Time

	config *config.Config

	Services Services

	db          *sqlx.DB
	transaction *sqlx.Tx
}

// Services represents the services that the postgres service provides
type Services struct {
	User User
}

// NewClient function
func NewClient() *Client {
	c := &Client{Now: time.Now, transaction: nil}

	c.Services.User.client = c

	// get configuration stucts via .env file
	config, err := config.NewConfig()
	if err != nil {
		panic(err)
	}
	c.config = config

	return c
}

// Open and initializes the PostgreSQL database.
func (c *Client) Open() error {

	connectionStr := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s search_path=%s sslmode=disable",
		c.config.Database.Host, c.config.Database.Port, c.config.Database.User,
		c.config.Database.Password, c.config.Database.DB, c.config.Database.Schema,
	)

	db, err := sqlx.Open("postgres", connectionStr)
	if err != nil {
		return err
	}
	c.db = db

	if err = db.Ping(); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	return nil
}

// Close closes then underlying postgres database.
func (c *Client) Close() error {
	if c.db != nil {
		return c.db.Close()
	}
	return nil
}

// UserService returns the user service associated with the client.
func (c *Client) UserService() user.Service { return &c.Services.User }
