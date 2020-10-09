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

// Service provides processes that can be achieved by user.
type Service interface {
	// CreateUser implementations should check if user exists before creating the user.
	CreateUser(*User) (*User, error)
	RetrieveUser(email string) (*User, error)
	// UpdateUser implementation should search the user by uuid before updating the user.
	UpdateUser(*User) error
	// DeleteUser implementation should search the user by uuid before deleting the user.
	DeleteUser(uuid string) error
}
