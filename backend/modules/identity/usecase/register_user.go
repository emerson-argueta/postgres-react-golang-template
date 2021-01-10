package usecase

import (
	"emersonargueta/m/v1/modules/identity/domain/user"
	"emersonargueta/m/v1/modules/identity/dto"
	"emersonargueta/m/v1/modules/identity/repository"
)

// RegisterUserUsecase performs registering
type RegisterUserUsecase struct {
	UserRepo repository.UserRepo
}

// NewRegisterUserUsecase to register user
func NewRegisterUserUsecase(userRepo repository.UserRepo) *RegisterUserUsecase {
	return &RegisterUserUsecase{UserRepo: userRepo}
}

// Execute using the following business logic
// Verify email and password.
// Check if user exists.
// Create user with hash password.
func (uc *RegisterUserUsecase) Execute(u dto.UserDTO) (e error) {

	email, e := user.NewEmail(u.Email)
	if e != nil {
		return e
	}
	if res, _ := uc.UserRepo.RetrieveUserByEmail(email); res != nil {
		return user.ErrUserExists
	}

	newUser, e := dto.ToDomain(u)
	if e != nil {
		return e
	}

	return uc.UserRepo.CreateUser(newUser)
}
