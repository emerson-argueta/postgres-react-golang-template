package identity

import (
	"emersonargueta/m/v1/identity/domain"
	"emersonargueta/m/v1/identity/user"
	"emersonargueta/m/v1/validation"

	"golang.org/x/crypto/bcrypt"
)

// DomainName of this package
const DomainName = "identity"

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
	AddDomain(*domain.Domain) error
	LookupDomain(*domain.Domain) (*domain.Domain, error)
	UpdateDomain(*domain.Domain) error
	RemoveDomain(*domain.Domain) error
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

//AddDomain using the following business logic
func (i *Identity) AddDomain(d *domain.Domain) (e error) {
	return e
}

//LookupDomain using the following business logic
func (i *Identity) LookupDomain(d *domain.Domain) (res *domain.Domain, e error) {
	return res, e
}

//UpdateDomain using the following business logic
func (i *Identity) UpdateDomain(d *domain.Domain) (e error) {
	return e
}

//RemoveDomain using the following business logic
func (i *Identity) RemoveDomain(d *domain.Domain) (e error) {
	return e
}
