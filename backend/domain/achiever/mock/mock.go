package mockachiever

import (
	"emersonargueta/m/v1/domain/achiever"
	"emersonargueta/m/v1/domain/goal"
)

//AchieverService is a mock of achiever.Service
type AchieverService struct {
	CreateGoalFn         func(*achiever.Achiever, *goal.Goal) (*goal.Goal, error)
	UpdateGoalProgressFn func(*achiever.Achiever, *goal.Goal) error
	AbandonGoalFn        func(*achiever.Achiever, *goal.Goal) error
	DeleteGoalFn         func(*achiever.Achiever, *goal.Goal) error
}

// CreateGoal mocks achiever.Service.CreateGoal
func (s *AchieverService) CreateGoal(a *achiever.Achiever, g *goal.Goal) (*goal.Goal, error) {
	return s.CreateGoalFn(a, g)
}

// UpdateGoalProgress mocks achiever.Service.UpdateGoalProgress
func (s *AchieverService) UpdateGoalProgress(a *achiever.Achiever, g *goal.Goal) error {
	return s.UpdateGoalProgressFn(a, g)
}

// AbandonGoal mocks achiever.Service.AbandonGoal
func (s *AchieverService) AbandonGoal(a *achiever.Achiever, g *goal.Goal) error {
	return s.AbandonGoalFn(a, g)
}

// DeleteGoal mocks achiever.Service.DeleteGoal
func (s *AchieverService) DeleteGoal(a *achiever.Achiever, g *goal.Goal) error {
	return s.DeleteGoalFn(a, g)
}
