package controller

import (
	"emersonargueta/m/v1/modules/communitygoaltracker/domain/goal"
	"emersonargueta/m/v1/modules/communitygoaltracker/infrastructure/persistence"
	"emersonargueta/m/v1/modules/communitygoaltracker/usecase"
	"emersonargueta/m/v1/shared/infrastructure/http/authorization"
	"emersonargueta/m/v1/shared/infrastructure/http/response"
	"log"
	"net/http"

	"github.com/labstack/echo"
)

var _ Controller = &abandonGoalController{}

type abandonGoalController struct {
	Usecase       *usecase.AbandonGoalUsecase
	Logger        *log.Logger
	Authorization authorization.JwtService
}

// NewAbandonGoalController for abandonGoal achiever usecase
func NewAbandonGoalController(cgtRepos *persistence.Services, logger *log.Logger, authorizationService authorization.JwtService) Controller {
	abandonGoalUsecase := usecase.NewAbandonGoalUsecase(&cgtRepos.Achiever, &cgtRepos.Goal, authorizationService)

	ctrl := &abandonGoalController{
		Usecase:       abandonGoalUsecase,
		Logger:        logger,
		Authorization: authorizationService,
	}
	return ctrl
}

// Execute the usecase
func (ctrl *abandonGoalController) Execute(ctx echo.Context) (e error) {
	goalName := ctx.QueryParam("goalName")
	if goalName == "" {
		return response.ErrorResponse(ctx.Response().Writer, goal.ErrGoalNotFound, http.StatusNotFound, ctrl.Logger)
	}

	// extract user id from authKey stored by JwtMiddleware handler func
	authKey := ctx.Get("user").(string)
	userID, err := ctrl.Authorization.VerifyTokenAndExtractIDClaim(authKey)
	if err != nil {
		return response.ErrorResponse(ctx.Response().Writer, err, http.StatusInternalServerError, ctrl.Logger)
	}

	dto := &usecase.AbandonGoalDTO{AchieverUserID: userID, Name: goalName}
	switch e := ctrl.Usecase.Execute(dto); e {
	case nil:
		response.EncodeJSON(
			ctx.Response().Writer,
			&abandonGoalResponse{Message: "Successfully abandoned goal"},
			ctrl.Logger,
		)
	case goal.ErrGoalNotFound:
		return response.ErrorResponse(ctx.Response().Writer, e, http.StatusNotFound, ctrl.Logger)
	case goal.ErrCannotAbandon:
		return response.ErrorResponse(ctx.Response().Writer, e, http.StatusConflict, ctrl.Logger)
	default:
		return response.ErrorResponse(ctx.Response().Writer, e, http.StatusInternalServerError, ctrl.Logger)
	}

	return nil
}
