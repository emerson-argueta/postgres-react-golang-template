package mockdonator

import (
	"trustdonations.org/m/v2/domain/donator"
)

//DonatorService is a mock of administrator.Service
type DonatorService struct {
	CreateManagementSessionFn func() error
	EndManagementSessionFn    func() error
	CreateFn                  func(d *donator.Donator) error
	ReadFn                    func(d *donator.Donator, byEmail bool) (*donator.Donator, error)
	ReadMultipleFn            func(donatorids []int64) ([]*donator.Donator, error)
	ReadWithFilterFn          func(d *donator.Donator, filterDonator *donator.Donator) (*donator.Donator, error)
	UpdateFn                  func(d *donator.Donator, byEmail bool) error
	DeleteFn                  func(d *donator.Donator, byEmail bool) error
}

// Create mocks donator.Service.Create
func (s *DonatorService) Create(d *donator.Donator) error {
	return s.CreateFn(d)
}

// Read mocks donator.Service.Read
func (s *DonatorService) Read(d *donator.Donator, byEmail bool) (*donator.Donator, error) {
	return s.ReadFn(d, byEmail)
}

// ReadMultiple mocks donator.Service.ReadMultiple
func (s *DonatorService) ReadMultiple(donatorids []int64) ([]*donator.Donator, error) {
	return s.ReadMultipleFn(donatorids)
}

// ReadWithFilter mocks donator.Service.ReadWithFilter
func (s *DonatorService) ReadWithFilter(d *donator.Donator, filterDonator *donator.Donator) (*donator.Donator, error) {
	return s.ReadWithFilterFn(d, filterDonator)
}

// Update mocks donator.Service.Update
func (s *DonatorService) Update(d *donator.Donator, byEmail bool) error {
	return s.UpdateFn(d, byEmail)
}

// Delete mocks donator.Service.Delete
func (s *DonatorService) Delete(d *donator.Donator, byEmail bool) error {
	return s.DeleteFn(d, byEmail)
}

// CreateManagementSession mocks donator.Service.CreateManagementSession
func (s *DonatorService) CreateManagementSession() error {
	return s.CreateManagementSessionFn()
}

// EndManagementSession mocks donator.Service.EndManagementSession
func (s *DonatorService) EndManagementSession() error {
	return s.EndManagementSessionFn()
}
