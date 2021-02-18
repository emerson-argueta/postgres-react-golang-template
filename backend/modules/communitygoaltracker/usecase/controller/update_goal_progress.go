package controller

import (
	"emersonargueta/m/v1/modules/communitygoaltracker/domain/achiever"
	"emersonargueta/m/v1/modules/communitygoaltracker/domain/goal"
	"emersonargueta/m/v1/modules/communitygoaltracker/infrastructure/persistence"
	"emersonargueta/m/v1/modules/communitygoaltracker/usecase"
	"emersonargueta/m/v1/shared/infrastructure/http/authorization"
	"emersonargueta/m/v1/shared/infrastructure/http/response"
	"log"
	"net/http"

	"github.com/labstack/echo"
)

var _ Controller = &updateGoalProgressController{}

type updateGoalProgressController struct {
	Usecase       *usecase.UpdateGoalProgressUsecase
	Logger        *log.Logger
	Authorization authorization.JwtService
}

// NewUpdateGoalProgressController for updateGoalProgress achiever usecase
func NewUpdateGoalProgressController(cgtRepos *persistence.Services, logger *log.Logger, authorizationService authorization.JwtService) Controller {
	updateGoalProgressUsecase := usecase.NewUpdateGoalProgressUsecase(&cgtRepos.Goal, authorizationService)

	ctrl := &updateGoalProgressController{
		Usecase:       updateGoalProgressUsecase,
		Logger:        logger,
		Authorization: authorizationService,
	}
	return ctrl
}

// Execute the usecase
func (ctrl *updateGoalProgressController) Execute(ctx echo.Context) (e error) {
	var req updateGoalProgressRequest

	// Decode the request.
	if err := ctx.Bind(&req); err != nil || req.Name == nil || req.Progress == nil {
		return response.ErrorResponse(ctx.Response().Writer, response.ErrInvalidJSON, http.StatusBadRequest, ctrl.Logger)
	}

	// extract user id from authKey stored by JwtMiddleware handler func
	authKey := ctx.Get("user").(string)
	userID, err := ctrl.Authorization.VerifyTokenAndExtractIDClaim(authKey)
	if err != nil {
		return response.ErrorResponse(ctx.Response().Writer, err, http.StatusInternalServerError, ctrl.Logger)
	}

	dto := &usecase.UpdateGoalProgressDTO{AchieverUserID: userID, Name: *req.Name, Progress: *req.Progress}
	switch e := ctrl.Usecase.Execute(dto); e {
	case nil:
		response.EncodeJSON(
			ctx.Response().Writer,
			&updateGoalProgressResponse{Message: "successfully updated goal progress"},
			ctrl.Logger,
		)
	case achiever.ErrAchieverNotFound:
		return response.ErrorResponse(ctx.Response().Writer, e, http.StatusNotFound, ctrl.Logger)
	case goal.ErrGoalNotFound:
		return response.ErrorResponse(ctx.Response().Writer, e, http.StatusNotFound, ctrl.Logger)
	case goal.ErrGoalCompleteCannotUpdateProgress:
		return response.ErrorResponse(ctx.Response().Writer, e, http.StatusConflict, ctrl.Logger)
	default:
		return response.ErrorResponse(ctx.Response().Writer, e, http.StatusInternalServerError, ctrl.Logger)
	}

	return nil
}
