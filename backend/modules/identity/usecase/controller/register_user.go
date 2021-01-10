package controller

import (
	"emersonargueta/m/v1/modules/identity/domain/user"
	"emersonargueta/m/v1/modules/identity/repository"
	"emersonargueta/m/v1/modules/identity/usecase"
	"emersonargueta/m/v1/shared/infrastructure/http/response"
	"log"
	"net/http"

	"github.com/labstack/echo"
)

var _ Controller = &registerUserController{}

// RegisterUserController executes register usecase
type registerUserController struct {
	Usecase *usecase.RegisterUserUsecase
	Logger  *log.Logger
}

// NewRegisterUserController for register user usecase
func NewRegisterUserController(userRepo repository.UserRepo, logger *log.Logger) Controller {
	registerUsecase := usecase.NewRegisterUserUsecase(userRepo)

	ctrl := &registerUserController{
		Usecase: registerUsecase,
		Logger:  logger,
	}
	return ctrl
}

// Execute the usecase
func (ctrl *registerUserController) Execute(ctx echo.Context) (e error) {
	var req userRequest

	// Decode the request.
	if err := ctx.Bind(&req); err != nil || req.User == nil {
		return response.ErrorResponse(ctx.Response().Writer, response.ErrInvalidJSON, http.StatusBadRequest, ctrl.Logger)
	}

	switch e := ctrl.Usecase.Execute(*req.User); e {
	case nil:
		response.EncodeJSON(ctx.Response().Writer, &registerResponse{Message: "registration successful"}, ctrl.Logger)
	case user.ErrUserExists:
		return response.ErrorResponse(ctx.Response().Writer, e, http.StatusConflict, ctrl.Logger)
	case user.ErrUserIncompleteDetails:
		return response.ErrorResponse(ctx.Response().Writer, e, http.StatusBadRequest, ctrl.Logger)
	default:
		return response.ErrorResponse(ctx.Response().Writer, e, http.StatusInternalServerError, ctrl.Logger)
	}

	return nil
}
