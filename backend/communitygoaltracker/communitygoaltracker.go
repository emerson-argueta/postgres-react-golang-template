package communitygoaltracker

import (
	"emersonargueta/m/v1/communitygoaltracker/achiever"
	"emersonargueta/m/v1/communitygoaltracker/goal"
	"emersonargueta/m/v1/identity"
	"emersonargueta/m/v1/identity/domain"
	"emersonargueta/m/v1/identity/user"
)

// DomainName of this package
const DomainName = "community-goal-tracker"

var _ Service = &Communitygoaltracker{}

// Communitygoaltracker exposes Communitygoaltracker domain processes.
type Communitygoaltracker struct {
	client *Client
	Service
	// Communitygoaltacker services used internally in processes.
	Achiever achiever.Service
	Goal     goal.Service
	// From identity domain used internally in communitygaoltracker processes.
	Identity identity.Service
}

// Service an interface for communitygoaltracker domain processes.
type Service interface {
	Register(a *achiever.Achiever) (*achiever.Achiever, error)
	Login(email string, password string) (*achiever.Achiever, error)
	UpdateAchiever(a *achiever.Achiever) error
	UnRegister(a *achiever.Achiever) error
	CreateGoal(g *goal.Goal) (*goal.Goal, error)
	UpdateGoalProgress(achieverUUID string, goalID int64, progress int) error
	AbandonGoal(achieverUUID string, goalID int64) error
	DeleteGoal(achieverUUID string, goalID int64) error
}

// Register using the following business logic:
// Using the identity subdomain service register the achiever using the
// provided achiever.
// If the identity service is unable to process the request return the error
func (cgt *Communitygoaltracker) Register(a *achiever.Achiever) (res *achiever.Achiever, e error) {
	if a.Email == nil || a.Password == nil {
		return nil, ErrIncompleteAchieverDetails
	}

	domainName := DomainName
	d, e := cgt.Identity.LookupDomain(&domain.Domain{Name: &domainName})
	if e != nil {
		return nil, e
	}

	domains := make(user.Domains, 0)
	role := a.Role.String()
	domains[*d.ID] = struct {
		Role *string "json:\"role,omitempty\""
	}{Role: &role}

	u, e := cgt.Identity.RegisterUser(&user.User{Email: a.Email, Password: a.Password, Domains: &domains})
	if e != nil {
		return nil, e
	}

	a.UUID = u.UUID

	return cgt.Achiever.CreateAchiever(a)
}

// Login using the following business logic: Using the identity subdomain
// service login the achiever using the provided achiever. If the identity
// service is unable to process the request return the error
func (cgt *Communitygoaltracker) Login(email string, password string) (res *achiever.Achiever, e error) {
	u, e := cgt.Identity.LoginUser(email, password)
	if e != nil {
		return nil, e
	}

	return cgt.Achiever.RetrieveAchiever(*u.UUID)
}

// UpdateAchiever using the following business logic: If email or password is
// updated update the identity user then update the communitygoaltracker
// achiever.
func (cgt *Communitygoaltracker) UpdateAchiever(a *achiever.Achiever) (e error) {
	if a.Email != nil || a.Password != nil {
		e = cgt.Identity.UpdateUser(&user.User{UUID: a.UUID, Email: a.Email, Password: a.Password})
	}
	if e != nil {
		return e
	}

	return cgt.Achiever.UpdateAchiever(a)
}

//UnRegister using the following business logic: delete the achiever and remove
//achiever from any goals they created.
func (cgt *Communitygoaltracker) UnRegister(a *achiever.Achiever) (e error) {
	gg, e := cgt.removeAchieverFromGoals(*a.Goals, *a.UUID)
	if e != nil {
		return e
	}

	if e = cgt.Achiever.DeleteAchiever(a); e != nil {
		return e
	}

	return cgt.Goal.UpdateGoals(gg)
}

