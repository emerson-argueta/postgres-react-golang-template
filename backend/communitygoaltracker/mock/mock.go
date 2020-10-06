package mockcommunitygoaltracker

import (
	"emersonargueta/m/v1/communitygoaltracker/achiever"
	"emersonargueta/m/v1/communitygoaltracker/goal"
)

// CommunityGoalTrackerService is a mock of communitygoaltracker.Service
type CommunityGoalTrackerService struct {
	CreateGoalFn         func(*achiever.Achiever, *goal.Goal) (*goal.Goal, error)
	UpdateGoalProgressFn func(*achiever.Achiever, *goal.Goal) error
	AbandonGoalFn        func(*achiever.Achiever, *goal.Goal) error
	DeleteGoalFn         func(*achiever.Achiever, *goal.Goal) error
	LoginFn              func(email string, password string) error
	RegisterFn           func(*achiever.Achiever) (*achiever.Achiever, error)
	UpdateFn             func(*achiever.Achiever) error
}

// CreateGoal mocks communitygoaltracker.Service.CreateGoal
func (s *CommunityGoalTrackerService) CreateGoal(a *achiever.Achiever, g *goal.Goal) (*goal.Goal, error) {
	return s.CreateGoalFn(a, g)
}

// UpdateGoalProgress mocks communitygoaltracker.Service.UpdateGoalProgress
func (s *CommunityGoalTrackerService) UpdateGoalProgress(a *achiever.Achiever, g *goal.Goal) error {
	return s.UpdateGoalProgressFn(a, g)
}

// AbandonGoal mocks communitygoaltracker.Service.AbandonGoal
func (s *CommunityGoalTrackerService) AbandonGoal(a *achiever.Achiever, g *goal.Goal) error {
	return s.AbandonGoalFn(a, g)
}

// DeleteGoal mocks communitygoaltracker.Service.DeleteGoal
func (s *CommunityGoalTrackerService) DeleteGoal(a *achiever.Achiever, g *goal.Goal) error {
	return s.DeleteGoalFn(a, g)
}

// Login mocks communitygoaltracker.Service.Login
func (s *CommunityGoalTrackerService) Login(email string, password string) error {
	return s.LoginFn(email, password)
}

// Register mocks communitygoaltracker.Service.Register
func (s *CommunityGoalTrackerService) Register(a *achiever.Achiever) (*achiever.Achiever, error) {
	return s.RegisterFn(a)
}

// Update mocks communitygoaltracker.Service.Update
func (s *CommunityGoalTrackerService) Update(a *achiever.Achiever) error {
	return s.UpdateFn(a)
}
