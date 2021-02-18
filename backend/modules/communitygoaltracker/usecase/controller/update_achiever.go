package controller

import (
	"emersonargueta/m/v1/modules/communitygoaltracker/domain/achiever"
	"emersonargueta/m/v1/modules/communitygoaltracker/infrastructure/persistence"
	"emersonargueta/m/v1/modules/communitygoaltracker/usecase"
	"emersonargueta/m/v1/shared/infrastructure/http/authorization"
	"emersonargueta/m/v1/shared/infrastructure/http/response"
	"log"
	"net/http"

	"github.com/labstack/echo"
)

var _ Controller = &updateAchieverController{}

type updateAchieverController struct {
	Usecase       *usecase.UpdateAchieverUsecase
	Logger        *log.Logger
	Authorization authorization.JwtService
}

// NewUpdateAchieverController for updateAchiever achiever usecase
func NewUpdateAchieverController(cgtRepos *persistence.Services, logger *log.Logger, authorizationService authorization.JwtService) Controller {
	updateAchieverUsecase := usecase.NewUpdateAchieverUsecase(&cgtRepos.Achiever, authorizationService)

	ctrl := &updateAchieverController{
		Usecase:       updateAchieverUsecase,
		Logger:        logger,
		Authorization: authorizationService,
	}
	return ctrl
}

// Execute the usecase
func (ctrl *updateAchieverController) Execute(ctx echo.Context) (e error) {
	var req achieverRequest

	// Decode the request.
	if err := ctx.Bind(&req); err != nil || req.Achiever == nil {
		return response.ErrorResponse(ctx.Response().Writer, response.ErrInvalidJSON, http.StatusBadRequest, ctrl.Logger)
	}

	// extract user id from authKey stored by JwtMiddleware handler func
	authKey := ctx.Get("user").(string)
	userID, err := ctrl.Authorization.VerifyTokenAndExtractIDClaim(authKey)
	if err != nil {
		return response.ErrorResponse(ctx.Response().Writer, err, http.StatusInternalServerError, ctrl.Logger)
	}

	req.Achiever.UserID = &userID
	switch e := ctrl.Usecase.Execute(*req.Achiever); e {
	case nil:
		response.EncodeJSON(ctx.Response().Writer, &updateAchieverResponse{Message: "successfully updated achiever"}, ctrl.Logger)
	case achiever.ErrAchieverNotFound:
		return response.ErrorResponse(ctx.Response().Writer, e, http.StatusNotFound, ctrl.Logger)
	case achiever.ErrAchieverExists:
		return response.ErrorResponse(ctx.Response().Writer, e, http.StatusConflict, ctrl.Logger)
	default:
		return response.ErrorResponse(ctx.Response().Writer, e, http.StatusInternalServerError, ctrl.Logger)
	}

	return nil
}
