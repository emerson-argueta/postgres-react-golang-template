package controller

import (
	"emersonargueta/m/v1/authorization"
	"emersonargueta/m/v1/modules/communitygoaltracker/domain/goal"
	"emersonargueta/m/v1/modules/communitygoaltracker/infrastructure/persistence"
	"emersonargueta/m/v1/modules/communitygoaltracker/usecase"
	"emersonargueta/m/v1/shared/infrastructure/http/response"
	"log"
	"net/http"

	"github.com/labstack/echo"
)

var _ Controller = &deleteGoalController{}

type deleteGoalController struct {
	Usecase       *usecase.DeleteGoalUsecase
	Logger        *log.Logger
	Authorization *authorization.Client
}

// NewDeleteGoalController for deleteGoal achiever usecase
func NewDeleteGoalController(cgtRepos *persistence.Services, logger *log.Logger, authorizationService *authorization.Client) Controller {
	deleteGoalUsecase := usecase.NewDeleteGoalUsecase(&cgtRepos.Achiever, &cgtRepos.Goal, authorizationService)

	ctrl := &deleteGoalController{
		Usecase:       deleteGoalUsecase,
		Logger:        logger,
		Authorization: authorizationService,
	}
	return ctrl
}

// Execute the usecase
func (ctrl *deleteGoalController) Execute(ctx echo.Context) (e error) {
	goalName := ctx.QueryParam("goalName")
	if goalName == "" {
		return response.ErrorResponse(ctx.Response().Writer, goal.ErrGoalNotFound, http.StatusNotFound, ctrl.Logger)
	}

	// extract user id from authKey stored by JwtMiddleware handler func
	authKey := ctx.Get("user")
	userID, err := ctrl.Authorization.JwtService().Authorize(authKey)
	if err != nil {
		return response.ErrorResponse(ctx.Response().Writer, err, http.StatusInternalServerError, ctrl.Logger)
	}

	dto := &usecase.DeleteGoalDTO{AchieverUserID: *userID, Name: goalName}
	switch e := ctrl.Usecase.Execute(dto); e {
	case nil:
		response.EncodeJSON(
			ctx.Response().Writer,
			&deletedGoalResponse{Message: "successfully deleted goal"},
			ctrl.Logger,
		)
	case goal.ErrGoalNotFound:
		return response.ErrorResponse(ctx.Response().Writer, e, http.StatusNotFound, ctrl.Logger)
	case goal.ErrGoalCannotDelete:
		return response.ErrorResponse(ctx.Response().Writer, e, http.StatusConflict, ctrl.Logger)
	default:
		return response.ErrorResponse(ctx.Response().Writer, e, http.StatusInternalServerError, ctrl.Logger)
	}

	return nil
}
