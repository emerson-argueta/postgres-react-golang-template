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

var _ Processes = &Service{}

// Service exposes identity domain processes.
type Service struct {
	client *Client
	// Processes from domain models used internally in identity domain processes
	User   user.Processes
	Domain domain.Processes
}

// Processes an interface for identity domain processes.
type Processes interface {
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
func (s *Service) RegisterUser(u *user.User) (res *user.User, e error) {
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

	return s.User.CreateUser(u)
}

// LoginUser using the following business logic:
// Check if user exists.
// Check if password matches.
func (s *Service) LoginUser(email string, password string) (res *user.User, e error) {
	res, e = s.User.RetrieveUser(email)
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
func (s *Service) UpdateUser(u *user.User) (e error) {
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

	return s.User.UpdateUser(u)
}

// UnRegisterUser using the following business logic
// Validate email and password.
// Delete the user.
func (s *Service) UnRegisterUser(u *user.User) (e error) {
	if u.UUID == nil || u.Email == nil || u.Password == nil {
		return ErrUserIncompleteDetails
	}

	if exists, e := s.User.RetrieveUser(*u.UUID); e != nil {
		return e
	} else if *exists.Email != *u.Email || bcrypt.CompareHashAndPassword([]byte(*exists.Password), []byte(*u.Password)) == bcrypt.ErrMismatchedHashAndPassword {
		return ErrUserIncorrectCredentials
	}

	return s.User.DeleteUser(*u.UUID)
}

// AddDomain using the following business logic:
// Verify domain name.
// Create the domain.
func (s *Service) AddDomain(d *domain.Domain) (res *domain.Domain, e error) {
	if d.Name == nil {
		return nil, ErrDomainIncompleteDetails
	}

	return s.Domain.CreateDomain(d)
}

// LookupDomain using the following business logic:
// Retrieve the domain.
func (s *Service) LookupDomain(name string) (res *domain.Domain, e error) {
	return s.Domain.RetrieveDomain(name)
}

// UpdateDomain using the following business logic
// Verify the domain name and id
func (s *Service) UpdateDomain(d *domain.Domain) (e error) {
	if d.Name == nil || d.ID == nil {
		return ErrDomainIncompleteDetails
	}
	return s.Domain.UpdateDomain(d)
}

//RemoveDomain using the following business logic:
// Delete the domain.
func (s *Service) RemoveDomain(id int64) (e error) {
	return s.Domain.DeleteDomain(id)
}
