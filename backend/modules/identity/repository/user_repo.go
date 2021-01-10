package repository

import "emersonargueta/m/v1/modules/identity/domain/user"

// UserRepo used to modify the user model.
type UserRepo interface {
	// CreateUser implementation must return ErrUserExists if the user exists.
	CreateUser(user.User) error
	// RetrieveUserByEmail implementation must return ErrUserNotFound if the user is not found.
	RetrieveUserByEmail(email user.Email) (user.User, error)
	// RetrieveUserID implementation must return ErrUserNotFound if the user is not found.
	RetrieveUserByID(id string) (user.User, error)
	// UpdateUser implementation must search user by uuid and return
	// ErrUserNotFound if user is not found.
	UpdateUser(user.User) error
	// DeleteUser implementation should search the user by id before deleting
	// the user and must return ErrUserNotFound if the user does not exists.
	DeleteUser(id string) error
}
