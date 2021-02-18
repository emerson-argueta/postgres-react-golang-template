package usecase

import (
	"emersonargueta/m/v1/modules/communitygoaltracker/domain/achiever"
	"emersonargueta/m/v1/modules/communitygoaltracker/domain/goal"
	"emersonargueta/m/v1/modules/communitygoaltracker/repository"
	"emersonargueta/m/v1/shared/infrastructure/http/authorization"
)

// RetrieveAchieversUsecase performs registering
type RetrieveAchieversUsecase struct {
	AchieverRepo         repository.AchieverRepo
	GoalRepo             repository.GoalRepo
	AuthorizationService authorization.JwtService
}

// NewRetrieveAchieversUsecase to register user
func NewRetrieveAchieversUsecase(
	achieverRepo repository.AchieverRepo,
	goalRepo repository.GoalRepo,
	authorizationService authorization.JwtService,
) *RetrieveAchieversUsecase {
	return &RetrieveAchieversUsecase{
		AchieverRepo:         achieverRepo,
		GoalRepo:             goalRepo,
		AuthorizationService: authorizationService,
	}
}

// Execute the usecase
// RetrieveAchievers using the following business logic
// For an achiever's goal retrieve all of its achievers. Do not send sensitive information
// like address, phone number, role, address, phone, email, and password.
func (uc *RetrieveAchieversUsecase) Execute(dto *RetrieveAchieversDTO) (res []achiever.Achiever, e error) {
	if a, e := uc.AchieverRepo.RetrieveAchieverByUserID(dto.AchieverUserID); e != nil {
		return nil, e
	} else if a.GetGoals() == nil {
		return nil, goal.ErrGoalNotFound
	} else if _, ok := a.GetGoals()[dto.AchieverGoalName]; !ok {
		return nil, goal.ErrGoalNotFound
	}

	g, e := uc.GoalRepo.RetrieveGoalByName(dto.AchieverGoalName)
	if e != nil {
		return nil, e
	}
	achievers := g.GetAchievers()
	achieverUserIDs := achievers.UserIDs()

	return uc.AchieverRepo.RetrieveAchieversByUserIDs(achieverUserIDs)
}
