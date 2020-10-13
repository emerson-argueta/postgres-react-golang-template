package postgres

import (
	"fmt"
	"log"
	"os"
	"time"

	"emersonargueta/m/v1/communitygoaltracker/achiever"
	"emersonargueta/m/v1/communitygoaltracker/goal"
	"emersonargueta/m/v1/config"
	"emersonargueta/m/v1/identity/domain"
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

	services services

	db *sqlx.DB
}

// services represents the services that the postgres service provides
type services struct {
	user     UserService
	domain   DomainService
	achiever AchieverService
	goal     GoalService
}

// NewClient function
func NewClient(config *config.Config) *Client {
	c := &Client{Now: time.Now}

	c.services.user.client = c
	c.services.domain.client = c
	c.services.achiever.client = c
	c.services.goal.client = c

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
func (c *Client) UserService() user.Processes { return &c.services.user }

// DomainService returns the domain service associated with the client.
func (c *Client) DomainService() domain.Processes { return &c.services.domain }

// AchieverService returns the achiever service associated with the client.
func (c *Client) AchieverService() achiever.Processes { return &c.services.achiever }

// GoalService returns the goal service associated with the client.
func (c *Client) GoalService() goal.Processes { return &c.services.goal }
