package subscriptions

import (
	"emersonargueta/m/v1/modules/communitygoaltracker/infrastructure/persistence"
	"emersonargueta/m/v1/modules/communitygoaltracker/usecase"
)

// Subscriptions holds all subscriptions
type Subscriptions struct {
	AfterUserCreated AfterUserCreated
}

// New subscriptions for communitygoaltracker domain
func New() *Subscriptions {
	subscriptions := &Subscriptions{}

	cgtRepos := persistence.CommunitygoaltrackerRepos
	CreateAchiever := usecase.NewCreateAchieverUsecase(&cgtRepos.Achiever)
	subscriptions.AfterUserCreated = NewAfterUserCreated(CreateAchiever)

	return subscriptions
}
