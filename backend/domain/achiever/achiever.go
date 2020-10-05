package achiever

import "emersonargueta/m/v1/domain/goal"

const UserRole = "user"
const AdministratorRole = "administrator"

// Achiever model in the communinty_goal_tracker domain
type Achiever struct {
	UUID      *string `db:"uuid" json:"uuid"`
	Firstname *string `db:"firstname" json:"firstname"`
	Lastname  *string `db:"lastname" json:"lastname"`
	Address   *string `db:"address" json:"address"`
	Phone     *string `db:"phone" json:"phone"`
	Goals     *Goals  `db:"goals" json:"goals,omitempty"`
	Email     *string `json:"email,omitempty"`
	Password  *string `json:"password,omitempty"`
}

// Goals represents a slice goals ids for an achiever
type Goals []int64

// Client creates a connection to a service. In this case, an administrator service.
type Client interface {
	Service() Service
}

// Service provides processes that can be achieved by an achiever.
type Service interface {
	CreateGoal(*Achiever, *goal.Goal) (*goal.Goal, error)
	UpdateGoalProgress(*Achiever, *goal.Goal) error
	AbandonGoal(*Achiever, *goal.Goal) error
	DeleteGoal(*Achiever, *goal.Goal) error
}
