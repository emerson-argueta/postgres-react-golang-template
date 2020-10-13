package communitygoaltracker

import (
	"emersonargueta/m/v1/communitygoaltracker/achiever"
	"emersonargueta/m/v1/communitygoaltracker/goal"
	"emersonargueta/m/v1/identity"
	"emersonargueta/m/v1/identity/user"
)

const (
	// DomainName of this package
	DomainName = "community-goal-tracker"
)

var _ Processes = &Service{}

// Service exposes domain and model processes
type Service struct {
	client *Client
	Processes
	Achiever achiever.Processes
	Goal     goal.Processes
	Identity identity.Service
}

// Processes for communitygoaltracker.
type Processes interface {
	Register(a *achiever.Achiever) (*achiever.Achiever, error)
	Login(email string, password string) (*achiever.Achiever, error)
	UpdateAchiever(a *achiever.Achiever) error
	UnRegister(a *achiever.Achiever) error
	CreateGoal(g *goal.Goal) (*goal.Goal, error)
	UpdateGoalProgress(achieverUUID string, goalID int64, progress int) (*goal.Goal, error)
	AbandonGoal(achieverUUID string, goalID int64) error
	DeleteGoal(achieverUUID string, goalID int64) error
}

// Register using the following business logic:
// Using the identity subdomain Service register the achiever.
// If there is an existing user, retrieve the user.
// Create the achiever.
func (s *Service) Register(a *achiever.Achiever) (res *achiever.Achiever, e error) {

	d, e := s.Identity.LookupDomain(DomainName)
	if e != nil {
		return nil, e
	}

	domains := make(user.Domains, 0)
	role := a.Role.String()
	domains[*d.ID] = struct {
		Role *string "json:\"role,omitempty\""
	}{Role: &role}

	u, e := s.Identity.RegisterUser(&user.User{Email: a.Email, Password: a.Password, Domains: &domains})
	if e == identity.ErrUserExists {
		u, e = s.Identity.LoginUser(*a.Email, *a.Password)
	} else if e != nil {
		return nil, e
	}

	a.UUID = u.UUID

	return s.Achiever.CreateAchiever(a)
}

// Login using the following business logic: Using the identity subdomain
// Service login the achiever using the provided achiever. If the identity
// Service is unable to process the request return the error
func (s *Service) Login(email string, password string) (res *achiever.Achiever, e error) {
	u, e := s.Identity.LoginUser(email, password)
	if e != nil {
		return nil, e
	}

	return s.Achiever.RetrieveAchiever(*u.UUID)
}

// UpdateAchiever using the following business logic: If email or password is
// updated update the identity domain user then update the communitygoaltracker
// achiever.
func (s *Service) UpdateAchiever(a *achiever.Achiever) (e error) {
	if a.UUID == nil {
		return ErrAchieverIncompleteDetails
	}

	if a.Email != nil || a.Password != nil {
		e = s.Identity.UpdateUser(&user.User{UUID: a.UUID, Email: a.Email, Password: a.Password})
	}
	if e != nil {
		return e
	}

	return s.Achiever.UpdateAchiever(a)
}

//UnRegister using the following business logic: delete the achiever and remove
//achiever from any goals they created.
func (s *Service) UnRegister(a *achiever.Achiever) (e error) {
	gg, e := s.removeAchieverFromGoals(*a.Goals, *a.UUID)
	if e != nil {
		return e
	}

	if e = s.Achiever.DeleteAchiever(*a.UUID); e != nil {
		return e
	}

	return s.Goal.UpdateGoals(gg)
}

// CreateGoal using the following business logic
// Create a goal with achiever as part of goal's achievers.
// Add goal to achiever's goals and update the achiever.
//
func (s *Service) CreateGoal(g *goal.Goal) (res *goal.Goal, e error) {
	if g.Achievers == nil {
		return nil, ErrGoalIncompleteDetails
	}
	achieverUUID := g.Achievers.Keys()[0]

	res, e = s.Goal.CreateGoal(g)
	e = s.addGoalToAchiever(&achiever.Achiever{UUID: &achieverUUID}, *res.ID)

	return res, e
}

// UpdateGoalProgress using the following business logic: Retrieve the goal and
// update the progress for the goal's achiever. Valid progress values are
// between 0 and 100 which can be interpreted as 0 percent to 100 percent. If
// the progress is 100 then update the state for the goal's achiever to
// complete.
func (s *Service) UpdateGoalProgress(achieverUUID string, goalID int64, progress int) (res *goal.Goal, e error) {
	if progress < 0 || progress > 100 {
		return nil, ErrGoalInvalidProgress
	}

	g, e := s.Goal.RetrieveGoal(goalID)
	if g.Achievers == nil {
		return nil, ErrGoalWithNoAchievers
	}

	if _, ok := (*g.Achievers)[achieverUUID]; !ok {
		return nil, ErrGoalNotFound
	}

	(*g.Achievers)[achieverUUID].Progress = &progress
	if progress == 100 {
		state := goal.Completed
		(*g.Achievers)[achieverUUID].State = &state
	}
	e = s.Goal.UpdateGoal(g)
	if e != nil {
		return nil, e
	}

	return g, e
}

// AbandonGoal using the following business logic: Retrieve the goal and
// update the state for the goal's achiever.
func (s *Service) AbandonGoal(achieverUUID string, goalID int64) (e error) {

	g, e := s.Goal.RetrieveGoal(goalID)
	if g.Achievers == nil {
		return ErrGoalWithNoAchievers
	}

	if _, ok := (*g.Achievers)[achieverUUID]; !ok {
		return ErrGoalNotFound

	}

	state := goal.Abondoned
	(*g.Achievers)[achieverUUID].State = &state

	return s.Goal.UpdateGoal(g)
}

// DeleteGoal using the following business logic: Retrieve the goal and if the
// goal has now achievers except for the one deleting then delete the goal.
func (s *Service) DeleteGoal(achieverUUID string, goalID int64) (e error) {

	g, e := s.Goal.RetrieveGoal(goalID)
	if g.Achievers == nil {
		return ErrGoalWithNoAchievers
	}

	if _, ok := (*g.Achievers)[achieverUUID]; !ok {
		return ErrGoalNotFound
	}
	if len((*g.Achievers).Keys()) > 1 {
		return ErrGoalCannotDelete
	}

	return s.Goal.DeleteGoal(goalID)
}

func (s *Service) removeGoalFromAchiever(a *achiever.Achiever, goalID int64) (e error) {
	a, e = s.Achiever.RetrieveAchiever(*a.UUID)
	if a.Goals == nil {
		return nil
	}

	_, ok := (*a.Goals)[goalID]
	if ok {
		delete((*a.Goals), goalID)
	}

	return s.Achiever.UpdateAchiever(a)
}
func (s *Service) addGoalToAchiever(a *achiever.Achiever, goalID int64) (e error) {
	a, e = s.Achiever.RetrieveAchiever(*a.UUID)
	if a.Goals == nil {
		achieverGoals := make(achiever.Goals)
		a.Goals = &achieverGoals
	}
	(*a.Goals)[goalID] = true

	return s.Achiever.UpdateAchiever(a)
}
func (s *Service) removeAchieverFromGoals(achieverGoals achiever.Goals, achieverUUID string) (res []*goal.Goal, e error) {
	goalIDs := achieverGoals.Keys()
	gg, e := s.Goal.RetrieveGoals(goalIDs)
	for _, g := range gg {
		_, ok := (*g.Achievers)[achieverUUID]
		if ok {
			delete((*g.Achievers), achieverUUID)
		}
	}

	return gg, e
}
