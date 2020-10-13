package achiever

import "errors"

// Role represents access role for an achiever within the communitygoaltracker
// domain
type Role int

const (
	// UserRole has access to manage their own information
	UserRole Role = iota
	// AdministratorRole has access to manage users within their own goals
	AdministratorRole
)

func (s Role) String() (res string, e error) {
	if s < UserRole || s > AdministratorRole {
		return res, errors.New("could not convert to Role to string")
	}
	return [...]string{"user", "administrator"}[s], e

}

// ToRole Converts string to Role enum
func ToRole(s string) (res Role, e error) {
	res, ok := map[string]Role{"user": UserRole, "administrator": AdministratorRole}[s]
	if !ok {
		e = errors.New("Could not covert to Role")
	}

	return res, e

}

// Achiever model in the communinty_goal_tracker domain
type Achiever struct {
	UUID      *string `db:"uuid" json:"uuid"`
	Role      *Role   `db:"role" json:"role"`
	Firstname *string `db:"firstname" json:"firstname"`
	Lastname  *string `db:"lastname" json:"lastname"`
	Address   *string `db:"address" json:"address"`
	Phone     *string `db:"phone" json:"phone"`
	Goals     *Goals  `db:"goals" json:"goals,omitempty"`
	Email     *string `json:"email,omitempty"`
	Password  *string `json:"password,omitempty"`
}

// Goals represents a slice goals ids for an achiever
type Goals map[int64]bool

// Keys represent the goal ids
func (g *Goals) Keys() []int64 {
	keys := make([]int64, len(*g))

	i := 0
	for k := range *g {
		keys[i] = k
		i++
	}
	return keys
}

// Processes used to modify the achiever model.
type Processes interface {
	// CreateAchiever implementation must return ErrAchieverExists if achiever
	// exists.
	CreateAchiever(*Achiever) (*Achiever, error)
	// RetrieveAchiever implementation must return ErrAchieverNotFound if the
	// achiever is not found.
	RetrieveAchiever(uuid string) (*Achiever, error)
	// UpdateAchiever implementation must search achiever by uuid and return
	// ErrAchieverNotFound if achiever is not found.
	UpdateAchiever(*Achiever) error
	// DeleteAchiever implementation should search the achiever by uuid before
	// deleting the achiever and must return ErrAchieverNotFound if the achiever
	// does not exists.
	DeleteAchiever(uuid string) error
}
