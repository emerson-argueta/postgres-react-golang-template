package user

// Access represents the access level a user within trustdonations org.
// The access level will determine how a user can manage other
// users within trustdonations.
type Access string

// NonRestricted enables a user to manage other users.
const NonRestricted = Access("non-restricted")

// Restricted disables a user to manage other users.
const Restricted = Access("restricted")

// Services is used for mapping to a postgres JSONB field
type Services map[string]*SingleService

// SingleService is used for services values
type SingleService struct {
	Role   *string `json:"role,omitempty"`
	Access *Access `json:"access,omitempty"`
}

//User from identity subdomain of trustdonations
type User struct {
	UUID     *string   `json:"uuid"`
	Email    *string   `db:"email" json:"email"`
	Password *string   `db:"password" json:"password"`
	Services *Services `db:"services" json:"services"`
}

// Client creates a connection to a service. In this case, an user service.
type Client interface {
	Service() Service
}

// Service provides functions that can be used for managing a user.
type Service interface {
	CreateManagementSession() error
	EndManagementSession() error
	Create(*User) error
	Read(u *User, byEmail bool) (*User, error)
	Update(u *User, byEmail bool) error
	Delete(u *User, byEmail bool) error
}
