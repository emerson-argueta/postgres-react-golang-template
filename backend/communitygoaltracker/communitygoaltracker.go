package communitygoaltracker

import (
	"emersonargueta/m/v1/communitygoaltracker/achiever"
	"emersonargueta/m/v1/communitygoaltracker/goal"
)

// Client creates a connection to a service. In this case, a communitygoaltracker service.
type Client interface {
	Service() Service
}

// Service provides processes provided by the communitygoaltracker domain.
type Service interface {
	achiever.Service
	goal.Service
}

// Usecase for process logic.
type Usecase struct {
	Services Services
}

// Services used by usecase
type Services struct {
	Achiever achiever.Service
	Goal     goal.Service
}

// Register using the following business logic
func (uc *Usecase) Register(a *achiever.Achiever) (e error) {

	return e
}

//Login using the following business logic
func (uc *Usecase) Login(email string, password string) (res *achiever.Achiever, e error) {
	return res, e
}

//UpdateAchiever using the following business logic
func (uc *Usecase) UpdateAchiever(a *achiever.Achiever) (e error) {
	return e
}

//UnRegister using the following business logic
func (uc *Usecase) UnRegister(a *achiever.Achiever) (e error) {
	return e
}

// CreateGoal using the following business logic
func (uc *Usecase) CreateGoal(a *achiever.Achiever, g *goal.Goal) (res *goal.Goal, e error) {
	return res, e
}

// UpdateGoalProgress using the following business logic
func (uc *Usecase) UpdateGoalProgress(a *achiever.Achiever, g *goal.Goal) (e error) {
	return e
}

// AbandonGoal using the following business logic
func (uc *Usecase) AbandonGoal(a *achiever.Achiever, g *goal.Goal) (e error) {
	return e
}

// DeleteGoal using the following business logic
func (uc *Usecase) DeleteGoal(a *achiever.Achiever, g *goal.Goal) (e error) {
	return e
}
