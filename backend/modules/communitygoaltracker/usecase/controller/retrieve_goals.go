package controller

import (
	"emersonargueta/m/v1/authorization"
	"emersonargueta/m/v1/modules/communitygoaltracker/domain/achiever"
	"emersonargueta/m/v1/modules/communitygoaltracker/domain/goal"
	"emersonargueta/m/v1/modules/communitygoaltracker/dto"
	"emersonargueta/m/v1/modules/communitygoaltracker/infrastructure/persistence"
	"emersonargueta/m/v1/modules/communitygoaltracker/usecase"
	"emersonargueta/m/v1/shared/infrastructure/http/response"
	"log"
	"net/http"

	"github.com/labstack/echo"
)

var _ Controller = &retrieveGoalsController{}

type retrieveGoalsController struct {
	Usecase       *usecase.RetrieveGoalsUsecase
	Logger        *log.Logger
	Authorization *authorization.Client
}

// NewRetrieveGoalsController for retrieveGoals achiever usecase
func NewRetrieveGoalsController(cgtRepos *persistence.Services, logger *log.Logger, authorizationService *authorization.Client) Controller {
	retrieveGoalsUsecase := usecase.NewRetrieveGoalsUsecase(&cgtRepos.Achiever, &cgtRepos.Goal, authorizationService)

	ctrl := &retrieveGoalsController{
		Usecase:       retrieveGoalsUsecase,
		Logger:        logger,
		Authorization: authorizationService,
	}
	return ctrl
}

// Execute the usecase
func (ctrl *retrieveGoalsController) Execute(ctx echo.Context) (e error) {

	// extract user id from authKey stored by JwtMiddleware handler func
	authKey := ctx.Get("user")
	userID, e := ctrl.Authorization.JwtService().Authorize(authKey)
	if e != nil || userID == nil {
		return response.ErrorResponse(ctx.Response().Writer, e, http.StatusInternalServerError, ctrl.Logger)
	}

	ggDto := &usecase.RetrieveGoalsDTO{AchieverUserID: *userID}
	switch retrievedGoals, e := ctrl.Usecase.Execute(ggDto); e {
	case nil:
		ggRes := make([]*dto.GoalDTO, len(retrievedGoals))
		for i, v := range retrievedGoals {
			gRes := dto.GoalToDTO(v)
			ggRes[i] = gRes
		}

		response.EncodeJSON(ctx.Response().Writer, &retrieveGoalsResponse{Goals: &ggRes}, ctrl.Logger)
	case achiever.ErrAchieverNotFound:
		return response.ErrorResponse(ctx.Response().Writer, e, http.StatusNotFound, ctrl.Logger)
	case goal.ErrGoalNotFound:
		return response.ErrorResponse(ctx.Response().Writer, e, http.StatusNotFound, ctrl.Logger)
	default:
		return response.ErrorResponse(ctx.Response().Writer, e, http.StatusInternalServerError, ctrl.Logger)
	}

	return nil
}
