package user

// user is part of the identity subdomain to support core domains which need it

//User model
type User struct {
	UUID     *string  `db:"uuid" dbignoreinsert:"" json:"uuid"`
	Email    *string  `db:"email" json:"email"`
	Password *string  `db:"password" json:"password"`
	Domains  *Domains `db:"domains" json:"domains"`
}

// Domains is a map where the key represents the domains's id
type Domains map[int64]struct {
	Role *string `json:"role,omitempty"`
}

// Keys represent the domain's id
func (d *Domains) Keys() []int64 {
	keys := make([]int64, len(*d))

	i := 0
	for k := range *d {
		keys[i] = k
		i++
	}
	return keys
}

// Processes used to modify the user model.
type Processes interface {
	// CreateUser implementation must return ErrUserExists if the user exists.
	CreateUser(*User) (*User, error)
	// RetrieveUser implementation must return ErrUserNotFound if the user is not found.
	RetrieveUser(email string) (*User, error)
	// UpdateUser implementation must search user by uuid and return
	// ErrUserNotFound if user is not found. Must return ErrUserExists if
	// user's new email to be updated conflicts with another user's email.
	UpdateUser(*User) error
	// DeleteUser implementation should search the user by uuid before deleting
	// the user and must return ErrUserNotFound if the user does not exists.
	DeleteUser(uuid string) error
}
