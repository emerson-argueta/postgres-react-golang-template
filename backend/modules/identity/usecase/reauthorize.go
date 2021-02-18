package usecase

import (
	"emersonargueta/m/v1/modules/identity/repository"
	"emersonargueta/m/v1/shared/infrastructure/http/authorization"
)

// ReauthorizeUsecase performs registering
type ReauthorizeUsecase struct {
	UserRepo             repository.UserRepo
	AuthorizationService authorization.JwtService
}

// NewReauthorizeUsecase to register user
func NewReauthorizeUsecase(userRepo repository.UserRepo, authorizationService authorization.JwtService) *ReauthorizeUsecase {
	return &ReauthorizeUsecase{
		UserRepo:             userRepo,
		AuthorizationService: authorizationService,
	}
}

// Execute Reauthorize using the following business logic:
// Check if user exists.
// Check if password matches.
func (uc *ReauthorizeUsecase) Execute(dto TokenDTO) (res *authorization.TokenPair, e error) {

	id, err := uc.AuthorizationService.VerifyTokenAndExtractIDClaim(dto.RefreshToken)
	if err != nil {
		return nil, err
	}
	return uc.AuthorizationService.IssueTokenPair(id, nil, nil)
}
