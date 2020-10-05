package mockchurch

import (
	"emersonargueta/m/v1/domain/church"
)

//ChurchService is a mock of church.Service
type ChurchService struct {
	CreateManagementSessionFn func() error
	EndManagementSessionFn    func() error
	CreateFn                  func(c *church.Church) error
	ReadFn                    func(c *church.Church, byEmail bool) (*church.Church, error)
	ReadMultipleFn            func(churchids []int64) ([]*church.Church, error)
	UpdateFn                  func(c *church.Church, byEmail bool) error
	DeleteFn                  func(c *church.Church, byEmail bool) error
}

// Create mocks church.Service.Create
func (s *ChurchService) Create(c *church.Church) error {
	return s.CreateFn(c)
}

// Read mocks church.Service.Read
func (s *ChurchService) Read(c *church.Church, byEmail bool) (*church.Church, error) {
	return s.ReadFn(c, byEmail)
}

// ReadMultiple mocks church.Service.ReadMultiple
func (s *ChurchService) ReadMultiple(churchids []int64) ([]*church.Church, error) {
	return s.ReadMultipleFn(churchids)
}

// Update mocks church.Service.Update
func (s *ChurchService) Update(c *church.Church, byEmail bool) error {
	return s.UpdateFn(c, byEmail)
}

// Delete mocks church.Service.Delete
func (s *ChurchService) Delete(c *church.Church, byEmail bool) error {
	return s.DeleteFn(c, byEmail)
}

// CreateManagementSession mocks church.Service.CreateManagementSession
func (s *ChurchService) CreateManagementSession() error {
	return s.CreateManagementSessionFn()
}

// EndManagementSession mocks church.Service.EndManagementSession
func (s *ChurchService) EndManagementSession() error {
	return s.EndManagementSessionFn()
}
