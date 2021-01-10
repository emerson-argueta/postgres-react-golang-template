package postgres

import (
	"database/sql"
	"emersonargueta/m/v1/shared/infrastructure"
	"emersonargueta/m/v1/shared/infrastructure/database"
	"fmt"
	"log"
	"os"

	"github.com/gchaincl/sqlhooks"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	_ "github.com/lib/pq" // here
)

// DatabaseClient is a new postgres database client
var DatabaseClient = new()

func init() {
	if err := DatabaseClient.Open(); err != nil {
		panic(err)
	}
}

// Client inherits from database client
type Client interface {
	database.Client
}

type client struct {
	config *infrastructure.Config
	query  database.Query
	db     *sqlx.DB
}

func new() Client {
	c := &client{}

	query := NewQuery()
	config := infrastructure.NewConfig

	c.config = config
	c.query = query

	return c
}

// Open and initializes the PostgreSQL database.
func (c *client) Open() error {

	sslmode := "require"
	if c.config.Environment == "DEV" {
		sslmode = "disable"
	}
	connectionStr := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s search_path=%s sslmode=%s",
		c.config.Database.Host, c.config.Database.Port, c.config.Database.User,
		c.config.Database.Password, c.config.Database.DB, c.config.Database.Schema, sslmode,
	)

	// pq.Driver
	driverName := c.config.Database.Dialect + "WithHooks"
	sql.Register(driverName, sqlhooks.Wrap(&pq.Driver{}, &Hooks{}))

	// driverName := c.config.Database.Dialect
	sqlDatabase, err := sql.Open(driverName, connectionStr)

	database := sqlx.NewDb(sqlDatabase, c.config.Database.Dialect)

	if err != nil {
		return err
	}
	c.db = database

	if err = database.Ping(); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	return nil
}

// Close closes then underlying postgres database.
func (c *client) Close() error {
	if c.db != nil {
		return c.db.Close()
	}
	return nil
}

func (c *client) Query() database.Query {
	return c.query
}
func (c *client) DB() *sqlx.DB {
	return c.db
}
func (c *client) Schema() string {
	return c.config.Database.Schema
}
