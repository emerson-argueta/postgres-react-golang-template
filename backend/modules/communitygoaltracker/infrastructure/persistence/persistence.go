package persistence

import "emersonargueta/m/v1/shared/infrastructure/database/postgres"

const (
	// CommunitygoaltrackerSchema used to group tables used in the identity domain
	CommunitygoaltrackerSchema = "COMMUNITY_GOAL_TRACKER"
	// AchieverTable stores Achiever information for the communitygoaltracker domain
	AchieverTable = "ACHIEVER"
	// GoalTable stores goal information for the communitygoaltracker domain
	GoalTable = "GOAL"
)

// CommunitygoaltrackerRepos implementation
var CommunitygoaltrackerRepos = new()

// Services represents the services that the persistence service provides
type Services struct {
	Achiever Achiever
	Goal     Goal
}

func new() *Services {
	services := &Services{}
	connection := postgres.DatabaseClient

	services.Achiever.Client = connection
	services.Goal.Client = connection

	return services
}
