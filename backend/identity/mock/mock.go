package mockuser

import (
	"emersonargueta/m/v1/identity/domain"
	"emersonargueta/m/v1/identity/user"
)

//IdentityService is a mock of identity.Service
type IdentityService struct {
	RegisterUserFn   func(*user.User) error
	LoginUserFn      func(email string, password string) (res *user.User, e error)
	UpdateUserFn     func(u *user.User) error
	UnRegisterUserFn func(u *user.User) error
	AddDomainFn      func(*domain.Domain) error
	LookUpDomainFn   func(*domain.Domain) (*domain.Domain, error)
	UpdateDomainFn   func(*domain.Domain) error
	RemoveDomainFn   func(*domain.Domain) error
}

// RegisterUser mocks identity.Service.RegisterUser
func (s *IdentityService) RegisterUser(u *user.User) error {
	return s.RegisterUserFn(u)
}

// LoginUser mocks identity.Service.LoginUser
func (s *IdentityService) LoginUser(email string, password string) (res *user.User, e error) {
	return s.LoginUserFn(email, password)
}

// UpdateUser mocks identity.Service.UpdateUser
func (s *IdentityService) UpdateUser(u *user.User) error {
	return s.UpdateUserFn(u)
}

// UnRegister mocks identity.Service.UnRegister
func (s *IdentityService) UnRegister(u *user.User) error {
	return s.UnRegisterUserFn(u)
}

// AddDomain mockes identity.Service.AddDomain
func (s *IdentityService) AddDomain(d *domain.Domain) error {
	return s.AddDomainFn(d)
}

// LookUpDomain mockes identity.Service.LookUpDomain
func (s *IdentityService) LookUpDomain(d *domain.Domain) (*domain.Domain, error) {
	return s.LookUpDomainFn(d)
}

// UpdateDomain mockes identity.Service.UpdateDomain
func (s *IdentityService) UpdateDomain(d *domain.Domain) error {
	return s.UpdateDomainFn(d)
}

// RemoveDomain mockes identity.Service.RemoveDomain
func (s *IdentityService) RemoveDomain(d *domain.Domain) error {
	return s.RemoveDomainFn(d)
}
