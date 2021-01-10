package usecase

import (
	"emersonargueta/m/v1/authorization"
	"emersonargueta/m/v1/modules/identity/repository"
)

// ReauthorizeUsecase performs registering
type ReauthorizeUsecase struct {
	UserRepo             repository.UserRepo
	AuthorizationService *authorization.Client
}

// NewReauthorizeUsecase to register user
func NewReauthorizeUsecase(userRepo repository.UserRepo, authorizationService *authorization.Client) *ReauthorizeUsecase {
	return &ReauthorizeUsecase{
		UserRepo:             userRepo,
		AuthorizationService: authorizationService,
	}
}

// Execute Reauthorize using the following business logic:
// Check if user exists.
// Check if password matches.
func (uc *ReauthorizeUsecase) Execute(dto ReauthorizeDTO) (res map[string]string, e error) {
	authorizeKey := make(map[string]string)
	for k, v := range dto {
		authorizeKey[k] = v
	}
	return uc.AuthorizationService.JwtService().ReAuthorize(authorizeKey)
}
