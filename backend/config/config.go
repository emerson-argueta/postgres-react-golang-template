package config

import "emersonargueta/m/v1/env"

// Config is a struct that contains configuration variables.
type Config struct {
	Environment    string
	Host           string
	Port           string
	Database       *Database
	PaymentGateway *PaymentGateway
	Authorization  *Authorization
}

// Database is a struct that contains the DB's configuration variables.
type Database struct {
	Host     string
	Port     string
	User     string
	DB       string
	Schema   string
	Password string
}

// PaymentGateway is a struct that contains the PaymentGateway configuration
// variables.
type PaymentGateway struct {
	APIKey string
}

// Authorization is a struct that contains the Auth configuration variables.
type Authorization struct {
	APIKey string
	Secret string
}

// NewConfig creates a new Config struct. Will panic unsuccessful in getting
// configuration variables.
func NewConfig() *Config {
	env.CheckDotEnv()
	port := env.Get("PORT")
	host := env.Get("HOST")
	return &Config{
		Environment: env.MustGet("ENV"),
		Port:        port,
		Host:        host,
		Database: &Database{
			Host:     env.MustGet("DB_HOST"),
			Port:     env.MustGet("DB_PORT"),
			User:     env.MustGet("DB_USER"),
			DB:       env.MustGet("DB_DB"),
			Schema:   env.MustGet("DB_SCHEMA"),
			Password: env.MustGet("DB_PASSWORD"),
		},
		PaymentGateway: &PaymentGateway{
			APIKey: env.MustGet("PAYMENT_GATEWAY_API_KEY"),
		},
		Authorization: &Authorization{
			APIKey: env.MustGet("AUTH_API_KEY"),
			Secret: env.MustGet("SECRET"),
		},
	}
}
