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
		return res, ErrAchieverInvalidRole
	}

	return res, e

}

// NewRole creates a new role, user role by default. Returns error if role is
// invalid.
func NewRole(s *string) (res Role, e error) {
	if s == nil {
		return UserRole, nil
	}
	return ToRole(*s)
}
