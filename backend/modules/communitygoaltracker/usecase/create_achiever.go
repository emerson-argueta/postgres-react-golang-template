package usecase

import (
	"emersonargueta/m/v1/modules/communitygoaltracker/domain/achiever"
	"emersonargueta/m/v1/modules/communitygoaltracker/dto"
	"emersonargueta/m/v1/modules/communitygoaltracker/repository"
)

// CreateAchieverUsecase performs registering
type CreateAchieverUsecase struct {
	AchieverRepo repository.AchieverRepo
}

// NewCreateAchieverUsecase to register user
func NewCreateAchieverUsecase(achieverRepo repository.AchieverRepo) *CreateAchieverUsecase {
	return &CreateAchieverUsecase{AchieverRepo: achieverRepo}
}

// Execute using the following business logic
// Verify email and password.
// Check if user exists.
// Create user with hash password.
func (uc *CreateAchieverUsecase) Execute(caDTO *CreateAchieverDTO) (e error) {

	if res, _ := uc.AchieverRepo.RetrieveAchieverByUserID(caDTO.UserID); res != nil {
		return achiever.ErrAchieverExists
	}

	userID := caDTO.UserID
	newAchiever, e := dto.ToDomain(dto.AchieverDTO{UserID: &userID})
	if e != nil {
		return e
	}

	return uc.AchieverRepo.CreateAchiever(newAchiever)
}
