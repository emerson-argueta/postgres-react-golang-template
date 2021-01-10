package usecase

import (
	"emersonargueta/m/v1/authorization"
	"emersonargueta/m/v1/modules/communitygoaltracker/repository"
)

// DeleteGoalUsecase performs registering
type DeleteGoalUsecase struct {
	AchieverRepo         repository.AchieverRepo
	GoalRepo             repository.GoalRepo
	AuthorizationService *authorization.Client
}

// NewDeleteGoalUsecase to register user
func NewDeleteGoalUsecase(achieverRepo repository.AchieverRepo, goalRepo repository.GoalRepo, authorizationService *authorization.Client) *DeleteGoalUsecase {
	return &DeleteGoalUsecase{
		AchieverRepo:         achieverRepo,
		GoalRepo:             goalRepo,
		AuthorizationService: authorizationService,
	}
}

// Execute usecase
// DeleteGoal using the following business logic: Retrieve the goal and if the
// goal has no achievers except for the one deleting then delete the goal.
func (uc *DeleteGoalUsecase) Execute(dto *DeleteGoalDTO) (e error) {
	retrievedAchiever, e := uc.AchieverRepo.RetrieveAchieverByUserID(dto.AchieverUserID)
	if e != nil {
		return e
	}
	achieverGoals := retrievedAchiever.GetGoals()
	if e = achieverGoals.Remove(dto.Name); e != nil {
		return e
	}

	retrievedGoal, e := uc.GoalRepo.RetrieveGoalByName(dto.Name)
	if e != nil {
		return e
	}

	if e = retrievedGoal.Delete(dto.AchieverUserID); e != nil {
		return e
	}

	if e = uc.AchieverRepo.UpdateAchiever(retrievedAchiever); e != nil {
		return e
	}

	return uc.GoalRepo.DeleteGoal(retrievedGoal.GetID())
}
