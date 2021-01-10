package persistence

import "emersonargueta/m/v1/shared/infrastructure/database/postgres"

const (
	// IdentitySchema used to group tables used in the identity domain
	IdentitySchema = "IDENTITY"
	// UserTable stores user information for the identity domain
	UserTable = "USER"
)

// IdentityRepos implementation
var IdentityRepos = new()

// Services represents the services that the persistence service provides
type Services struct {
	User User
}

func new() *Services {
	services := &Services{}
	connection := postgres.DatabaseClient

	services.User.Client = connection

	return services
}