// CreateGoal using the following business logic
// Create a goal with achiever as part of goal's achievers.
// Add goal to achiever's goals and update the achiever.
//
func (cgt *Communitygoaltracker) CreateGoal(g *goal.Goal) (res *goal.Goal, e error) {
	if g.Achievers == nil {
		return nil, ErrIncompleteGoalDetails
	}
	achieverUUID := g.Achievers.Keys()[0]

	res, e = cgt.Goal.CreateGoal(g)
	e = cgt.addGoalToAchiever(&achiever.Achiever{UUID: &achieverUUID}, *res.ID)

	return res, e
}

// UpdateGoalProgress using the following business logic: Retrieve the goal and
// update the progress for the goal's achiever. Valid progress values are
// between 0 and 100 which can be interpreted as 0 percent to 100 percent. If
// the progress is 100 then update the state for the goal's achiever to
// complete.
func (cgt *Communitygoaltracker) UpdateGoalProgress(achieverUUID string, goalID int64, progress int) (e error) {
	if progress < 0 || progress > 100 {
		return ErrInvalidProgress
	}

	g, e := cgt.Goal.RetrieveGoal(goalID)
	if g.Achievers == nil {
		return ErrGoalWithNoAchievers
	}

	if _, ok := (*g.Achievers)[achieverUUID]; !ok {
		return ErrGoalNotFound
	}

	(*g.Achievers)[achieverUUID].Progress = &progress
	if progress == 100 {
		state := goal.Completed
		(*g.Achievers)[achieverUUID].State = &state
	}

	return cgt.Goal.UpdateGoal(g)
}

// AbandonGoal using the following business logic: Retrieve the goal and
// update the state for the goal's achiever.
func (cgt *Communitygoaltracker) AbandonGoal(achieverUUID string, goalID int64) (e error) {

	g, e := cgt.Goal.RetrieveGoal(goalID)
	if g.Achievers == nil {
		return ErrGoalWithNoAchievers
	}

	if _, ok := (*g.Achievers)[achieverUUID]; !ok {
		return ErrGoalNotFound

	}

	state := goal.Abondoned
	(*g.Achievers)[achieverUUID].State = &state

	return cgt.Goal.UpdateGoal(g)
}

// DeleteGoal using the following business logic: Retrieve the goal and if the
// goal has now achievers except for the one deleting then delete the goal.
func (cgt *Communitygoaltracker) DeleteGoal(achieverUUID string, goalID int64) (e error) {

	g, e := cgt.Goal.RetrieveGoal(goalID)
	if g.Achievers == nil {
		return ErrGoalWithNoAchievers
	}

	if _, ok := (*g.Achievers)[achieverUUID]; !ok {
		return ErrGoalNotFound
	}
	if len((*g.Achievers).Keys()) > 1 {
		return ErrCannotDeleteGoal
	}

	return cgt.Goal.DeleteGoal(goalID)
}

func (cgt *Communitygoaltracker) removeGoalFromAchiever(a *achiever.Achiever, goalID int64) (e error) {
	a, e = cgt.Achiever.RetrieveAchiever(*a.UUID)
	if a.Goals == nil {
		return nil
	}

	_, ok := (*a.Goals)[goalID]
	if ok {
		delete((*a.Goals), goalID)
	}

	return cgt.Achiever.UpdateAchiever(a)
}
func (cgt *Communitygoaltracker) addGoalToAchiever(a *achiever.Achiever, goalID int64) (e error) {
	a, e = cgt.Achiever.RetrieveAchiever(*a.UUID)
	if a.Goals == nil {
		achieverGoals := make(achiever.Goals)
		a.Goals = &achieverGoals
	}
	(*a.Goals)[goalID] = true

	return cgt.Achiever.UpdateAchiever(a)
}
func (cgt *Communitygoaltracker) removeAchieverFromGoals(achieverGoals achiever.Goals, achieverUUID string) (res []*goal.Goal, e error) {
	goalIDs := achieverGoals.Keys()
	gg, e := cgt.Goal.RetrieveGoals(goalIDs)
	for _, g := range gg {
		_, ok := (*g.Achievers)[achieverUUID]
		if ok {
			delete((*g.Achievers), achieverUUID)
		}
	}

	return gg, e
}
