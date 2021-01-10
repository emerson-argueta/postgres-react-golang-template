package mockcommunitygoaltracker

import (
	"emersonargueta/m/v1/modules/communitygoaltracker/domain/achiever"
	"emersonargueta/m/v1/modules/communitygoaltracker/domain/goal"
)

// CommunitygoaltrackerProcesses is a mock of communitygoaltracker.Processes
type CommunitygoaltrackerProcesses struct {
	RegisterFn           func(a *achiever.Achiever) (*achiever.Achiever, error)
	LoginFn              func(email string, password string) (*achiever.Achiever, error)
	UpdateAchieverFn     func(a *achiever.Achiever) error
	UnRegisterFn         func(a *achiever.Achiever) error
	CreateGoalFn         func(g *goal.Goal) (*goal.Goal, error)
	UpdateGoalProgressFn func(achieverUUID string, goalID int64, progress int) (*goal.Goal, error)
	AbandonGoalFn        func(achieverUUID string, goalID int64) error
	DeleteGoalFn         func(achieverUUID string, goalID int64) error
}

// Register is a mock communitygoaltracker.Processes.Register
func (s *CommunitygoaltrackerProcesses) Register(a *achiever.Achiever) (res *achiever.Achiever, e error) {
	return s.RegisterFn(a)
}

// Login is a mock communitygoaltracker.Processes.Login
func (s *CommunitygoaltrackerProcesses) Login(email string, password string) (res *achiever.Achiever, e error) {
	return s.LoginFn(email, password)
}

// UpdateAchiever is a mock communitygoaltracker.Processes.UpdateAchiever
func (s *CommunitygoaltrackerProcesses) UpdateAchiever(a *achiever.Achiever) (e error) {
	return s.UpdateAchieverFn(a)
}

// UnRegister is a mock communitygoaltracker.Processes.UnRegister
func (s *CommunitygoaltrackerProcesses) UnRegister(a *achiever.Achiever) (e error) {
	return s.UnRegisterFn(a)
}

// CreateGoal is a mock communitygoaltracker.Processes.CreateGoal
func (s *CommunitygoaltrackerProcesses) CreateGoal(g *goal.Goal) (res *goal.Goal, e error) {
	return s.CreateGoalFn(g)
}

// UpdateGoalProgress is a mock communitygoaltracker.Processes.UpdateGoalProgress
func (s *CommunitygoaltrackerProcesses) UpdateGoalProgress(achieverUUID string, goalID int64, progress int) (res *goal.Goal, e error) {
	return s.UpdateGoalProgressFn(achieverUUID, goalID, progress)
}

// AbandonGoal is a mock communitygoaltracker.Processes.AbandonGoal
func (s *CommunitygoaltrackerProcesses) AbandonGoal(achieverUUID string, goalID int64) (e error) {
	return s.AbandonGoalFn(achieverUUID, goalID)
}

// DeleteGoal is a mock communitygoaltracker.Processes.DeleteGoal
func (s *CommunitygoaltrackerProcesses) DeleteGoal(achieverUUID string, goalID int64) (e error) {
	return s.DeleteGoalFn(achieverUUID, goalID)
}
