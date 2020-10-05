package user

// user is part of the identity subdomain to support core domains which need it

//User model
type User struct {
	UUID     *string  `json:"uuid"`
	Email    *string  `db:"email" json:"email"`
	Password *string  `db:"password" json:"password"`
	Domains  *Domains `db:"domains" json:"domains"`
}

//Domain model represents domains which users can be part of
type Domain struct {
	ID   *int64  `db:"id" dbignoreinsert:"" json:"id"`
	Name *string `db:"name" json:"name"`
}

// Domains is a map where the key represents the domains's id
type Domains map[int64]struct {
	Role *string `json:"role,omitempty"`
}

// Client creates a connection to a service. In this case, an user service.
type Client interface {
	Service() Service
}

// Service provides processes that can be achieved by user.
type Service interface {
	Register(*User) error
	Retrieve(u *User, byEmail bool) (*User, error)
	Update(u *User, byEmail bool) error
	UnRegister(u *User, byEmail bool) error
	LookUpDomain(*Domain) (*Domain, error)
}
