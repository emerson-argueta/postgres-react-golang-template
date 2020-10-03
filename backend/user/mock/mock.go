package mockuser

import "trustdonations.org/m/v2/user"

//UserService is a mock of user.Service
type UserService struct {
	CreateManagementSessionFn func() error
	EndManagementSessionFn    func() error
	CreateFn                  func(a *user.User) error
	ReadFn                    func(a *user.User, byEmail bool) (*user.User, error)
	UpdateFn                  func(a *user.User, byEmail bool) error
	DeleteFn                  func(a *user.User, byEmail bool) error
}

// Create mocks user.Service.Create
func (s *UserService) Create(u *user.User) error {
	return s.CreateFn(u)
}

// Read mocks user.Service.Read
func (s *UserService) Read(u *user.User, byEmail bool) (*user.User, error) {
	return s.ReadFn(u, byEmail)
}

// Update mocks user.Service.Update
func (s *UserService) Update(u *user.User, byEmail bool) error {
	return s.UpdateFn(u, byEmail)
}

// Delete mocks user.Service.Delete
func (s *UserService) Delete(u *user.User, byEmail bool) error {
	return s.DeleteFn(u, byEmail)
}

// CreateManagementSession mocks user.Service.CreateManagementSession
func (s *UserService) CreateManagementSession() error {
	return s.CreateManagementSessionFn()
}

// EndManagementSession mocks user.Service.EndManagementSession
func (s *UserService) EndManagementSession() error {
	return s.EndManagementSessionFn()
}
