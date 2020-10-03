package mockadministrator

import (
	"trustdonations.org/m/v2/domain/administrator"
)

//AdministratorService is a mock of administrator.Service
type AdministratorService struct {
	CreateManagementSessionFn func() error
	EndManagementSessionFn    func() error
	CreateFn                  func(a *administrator.Administrator) error
	ReadMultipleFn            func(administratorUUIDs []string) ([]*administrator.Administrator, error)
	ReadFn                    func(a *administrator.Administrator) (*administrator.Administrator, error)
	UpdateFn                  func(a *administrator.Administrator) error
	DeleteFn                  func(a *administrator.Administrator) error
}

// Create mocks administrator.Service.Create
func (s *AdministratorService) Create(a *administrator.Administrator) error {
	return s.CreateFn(a)
}

// Read mocks administrator.Service.Read
func (s *AdministratorService) Read(a *administrator.Administrator) (*administrator.Administrator, error) {
	return s.ReadFn(a)
}

// ReadMultiple mocks administrator.Service.ReadMultiple
func (s *AdministratorService) ReadMultiple(administratorUUIDs []string) ([]*administrator.Administrator, error) {
	return s.ReadMultipleFn(administratorUUIDs)
}

// Update mocks administrator.Service.Update
func (s *AdministratorService) Update(a *administrator.Administrator) error {
	return s.UpdateFn(a)
}

// Delete mocks administrator.Service.Delete
func (s *AdministratorService) Delete(a *administrator.Administrator) error {
	return s.DeleteFn(a)
}

// CreateManagementSession mocks administrator.Service.CreateManagementSession
func (s *AdministratorService) CreateManagementSession() error {
	return s.CreateManagementSessionFn()
}

// EndManagementSession mocks administrator.Service.EndManagementSession
func (s *AdministratorService) EndManagementSession() error {
	return s.EndManagementSessionFn()
}
