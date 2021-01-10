package controller

import (
	"emersonargueta/m/v1/authorization"
	"emersonargueta/m/v1/modules/communitygoaltracker/infrastructure/persistence"
	"log"

	"github.com/labstack/echo"
)

// Controller for usescases
type Controller interface {
	Execute(ctx echo.Context) error
}

// Controllers holds all controllers
type Controllers struct {
	RetrieveAchieverController   Controller
	RetrieveAchieversController  Controller
	UpdateAchieverController     Controller
	CreateGoalController         Controller
	RetrieveGoalsController      Controller
	UpdateGoalProgressController Controller
	AbandonGoalController        Controller
	DeleteGoalController         Controller
}

// New controller holds all necessary controllers
func New(
	authorizationService *authorization.Client,
	logger *log.Logger,
) *Controllers {
	controllers := &Controllers{}

	cgtRepos := persistence.CommunitygoaltrackerRepos

	// TODO
	controllers.RetrieveAchieverController = NewRetrieveAchieverController(cgtRepos, logger, authorizationService)
	controllers.RetrieveAchieversController = NewRetrieveAchieversController(cgtRepos, logger, authorizationService)
	controllers.UpdateAchieverController = NewUpdateAchieverController(cgtRepos, logger, authorizationService)
	controllers.CreateGoalController = NewCreateGoalController(cgtRepos, logger, authorizationService)
	controllers.RetrieveGoalsController = NewRetrieveGoalsController(cgtRepos, logger, authorizationService)
	controllers.UpdateGoalProgressController = NewUpdateGoalProgressController(cgtRepos, logger, authorizationService)
	controllers.AbandonGoalController = NewAbandonGoalController(cgtRepos, logger, authorizationService)
	controllers.DeleteGoalController = NewDeleteGoalController(cgtRepos, logger, authorizationService)

	return controllers
}
