package controller

import (
	"emersonargueta/m/v1/authorization"
	"emersonargueta/m/v1/modules/identity/domain/user"
	"emersonargueta/m/v1/modules/identity/repository"
	"emersonargueta/m/v1/modules/identity/usecase"
	"emersonargueta/m/v1/shared/infrastructure/http/response"
	"log"
	"net/http"

	"github.com/labstack/echo"
)

var _ Controller = &loginUserController{}

// LoginUserController executes login usecase
type loginUserController struct {
	Usecase *usecase.LoginUserUsecase
	Logger  *log.Logger
}

// NewLoginUserController for login user usecase
func NewLoginUserController(userRepo repository.UserRepo, logger *log.Logger, authorizationService *authorization.Client) Controller {
	loginUsecase := usecase.NewLoginUserUsecase(userRepo, authorizationService)

	ctrl := &loginUserController{
		Usecase: loginUsecase,
		Logger:  logger,
	}
	return ctrl
}

// Execute the usecase
func (ctrl *loginUserController) Execute(ctx echo.Context) (e error) {
	var req loginRequest

	// Decode the request.
	if err := ctx.Bind(&req); err != nil || req.Email == nil || req.Password == nil {
		return response.ErrorResponse(ctx.Response().Writer, response.ErrInvalidJSON, http.StatusBadRequest, ctrl.Logger)
	}

	dto := &usecase.LoginUserDTO{
		Email:    *req.Email,
		Password: *req.Password,
	}
	switch newKey, e := ctrl.Usecase.Execute(dto); e {
	case nil:
		response.EncodeJSON(ctx.Response().Writer, &loginResponse{Authorization: &newKey}, ctrl.Logger)
	case user.ErrUserIncorrectCredentials:
		return response.ErrorResponse(ctx.Response().Writer, e, http.StatusUnauthorized, ctrl.Logger)
	case user.ErrUserNotFound:
		return response.ErrorResponse(ctx.Response().Writer, e, http.StatusNotFound, ctrl.Logger)
	default:
		return response.ErrorResponse(ctx.Response().Writer, e, http.StatusInternalServerError, ctrl.Logger)
	}

	return nil
}
