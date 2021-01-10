package usecase

import (
	"emersonargueta/m/v1/authorization"
	"emersonargueta/m/v1/modules/communitygoaltracker/domain/achiever"
	"emersonargueta/m/v1/modules/communitygoaltracker/repository"
)

// RetrieveAchieverUsecase performs registering
type RetrieveAchieverUsecase struct {
	AchieverRepo         repository.AchieverRepo
	GoalRepo             repository.GoalRepo
	AuthorizationService *authorization.Client
}

// NewRetrieveAchieverUsecase to register user
func NewRetrieveAchieverUsecase(
	achieverRepo repository.AchieverRepo,
	authorizationService *authorization.Client,
) *RetrieveAchieverUsecase {
	return &RetrieveAchieverUsecase{
		AchieverRepo:         achieverRepo,
		AuthorizationService: authorizationService,
	}
}

// Execute the usecase
func (uc *RetrieveAchieverUsecase) Execute(dto *RetrieveAchieverDTO) (res achiever.Achiever, e error) {
	return uc.AchieverRepo.RetrieveAchieverByUserID(dto.AchieverUserID)
}
