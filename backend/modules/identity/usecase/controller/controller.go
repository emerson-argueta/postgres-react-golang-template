package controller

import (
	"emersonargueta/m/v1/authorization"
	"emersonargueta/m/v1/modules/identity/infrastructure/persistence"
	"log"

	"github.com/labstack/echo"
)

// Controller for usescases
type Controller interface {
	Execute(ctx echo.Context) error
}

// Controllers holds all controllers
type Controllers struct {
	RegisterUserController   Controller
	LoginUserController      Controller
	ReauthorizeController    Controller
	UpdateUserController     Controller
	UnRegisterUserController Controller
}

// New controller holds all necessary controllers
func New(
	authorizationService *authorization.Client,
	logger *log.Logger,
) *Controllers {
	controllers := &Controllers{}

	identityRepos := persistence.IdentityRepos

	controllers.RegisterUserController = NewRegisterUserController(&identityRepos.User, logger)
	controllers.LoginUserController = NewLoginUserController(&identityRepos.User, logger, authorizationService)
	controllers.ReauthorizeController = NewReauthorizeController(&identityRepos.User, logger, authorizationService)
	controllers.UpdateUserController = NewUpdateUserController(&identityRepos.User, logger, authorizationService)
	controllers.UnRegisterUserController = nil

	return controllers
}
