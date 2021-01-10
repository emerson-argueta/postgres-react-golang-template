package mockuser

import (
	"os/user"
)

//IdentityProcesses is a mock of identity.Processes
type IdentityProcesses struct {
	RegisterUserFn   func(*user.User) error
	LoginUserFn      func(email string, password string) (res *user.User, e error)
	UpdateUserFn     func(u *user.User) error
	UnRegisterUserFn func(u *user.User) error
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
