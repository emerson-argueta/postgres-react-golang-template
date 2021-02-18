package usecase

import (
	"emersonargueta/m/v1/modules/communitygoaltracker/domain/achiever"
	"emersonargueta/m/v1/modules/communitygoaltracker/domain/goal"
	"emersonargueta/m/v1/modules/communitygoaltracker/repository"
	"emersonargueta/m/v1/shared/infrastructure/http/authorization"
)

// CreateGoalUsecase performs registering
type CreateGoalUsecase struct {
	AchieverRepo         repository.AchieverRepo
	GoalRepo             repository.GoalRepo
	AuthorizationService authorization.JwtService
}

// NewCreateGoalUsecase to register user
func NewCreateGoalUsecase(achieverRepo repository.AchieverRepo, goalRepo repository.GoalRepo, authorizationService authorization.JwtService) *CreateGoalUsecase {
	return &CreateGoalUsecase{
		AchieverRepo:         achieverRepo,
		GoalRepo:             goalRepo,
		AuthorizationService: authorizationService,
	}
}

// Execute usecase
// CreateGoal using the following business logic
// Create a goal with achiever as part of goal's achievers.
// Add goal to achiever's goals and update the achiever.
func (uc *CreateGoalUsecase) Execute(dto *CreateGoalDTO) (e error) {
	name, _ := goal.NewName(&dto.Name)
	retrievedGoal, _ := uc.GoalRepo.RetrieveGoalByName(name.ToString())
	if retrievedGoal != nil {
		return goal.ErrGoalExists
	}

	retrievedAchiever, e := uc.AchieverRepo.RetrieveAchieverByUserID(dto.AchieverUserID)
	if e != nil {
		return e
	}
	goals := retrievedAchiever.GetGoals()
	if goals == nil {
		achieverGoals := achiever.NewGoals(nil)
		retrievedAchiever.SetGoals(achieverGoals)
	}

	goals[name.ToString()] = true
	e = uc.AchieverRepo.UpdateAchiever(retrievedAchiever)
	if e != nil {
		return e
	}

	achievers := goal.NewAchievers(dto.AchieverUserID)
	goalFields := &goal.Fields{
		Name:      &name,
		Achievers: &achievers,
	}
	newGoal, e := goal.Create(goalFields, nil)
	if e != nil {
		return e
	}

	return uc.GoalRepo.CreateGoal(newGoal)

}
