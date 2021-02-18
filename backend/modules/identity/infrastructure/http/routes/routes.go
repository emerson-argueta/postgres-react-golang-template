package routes

import (
	"emersonargueta/m/v1/modules/identity/usecase/controller"
	"emersonargueta/m/v1/shared/infrastructure"
	"emersonargueta/m/v1/shared/infrastructure/http/authorization"
	"emersonargueta/m/v1/shared/infrastructure/http/middleware"
	"log"
	"os"

	"github.com/labstack/echo"
)

const (
	// IdentityURLPrefix used for identity routes
	IdentityURLPrefix = "/identity"
)

// IdentityHandler represents an HTTP API handler.
type IdentityHandler struct {
	*echo.Echo
	*controller.Controllers
	Logger *log.Logger
}

// NewIdentityHandler uses the labstack echo router.
func NewIdentityHandler(apiBaseURL string) *IdentityHandler {
	h := new(IdentityHandler)

	echoRouter := echo.New()
	logger := log.New(os.Stderr, "", log.LstdFlags)

	authorizationService := authorization.NewJWTService(infrastructure.GlobalConfig)
	controllers := controller.New(authorizationService, logger)

	h.Echo = echoRouter
	h.Logger = logger
	h.Controllers = controllers

	public := h.Group(apiBaseURL + IdentityURLPrefix)
	public.POST(RegisterURL, h.HandleRegisterUser)
	public.POST(LoginURL, h.handleLoginUser)
	public.POST(ReAuthorizeURL, h.handleReauthorize)

	restricted := h.Group(apiBaseURL + IdentityURLPrefix)
	restricted.Use(middleware.JwtMiddleware)
	restricted.PATCH(UserURL, h.handleUpdateUser)

	return h
}
