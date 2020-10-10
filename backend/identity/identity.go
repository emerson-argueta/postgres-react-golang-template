package identity

import (
	"emersonargueta/m/v1/identity/domain"
	"emersonargueta/m/v1/identity/user"
	"emersonargueta/m/v1/validation"

	"golang.org/x/crypto/bcrypt"
)

const (
	// DomainName of this package
	DomainName = "identity"
)

var _ Service = &Identity{}

// Identity exposes identity domain processes.
type Identity struct {
	client *Client
	// Services from domain models used internally in identity domain processes
	User   user.Service
	Domain domain.Service
}

// Service an interface for identity domain processes.
type Service interface {
	RegisterUser(*user.User) (*user.User, error)
	LoginUser(email string, password string) (*user.User, error)
	UpdateUser(*user.User) error
	UnRegisterUser(*user.User) error
	AddDomain(*domain.Domain) (*domain.Domain, error)
	LookupDomain(name string) (*domain.Domain, error)
	UpdateDomain(*domain.Domain) error
	RemoveDomain(id int64) error
}

// RegisterUser using the following business logic
// Verify email and password.
// Check if user exists.
// Create user with hash password.
func (i *Identity) RegisterUser(u *user.User) (res *user.User, e error) {
	if u.Email == nil || u.Password == nil {
		return nil, ErrUserIncompleteDetails
	}

	if e = validation.ValidateUserEmail(*u.Email); e != nil {
		return nil, e
	}
	if e = validation.ValidatePassword(*u.Password); e != nil {
		return nil, e
	}

	hash, e := bcrypt.GenerateFromPassword([]byte(*u.Password), bcrypt.DefaultCost)
	hashString := string(hash)
	u.Password = &hashString

	return i.User.CreateUser(u)
}

// LoginUser using the following business logic:
// Check if user exists.
// Check if password matches.
func (i *Identity) LoginUser(email string, password string) (res *user.User, e error) {
	res, e = i.User.RetrieveUser(email)
	if e != nil {
		return nil, e
	}

	e = bcrypt.CompareHashAndPassword([]byte(*res.Password), []byte(password))
	if e != nil && e == bcrypt.ErrMismatchedHashAndPassword {
		return nil, ErrUserIncorrectCredentials
	} else if e != nil {
		return nil, e
	}

	return res, e
}

// UpdateUser using the following business logic
// Validate email and password
// Update the user searching by uuid.
func (i *Identity) UpdateUser(u *user.User) (e error) {
	if u.UUID == nil {
		return ErrUserIncompleteDetails
	}

	if u.Email != nil {
		if e = validation.ValidateUserEmail(*u.Email); e != nil {
			return e
		}
	}
	if u.Password != nil {
		if e = validation.ValidatePassword(*u.Password); e != nil {
			return e
		}
	}

	return i.User.UpdateUser(u)
}

// UnRegisterUser using the following business logic
// Validate email and password.
// Delete the user.
func (i *Identity) UnRegisterUser(u *user.User) (e error) {
	if u.UUID == nil || u.Email == nil || u.Password == nil {
		return ErrUserIncompleteDetails
	}

	if exists, e := i.User.RetrieveUser(*u.UUID); e != nil {
		return e
	} else if *exists.Email != *u.Email || bcrypt.CompareHashAndPassword([]byte(*exists.Password), []byte(*u.Password)) == bcrypt.ErrMismatchedHashAndPassword {
		return ErrUserIncorrectCredentials
	}

	return i.User.DeleteUser(*u.UUID)
}

// AddDomain using the following business logic:
// Verify domain name.
// Create the domain.
func (i *Identity) AddDomain(d *domain.Domain) (res *domain.Domain, e error) {
	if d.Name == nil {
		return nil, ErrDomainIncompleteDetails
	}

	return i.Domain.CreateDomain(d)
}

// LookupDomain using the following business logic:
// Retrieve the domain.
func (i *Identity) LookupDomain(name string) (res *domain.Domain, e error) {
	return i.Domain.RetrieveDomain(name)
}

// UpdateDomain using the following business logic
// Verify the domain name and id
func (i *Identity) UpdateDomain(d *domain.Domain) (e error) {
	if d.Name == nil || d.ID == nil {
		return ErrDomainIncompleteDetails
	}
	return i.Domain.UpdateDomain(d)
}

//RemoveDomain using the following business logic:
// Delete the domain.
func (i *Identity) RemoveDomain(id int64) (e error) {
	return i.Domain.DeleteDomain(id)
}
