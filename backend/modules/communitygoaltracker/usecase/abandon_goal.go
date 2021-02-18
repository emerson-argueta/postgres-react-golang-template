package usecase

import (
	"emersonargueta/m/v1/modules/communitygoaltracker/domain/goal"
	"emersonargueta/m/v1/modules/communitygoaltracker/repository"
	"emersonargueta/m/v1/shared/infrastructure/http/authorization"
)

// AbandonGoalUsecase performs registering
type AbandonGoalUsecase struct {
	AchieverRepo         repository.AchieverRepo
	GoalRepo             repository.GoalRepo
	AuthorizationService authorization.JwtService
}

// NewAbandonGoalUsecase to register user
func NewAbandonGoalUsecase(achieverRepo repository.AchieverRepo, goalRepo repository.GoalRepo, authorizationService authorization.JwtService) *AbandonGoalUsecase {
	return &AbandonGoalUsecase{
		AchieverRepo:         achieverRepo,
		GoalRepo:             goalRepo,
		AuthorizationService: authorizationService,
	}
}

// Execute usecase
// AbandonGoal using the following business logic: Retrieve the goal and
// update the state for the goal's achiever.
func (uc *AbandonGoalUsecase) Execute(dto *AbandonGoalDTO) (e error) {
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

	goalAchievers := retrievedGoal.GetAchievers()

	achieverForGoal, ok := goalAchievers[dto.AchieverUserID]
	if !ok {
		return goal.ErrGoalNotFound
	}

	if e = achieverForGoal.AbandonGoal(); e != nil {
		return e
	}

	if e = uc.GoalRepo.UpdateGoal(retrievedGoal); e != nil {
		return e
	}

	return uc.AchieverRepo.UpdateAchiever(retrievedAchiever)
}
