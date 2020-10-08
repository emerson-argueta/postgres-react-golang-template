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

// Communitygoaltracker exposes Communitygoaltracker domain processes.
type Communitygoaltracker struct {
	client *Client
	// Communitygoaltacker services used internally in processes.
	Achiever achiever.Service
	Goal     goal.Service
	SupportingServices
}

// SupportingServices from supporting domains used internally in communitygaoltracker processes.
type SupportingServices struct {
	identity.Services
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
	gg, e := cgt.deleteAchieverFromGoals(*a.Goals, *a.UUID)
	if e != nil {
		return e
	}

	if e = cgt.Achiever.DeleteAchiever(a); e != nil {
		return e
	}

	return cgt.Goal.UpdateGoals(gg)
}
func (cgt *Communitygoaltracker) deleteAchieverFromGoals(achieverGoals achiever.Goals, achieverUUID string) (res []*goal.Goal, e error) {
	goalIDs := achieverGoals
	gg, e := cgt.Goal.RetrieveGoals(goalIDs)
	for _, g := range gg {
		_, ok := (*g.Achievers)[achieverUUID]
		if ok {
			delete((*g.Achievers), achieverUUID)
		}
	}

	return gg, e
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

	e = cgt.UpdateAchiever(&achiever.Achiever{UUID: &achieverUUID, Goals: &achiever.Goals{*res.ID}})

	return res, e
}

// UpdateGoalProgress using the following business logic
func (cgt *Communitygoaltracker) UpdateGoalProgress(achieverUUID, g *goal.Goal) (e error) {
	return e
}

// AbandonGoal using the following business logic
func (cgt *Communitygoaltracker) AbandonGoal(a *achiever.Achiever, g *goal.Goal) (e error) {
	return e
}
