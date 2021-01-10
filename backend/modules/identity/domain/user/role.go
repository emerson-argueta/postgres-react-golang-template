package user

import "errors"

// Role represents access role for a user within the indentity
// domain
type Role int

const (
	// UserRole has access to manage their own information
	UserRole Role = iota
	// AdministratorRole has access to manage users within other domains
	AdministratorRole
)

func (r Role) String() (res string, e error) {
	if r < UserRole || r > AdministratorRole {
		return res, errors.New("could not convert to Role to string")
	}
	return [...]string{"user", "administrator"}[r], e

}

// ToRole Converts string to Role enum
func ToRole(r string) (res Role, e error) {
	res, ok := map[string]Role{"user": UserRole, "administrator": AdministratorRole}[r]
	if !ok {
		e = errors.New("Could not covert to Role")
	}

	return res, e

}

// NewRole creates a new role, user role by default
func NewRole(r *string) (res Role, e error) {
	if r == nil {
		return UserRole, nil
	}
	return ToRole(*r)
}
