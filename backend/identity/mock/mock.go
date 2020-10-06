package mockuser

import (
	"emersonargueta/m/v1/identity/domain"
	"emersonargueta/m/v1/identity/user"
)

//IdentityService is a mock of identity.Service
type IdentityService struct {
	RegisterFn       func(*user.User) error
	RetrieveFn       func(u *user.User, byEmail bool) (*user.User, error)
	UpdateFn         func(u *user.User, byEmail bool) error
	UnRegisterFn     func(u *user.User, byEmail bool) error
	CreateDomainFn   func(*domain.Domain) error
	RetrieveDomainFn func(*domain.Domain) (*domain.Domain, error)
	UpdateDomainFn   func(*domain.Domain) error
	DeleteDomainFn   func(*domain.Domain) error
}

// Register mocks identity.Service.Register
func (s *IdentityService) Register(u *user.User) error {
	return s.RegisterFn(u)
}

// Retrieve mocks identity.Service.Retrieve
func (s *IdentityService) Retrieve(u *user.User, byEmail bool) (*user.User, error) {
	return s.RetrieveFn(u, byEmail)
}

// Update mocks identity.Service.Update
func (s *IdentityService) Update(u *user.User, byEmail bool) error {
	return s.UpdateFn(u, byEmail)
}

// UnRegister mocks identity.Service.UnRegister
func (s *IdentityService) UnRegister(u *user.User, byEmail bool) error {
	return s.UnRegisterFn(u, byEmail)
}

// CreateDomain mockes identity.Service.CreateDomain
func (s *IdentityService) CreateDomain(d *domain.Domain) error {
	return s.CreateDomainFn(d)
}

// RetrieveDomain mockes identity.Service.RetrieveDomain
func (s *IdentityService) RetrieveDomain(d *domain.Domain) (*domain.Domain, error) {
	return s.RetrieveDomainFn(d)
}

// UpdateDomain mockes identity.Service.UpdateDomain
func (s *IdentityService) UpdateDomain(d *domain.Domain) error {
	return s.UpdateDomainFn(d)
}

// DeleteDomain mockes identity.Service.DeleteDomain
func (s *IdentityService) DeleteDomain(d *domain.Domain) error {
	return s.DeleteDomainFn(d)
}
