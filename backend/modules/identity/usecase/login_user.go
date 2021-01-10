package usecase

import (
	"emersonargueta/m/v1/authorization"
	"emersonargueta/m/v1/modules/identity/domain/user"
	"emersonargueta/m/v1/modules/identity/repository"
)

// LoginUserUsecase performs registering
type LoginUserUsecase struct {
	UserRepo             repository.UserRepo
	AuthorizationService *authorization.Client
}

// NewLoginUserUsecase to register user
func NewLoginUserUsecase(userRepo repository.UserRepo, authorizationService *authorization.Client) *LoginUserUsecase {
	return &LoginUserUsecase{
		UserRepo:             userRepo,
		AuthorizationService: authorizationService,
	}
}

// Execute LoginUser using the following business logic:
// Check if user exists.
// Check if password matches.
func (uc *LoginUserUsecase) Execute(dto *LoginUserDTO) (res map[string]string, e error) {
	email, e := user.NewEmail(&dto.Email)
	if e != nil {
		return nil, e
	}
	retrievedUser, e := uc.UserRepo.RetrieveUserByEmail(email)
	if e != nil {
		return nil, e
	}
	e = user.CompareHashAndPassword(retrievedUser.GetHashPassword(), &dto.Password)
	if e != nil {
		return nil, e
	}

	return uc.AuthorizationService.JwtService().NewKey(retrievedUser.GetID())
}
