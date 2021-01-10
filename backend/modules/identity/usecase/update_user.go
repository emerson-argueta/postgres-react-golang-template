package usecase

import (
	"emersonargueta/m/v1/authorization"
	"emersonargueta/m/v1/modules/identity/domain/user"
	"emersonargueta/m/v1/modules/identity/dto"
	"emersonargueta/m/v1/modules/identity/repository"
)

// UpdateUserUsecase performs registering
type UpdateUserUsecase struct {
	UserRepo             repository.UserRepo
	AuthorizationService *authorization.Client
}

// NewUpdateUserUsecase to register user
func NewUpdateUserUsecase(userRepo repository.UserRepo, authorizationService *authorization.Client) *UpdateUserUsecase {
	return &UpdateUserUsecase{
		UserRepo:             userRepo,
		AuthorizationService: authorizationService,
	}
}

// Execute update user using the following business logic
// Validate email and password
// Update the user searching by uuid.
// Must return ErrUserExists if
// user's new email to be updated conflicts with another user's email.
func (uc *UpdateUserUsecase) Execute(u dto.UserDTO) (e error) {
	if u.ID == nil {
		return user.ErrUserNotFound
	}
	retreivedUser, e := uc.UserRepo.RetrieveUserByID(*u.ID)
	if e != nil {
		return e
	}

	if u.Email != nil {
		email, e := user.NewEmail(u.Email)
		if e != nil {
			return e
		}
		if userWithEmailExists, _ := uc.UserRepo.RetrieveUserByEmail(email); userWithEmailExists != nil {
			return user.ErrUserExists
		}
		retreivedUser.SetEmail(email)
	}

	if u.Password != nil {
		hashPassword, e := user.NewHashPassword(u.Password)
		if e != nil {
			return e
		}
		retreivedUser.SetHashPassword(hashPassword)
	}

	role, e := user.NewRole(u.Role)
	if e != nil {
		return e
	}
	if u.Role != nil {
		retreivedUser.SetRole(role)
	}

	return uc.UserRepo.UpdateUser(retreivedUser)
}
