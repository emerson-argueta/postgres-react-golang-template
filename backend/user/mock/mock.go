package mockuser

import "emersonargueta/m/v1/user"

//UserService is a mock of user.Service
type UserService struct {
	RegisterFn   func(a *user.User) error
	RetrieveFn   func(a *user.User, byEmail bool) (*user.User, error)
	UpdateFn     func(a *user.User, byEmail bool) error
	UnRegisterFn func(a *user.User, byEmail bool) error
}

// Register mocks user.Service.Register
func (s *UserService) Register(u *user.User) error {
	return s.RegisterFn(u)
}

// Retrieve mocks user.Service.Retrieve
func (s *UserService) Retrieve(u *user.User, byEmail bool) (*user.User, error) {
	return s.RetrieveFn(u, byEmail)
}

// Update mocks user.Service.Update
func (s *UserService) Update(u *user.User, byEmail bool) error {
	return s.UpdateFn(u, byEmail)
}

// UnRegister mocks user.Service.UnRegister
func (s *UserService) UnRegister(u *user.User, byEmail bool) error {
	return s.UnRegisterFn(u, byEmail)
}
