package controller

import (
	"emersonargueta/m/v1/authorization"
	"emersonargueta/m/v1/modules/communitygoaltracker/domain/achiever"
	"emersonargueta/m/v1/modules/communitygoaltracker/domain/goal"
	"emersonargueta/m/v1/modules/communitygoaltracker/infrastructure/persistence"
	"emersonargueta/m/v1/modules/communitygoaltracker/usecase"
	"emersonargueta/m/v1/shared/infrastructure/http/response"
	"log"
	"net/http"

	"github.com/labstack/echo"
)

var _ Controller = &retrieveAchieversController{}

type retrieveAchieversController struct {
	Usecase       *usecase.RetrieveAchieversUsecase
	Logger        *log.Logger
	Authorization *authorization.Client
}

// NewRetrieveAchieversController for retrieveAchievers achiever usecase
func NewRetrieveAchieversController(cgtRepos *persistence.Services, logger *log.Logger, authorizationService *authorization.Client) Controller {
	retrieveAchieversUsecase := usecase.NewRetrieveAchieversUsecase(&cgtRepos.Achiever, &cgtRepos.Goal, authorizationService)

	ctrl := &retrieveAchieversController{
		Usecase:       retrieveAchieversUsecase,
		Logger:        logger,
		Authorization: authorizationService,
	}
	return ctrl
}

// Execute the usecase
func (ctrl *retrieveAchieversController) Execute(ctx echo.Context) (e error) {
	goalName := ctx.QueryParam("goalName")
	if goalName == "" {
		return response.ErrorResponse(ctx.Response().Writer, goal.ErrGoalNotFound, http.StatusNotFound, ctrl.Logger)
	}

	// extract user id from authKey stored by JwtMiddleware handler func
	authKey := ctx.Get("user")
	userID, e := ctrl.Authorization.JwtService().Authorize(authKey)
	if e != nil || userID == nil {
		return response.ErrorResponse(ctx.Response().Writer, e, http.StatusInternalServerError, ctrl.Logger)
	}

	dto := &usecase.RetrieveAchieversDTO{AchieverUserID: *userID, AchieverGoalName: goalName}
	switch achievers, e := ctrl.Usecase.Execute(dto); e {
	case nil:
		achieversRes := make(achieversDTO)
		for _, a := range achievers {
			firstname := a.GetFirstname().ToString()
			lastname := a.GetLastname().ToString()
			userID := a.GetUserID().ToString()
			aDTO := &achieverDTO{
				Firstname: &firstname,
				Lastname:  &lastname,
				UserID:    &userID,
			}
			achieversRes[a.GetUserID().ToString()] = aDTO
		}
		response.EncodeJSON(ctx.Response().Writer, &retrieveAchieversResponse{Achievers: &achieversRes}, ctrl.Logger)
	case achiever.ErrAchieverNotFound:
		return response.ErrorResponse(ctx.Response().Writer, e, http.StatusNotFound, ctrl.Logger)
	case goal.ErrGoalNotFound:
		return response.ErrorResponse(ctx.Response().Writer, e, http.StatusNotFound, ctrl.Logger)
	default:
		return response.ErrorResponse(ctx.Response().Writer, e, http.StatusInternalServerError, ctrl.Logger)
	}

	return nil
}
