package database

import "github.com/jmoiron/sqlx"

// Client for any database
type Client interface {
	Open() error
	Close() error
	Query() Query
	DB() *sqlx.DB
	Schema() string
}
