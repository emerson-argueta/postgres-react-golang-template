package mockuser

import "emersonargueta/m/v1/identity"

//UserService is a mock of identity.Service
type UserService struct {
	RegisterFn     func(a *identity.User) error
	RetrieveFn     func(a *identity.User, byEmail bool) (*identity.User, error)
	UpdateFn       func(a *identity.User, byEmail bool) error
	UnRegisterFn   func(a *identity.User, byEmail bool) error
	LookUpDomainFn func(a *identity.Domain) (*identity.Domain, error)
}

// Register mocks identity.Service.Register
func (s *UserService) Register(u *identity.User) error {
	return s.RegisterFn(u)
}

// Retrieve mocks identity.Service.Retrieve
func (s *UserService) Retrieve(u *identity.User, byEmail bool) (*identity.User, error) {
	return s.RetrieveFn(u, byEmail)
}

// Update mocks identity.Service.Update
func (s *UserService) Update(u *identity.User, byEmail bool) error {
	return s.UpdateFn(u, byEmail)
}

// UnRegister mocks identity.Service.UnRegister
func (s *UserService) UnRegister(u *identity.User, byEmail bool) error {
	return s.UnRegisterFn(u, byEmail)
}

// LookUpDomain mocks identity.Service.UnRegister
func (s *UserService) LookUpDomain(d *identity.Domain) (*identity.Domain, error) {
	return s.LookUpDomainFn(d)
}
