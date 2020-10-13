package mockuser

import (
	"emersonargueta/m/v1/identity/domain"
	"emersonargueta/m/v1/identity/user"
)

//IdentityProcesses is a mock of identity.Processes
type IdentityProcesses struct {
	RegisterUserFn   func(*user.User) error
	LoginUserFn      func(email string, password string) (res *user.User, e error)
	UpdateUserFn     func(u *user.User) error
	UnRegisterUserFn func(u *user.User) error
	AddDomainFn      func(*domain.Domain) error
	LookUpDomainFn   func(*domain.Domain) (*domain.Domain, error)
	UpdateDomainFn   func(*domain.Domain) error
	RemoveDomainFn   func(*domain.Domain) error
}

// RegisterUser mocks identity.Processes.RegisterUser
func (s *IdentityProcesses) RegisterUser(u *user.User) error {
	return s.RegisterUserFn(u)
}

// LoginUser mocks identity.Processes.LoginUser
func (s *IdentityProcesses) LoginUser(email string, password string) (res *user.User, e error) {
	return s.LoginUserFn(email, password)
}

// UpdateUser mocks identity.Processes.UpdateUser
func (s *IdentityProcesses) UpdateUser(u *user.User) error {
	return s.UpdateUserFn(u)
}

// UnRegister mocks identity.Processes.UnRegister
func (s *IdentityProcesses) UnRegister(u *user.User) error {
	return s.UnRegisterUserFn(u)
}

// AddDomain mockes identity.Processes.AddDomain
func (s *IdentityProcesses) AddDomain(d *domain.Domain) error {
	return s.AddDomainFn(d)
}

// LookUpDomain mockes identity.Processes.LookUpDomain
func (s *IdentityProcesses) LookUpDomain(d *domain.Domain) (*domain.Domain, error) {
	return s.LookUpDomainFn(d)
}

// UpdateDomain mockes identity.Processes.UpdateDomain
func (s *IdentityProcesses) UpdateDomain(d *domain.Domain) error {
	return s.UpdateDomainFn(d)
}

// RemoveDomain mockes identity.Processes.RemoveDomain
func (s *IdentityProcesses) RemoveDomain(d *domain.Domain) error {
	return s.RemoveDomainFn(d)
}
