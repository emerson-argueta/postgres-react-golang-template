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

var _ Controller = &createGoalController{}

type createGoalController struct {
	Usecase       *usecase.CreateGoalUsecase
	Logger        *log.Logger
	Authorization authorization.JwtService
}

// NewCreateGoalController for createGoal achiever usecase
func NewCreateGoalController(cgtRepos *persistence.Services, logger *log.Logger, authorizationService authorization.JwtService) Controller {
	createGoalUsecase := usecase.NewCreateGoalUsecase(&cgtRepos.Achiever, &cgtRepos.Goal, authorizationService)

	ctrl := &createGoalController{
		Usecase:       createGoalUsecase,
		Logger:        logger,
		Authorization: authorizationService,
	}
	return ctrl
}

// Execute the usecase
func (ctrl *createGoalController) Execute(ctx echo.Context) (e error) {
	var req createGoalRequest

	// Decode the request.
	if err := ctx.Bind(&req); err != nil || req.Name == nil {
		return response.ErrorResponse(ctx.Response().Writer, response.ErrInvalidJSON, http.StatusBadRequest, ctrl.Logger)
	}

	// extract user id from authKey stored by JwtMiddleware handler func
	authKey := ctx.Get("user").(string)
	userID, e := ctrl.Authorization.VerifyTokenAndExtractIDClaim(authKey)
	if e != nil {
		return response.ErrorResponse(ctx.Response().Writer, e, http.StatusInternalServerError, ctrl.Logger)
	}

	dto := &usecase.CreateGoalDTO{AchieverUserID: userID, Name: *req.Name}
	switch e := ctrl.Usecase.Execute(dto); e {
	case nil:
		response.EncodeJSON(ctx.Response().Writer, &createGoalResponse{Message: "Successfully created goal"}, ctrl.Logger)
	case goal.ErrGoalExists:
		return response.ErrorResponse(ctx.Response().Writer, e, http.StatusConflict, ctrl.Logger)
	default:
		return response.ErrorResponse(ctx.Response().Writer, e, http.StatusInternalServerError, ctrl.Logger)
	}

	return nil
}
