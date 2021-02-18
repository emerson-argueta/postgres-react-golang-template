package usecase

import (
	"emersonargueta/m/v1/modules/communitygoaltracker/domain/goal"
	"emersonargueta/m/v1/modules/communitygoaltracker/repository"
	"emersonargueta/m/v1/shared/infrastructure/http/authorization"
)

// UpdateGoalProgressUsecase performs registering
type UpdateGoalProgressUsecase struct {
	GoalRepo             repository.GoalRepo
	AuthorizationService authorization.JwtService
}

// NewUpdateGoalProgressUsecase to register user
func NewUpdateGoalProgressUsecase(goalRepo repository.GoalRepo, authorizationService authorization.JwtService) *UpdateGoalProgressUsecase {
	return &UpdateGoalProgressUsecase{
		GoalRepo:             goalRepo,
		AuthorizationService: authorizationService,
	}
}

// Execute usecase
// UpdateGoalProgress using the following business logic: Retrieve the goal and
// update the progress for the goal's achiever. Valid progress values are
// between 0 and 100 which can be interpreted as 0 percent to 100 percent. If
// the progress is 100 then update the state for the goal's achiever to
// complete.
func (uc *UpdateGoalProgressUsecase) Execute(dto *UpdateGoalProgressDTO) (e error) {
	retrievedGoal, e := uc.GoalRepo.RetrieveGoalByName(dto.Name)
	if e != nil {
		return e
	}

	goalAchievers := retrievedGoal.GetAchievers()

	achieverForGoal, ok := goalAchievers[dto.AchieverUserID]
	if !ok {
		return goal.ErrGoalNotFound
	}

	progress, e := goal.NewProgress(&dto.Progress)
	if e != nil {
		return e
	}

	if e = achieverForGoal.UpdateProgress(progress); e != nil {
		return e
	}

	return uc.GoalRepo.UpdateGoal(retrievedGoal)
}
