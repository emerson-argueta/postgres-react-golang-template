package controller

import (
	"emersonargueta/m/v1/authorization"
	"emersonargueta/m/v1/modules/identity/repository"
	"emersonargueta/m/v1/modules/identity/usecase"
	"emersonargueta/m/v1/shared/infrastructure/http/response"
	"log"
	"net/http"

	"github.com/labstack/echo"
)

var _ Controller = &reauthorizeController{}

// ReathorizeController executes reauthorize usecase
type reauthorizeController struct {
	Usecase *usecase.ReauthorizeUsecase
	Logger  *log.Logger
}

// NewReauthorizeController for reauthorize user usecase
func NewReauthorizeController(userRepo repository.UserRepo, logger *log.Logger, authorizationService *authorization.Client) Controller {
	reauthorizeUsecase := usecase.NewReauthorizeUsecase(userRepo, authorizationService)

	ctrl := &reauthorizeController{
		Usecase: reauthorizeUsecase,
		Logger:  logger,
	}
	return ctrl
}

// Execute the usecase
func (ctrl *reauthorizeController) Execute(ctx echo.Context) (e error) {
	var req userRequest

	// Decode the request.
	if err := ctx.Bind(&req); err != nil || req.Authorization == nil {
		return response.ErrorResponse(ctx.Response().Writer, response.ErrInvalidJSON, http.StatusBadRequest, ctrl.Logger)
	}

	switch newKey, e := ctrl.Usecase.Execute(*req.Authorization); e {
	case nil:
		response.EncodeJSON(ctx.Response().Writer, &userResponse{Authorization: &newKey}, ctrl.Logger)
	case authorization.ErrAuthorizationInvalidKey:
		response.ErrorResponse(ctx.Response().Writer, e, http.StatusUnauthorized, ctrl.Logger)
	default:
		return response.ErrorResponse(ctx.Response().Writer, e, http.StatusInternalServerError, ctrl.Logger)
	}

	return nil
}
