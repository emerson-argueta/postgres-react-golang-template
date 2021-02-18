package usecase

import (
	"emersonargueta/m/v1/modules/communitygoaltracker/domain/achiever"
	"emersonargueta/m/v1/modules/communitygoaltracker/dto"
	"emersonargueta/m/v1/modules/communitygoaltracker/repository"
	"emersonargueta/m/v1/shared/infrastructure/http/authorization"
)

// UpdateAchieverUsecase performs registering
type UpdateAchieverUsecase struct {
	AchieverRepo         repository.AchieverRepo
	AuthorizationService authorization.JwtService
}

// NewUpdateAchieverUsecase to register user
func NewUpdateAchieverUsecase(achieverRepo repository.AchieverRepo, authorizationService authorization.JwtService) *UpdateAchieverUsecase {
	return &UpdateAchieverUsecase{
		AchieverRepo:         achieverRepo,
		AuthorizationService: authorizationService,
	}
}

// Execute usecase
func (uc *UpdateAchieverUsecase) Execute(aDTO dto.AchieverDTO) (e error) {
	if aDTO.UserID == nil {
		return achiever.ErrAchieverNotFound
	}
	retreivedAchiever, e := uc.AchieverRepo.RetrieveAchieverByUserID(*aDTO.UserID)
	if e != nil {
		return e
	}

	if aDTO.Firstname != nil {
		firstname := achiever.NewName(aDTO.Firstname)
		retreivedAchiever.SetFirstname(firstname)
	}

	if aDTO.Lastname != nil {
		lastname := achiever.NewName(aDTO.Lastname)
		retreivedAchiever.SetLastname(lastname)
	}

	if aDTO.Address != nil {
		address := achiever.NewAddress(aDTO.Address)
		retreivedAchiever.SetAddress(address)
	}
	if aDTO.Phone != nil {
		phone, e := achiever.NewPhone(aDTO.Phone)
		if e != nil {
			return e
		}
		retreivedAchiever.SetPhone(phone)
	}

	return uc.AchieverRepo.UpdateAchiever(retreivedAchiever)
}
