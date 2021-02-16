package routes

import (
	"github.com/labstack/echo"
)

const (
	// UserURL used for identity processes that modify a user
	UserURL = "/user"
	// RegisterURL used for identity login process
	RegisterURL = "/register"
	// LoginURL used for identity login process
	LoginURL = "/login"
	// ReAuthorizeURL used for identity re-authorization process
	ReAuthorizeURL = "/reauthorize"
)

func (h *IdentityHandler) HandleRegisterUser(ctx echo.Context) error {
	return h.Controllers.RegisterUserController.Execute(ctx)
}
func (h *IdentityHandler) handleLoginUser(ctx echo.Context) error {
	return h.Controllers.LoginUserController.Execute(ctx)
}
func (h *IdentityHandler) handleReauthorize(ctx echo.Context) error {
	return h.Controllers.ReauthorizeController.Execute(ctx)
}
func (h *IdentityHandler) handleUpdateUser(ctx echo.Context) error {
	return h.Controllers.UpdateUserController.Execute(ctx)
}
