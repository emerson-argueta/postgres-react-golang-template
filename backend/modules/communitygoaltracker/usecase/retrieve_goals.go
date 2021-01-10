package usecase

import (
	"emersonargueta/m/v1/authorization"
	"emersonargueta/m/v1/modules/communitygoaltracker/domain/goal"
	"emersonargueta/m/v1/modules/communitygoaltracker/repository"
)

// RetrieveGoalsUsecase performs registering
type RetrieveGoalsUsecase struct {
	AchieverRepo         repository.AchieverRepo
	GoalRepo             repository.GoalRepo
	AuthorizationService *authorization.Client
}

// NewRetrieveGoalsUsecase to register user
func NewRetrieveGoalsUsecase(achieverRepo repository.AchieverRepo, goalRepo repository.GoalRepo, authorizationService *authorization.Client) *RetrieveGoalsUsecase {
	return &RetrieveGoalsUsecase{
		AchieverRepo:         achieverRepo,
		GoalRepo:             goalRepo,
		AuthorizationService: authorizationService,
	}
}

// Execute usecase
// RetrieveGoals using the following business logic
// For an achiever retrieve all of their goals.
func (uc *RetrieveGoalsUsecase) Execute(dto *RetrieveGoalsDTO) (res []goal.Goal, e error) {
	retrievedAchiever, e := uc.AchieverRepo.RetrieveAchieverByUserID(dto.AchieverUserID)
	if e != nil {
		return nil, e
	}
	achieverGoals := retrievedAchiever.GetGoals()
	goalNames := achieverGoals.Names()

	return uc.GoalRepo.RetrieveGoalsByNames(goalNames)
}
