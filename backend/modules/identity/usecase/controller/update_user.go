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

var _ Controller = &updateUserController{}

type updateUserController struct {
	Usecase       *usecase.UpdateUserUsecase
	Logger        *log.Logger
	Authorization *authorization.Client
}

// NewUpdateUserController for updateUser user usecase
func NewUpdateUserController(userRepo repository.UserRepo, logger *log.Logger, authorizationService *authorization.Client) Controller {
	updateUserUsecase := usecase.NewUpdateUserUsecase(userRepo, authorizationService)

	ctrl := &updateUserController{
		Usecase:       updateUserUsecase,
		Logger:        logger,
		Authorization: authorizationService,
	}
	return ctrl
}

// Execute the usecase
func (ctrl *updateUserController) Execute(ctx echo.Context) (e error) {
	var req userRequest

	// Decode the request.
	if err := ctx.Bind(&req); err != nil || req.User == nil {
		return response.ErrorResponse(ctx.Response().Writer, response.ErrInvalidJSON, http.StatusBadRequest, ctrl.Logger)
	}

	// extract user id from authKey stored by JwtMiddleware handler func
	authKey := ctx.Get("user")
	id, err := ctrl.Authorization.JwtService().Authorize(authKey)
	if err != nil {
		return response.ErrorResponse(ctx.Response().Writer, err, http.StatusInternalServerError, ctrl.Logger)
	}

	req.User.ID = id
	switch e := ctrl.Usecase.Execute(*req.User); e {
	case nil:
		response.EncodeJSON(ctx.Response().Writer, &updateUserResponse{Message: "successfully updated user"}, ctrl.Logger)
	case user.ErrUserNotFound:
		return response.ErrorResponse(ctx.Response().Writer, e, http.StatusNotFound, ctrl.Logger)
	case user.ErrUserExists:
		return response.ErrorResponse(ctx.Response().Writer, e, http.StatusConflict, ctrl.Logger)
	default:
		return response.ErrorResponse(ctx.Response().Writer, e, http.StatusInternalServerError, ctrl.Logger)
	}

	return nil
}
