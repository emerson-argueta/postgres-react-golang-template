package mockcommunitygoaltracker

import (
	"emersonargueta/m/v1/communitygoaltracker/achiever"
	"emersonargueta/m/v1/communitygoaltracker/goal"
)

// CommunityGoalTrackerService is a mock of communitygoaltracker.Service
type CommunityGoalTrackerService struct {
	RegisterFn           func(a *achiever.Achiever) (*achiever.Achiever, error)
	LoginFn              func(email string, password string) (*achiever.Achiever, error)
	UpdateAchieverFn     func(a *achiever.Achiever) error
	UnRegisterFn         func(a *achiever.Achiever) error
	CreateGoalFn         func(g *goal.Goal) (*goal.Goal, error)
	UpdateGoalProgressFn func(achieverUUID string, goalID int64, progress int) error
	AbandonGoalFn        func(achieverUUID string, goalID int64) error
	DeleteGoalFn         func(achieverUUID string, goalID int64) error
}

// Register is a mock communitygoaltracker.Service.Register
func (s *CommunityGoalTrackerService) Register(a *achiever.Achiever) (res *achiever.Achiever, e error) {
	return s.RegisterFn(a)
}

// Login is a mock communitygoaltracker.Service.Login
func (s *CommunityGoalTrackerService) Login(email string, password string) (res *achiever.Achiever, e error) {
	return s.LoginFn(email, password)
}

// UpdateAchiever is a mock communitygoaltracker.Service.UpdateAchiever
func (s *CommunityGoalTrackerService) UpdateAchiever(a *achiever.Achiever) (e error) {
	return s.UpdateAchieverFn(a)
}

// UnRegister is a mock communitygoaltracker.Service.UnRegister
func (s *CommunityGoalTrackerService) UnRegister(a *achiever.Achiever) (e error) {
	return s.UnRegisterFn(a)
}

// CreateGoal is a mock communitygoaltracker.Service.CreateGoal
func (s *CommunityGoalTrackerService) CreateGoal(g *goal.Goal) (res *goal.Goal, e error) {
	return s.CreateGoalFn(g)
}

// UpdateGoalProgress is a mock communitygoaltracker.Service.UpdateGoalProgress
func (s *CommunityGoalTrackerService) UpdateGoalProgress(achieverUUID string, goalID int64, progress int) (e error) {
	return s.UpdateGoalProgressFn(achieverUUID, goalID, progress)
}

// AbandonGoal is a mock communitygoaltracker.Service.AbandonGoal
func (s *CommunityGoalTrackerService) AbandonGoal(achieverUUID string, goalID int64) (e error) {
	return s.AbandonGoalFn(achieverUUID, goalID)
}

// DeleteGoal is a mock communitygoaltracker.Service.DeleteGoal
func (s *CommunityGoalTrackerService) DeleteGoal(achieverUUID string, goalID int64) (e error) {
	return s.DeleteGoalFn(achieverUUID, goalID)
}
