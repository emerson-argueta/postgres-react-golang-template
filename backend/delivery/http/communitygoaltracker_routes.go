package http

import (
	"emersonargueta/m/v1/authorization/jwt"
	"emersonargueta/m/v1/communitygoaltracker"
	"emersonargueta/m/v1/delivery/middleware"
	"emersonargueta/m/v1/identity"
	"emersonargueta/m/v1/paymentgateway/stripe"
	"log"
	"net/http"
	"os"

	"github.com/labstack/echo"
)

const (
	// AchieverURL used for communitygoaltracker processes that modify an achiever
	AchieverURL = "/communitygoaltracker/achiever"
	// AchieverLoginURL used for communitygoaltracker login process
	AchieverLoginURL = "/communitygoaltracker/achiever/login"
	// GoalURL used for communitygoaltracker processes that modify a goal
	GoalURL = "/communitygoaltracker/goal"
	// GoalAbandonURL used for communitygoaltracker goal abandon process
	GoalAbandonURL = "/communitygoaltracker/goal/abandon"
)

// CommunitygoaltrackerHandler represents an HTTP API handler.
type CommunitygoaltrackerHandler struct {
	*echo.Echo

	communitygoaltracker.Services

	Authorization jwt.Services

	PaymentGateway stripe.Services

	Logger *log.Logger
}

// NewCommunitygoaltrackerHandler returns CommunitygoaltrackerHandler.
func NewCommunitygoaltrackerHandler() *CommunitygoaltrackerHandler {
	h := &CommunitygoaltrackerHandler{
		Echo:   echo.New(),
		Logger: log.New(os.Stderr, "", log.LstdFlags),
	}

	public := h.Group(RoutePrefix)
	public.POST(AchieverURL, h.handleRegister)
	public.POST(AchieverLoginURL, h.handleLogin)
	// TODO: post method to handle authorization for achiever

	restricted := h.Group(RoutePrefix)
	restricted.Use(middleware.JwtMiddleware)
	restricted.PATCH(AchieverURL, h.handleUpdateAchiever)
	restricted.DELETE(AchieverURL, h.handleUnRegister)

	restricted.POST(GoalURL, h.handleCreateGoal)
	restricted.PATCH(GoalURL, h.handleUpdateGoalProgress)
	restricted.DELETE(GoalURL, h.handleAbandonGoal)
	restricted.DELETE(GoalURL, h.handleDeleteGoal)

	return h
}

func (h *CommunitygoaltrackerHandler) handleRegister(ctx echo.Context) (e error) {
	var req achieverRequest

	// Decode the request.
	if err := ctx.Bind(&req); err != nil || req.Achiever == nil {
		return ResponseError(ctx.Response().Writer, ErrInvalidJSON, http.StatusBadRequest, h.Logger)
	}

	switch registeredAchiever, e := h.Communitygoaltracker.Register(req.Achiever); e {
	case nil:
		tokenPair, e := middleware.GenerateTokenPair(*registeredAchiever.UUID, middleware.AccestokenLimit, middleware.RefreshtokenLimit)
		if e != nil {
			return ResponseError(ctx.Response().Writer, e, http.StatusInternalServerError, h.Logger)
		}
		// Do not send password to avoid exploits.
		registeredAchiever.Password = nil
		encodeJSON(ctx.Response().Writer, &achieverResponse{Achiever: registeredAchiever, Token: tokenPair}, h.Logger)
	case communitygoaltracker.ErrIncompleteAchieverDetails:
		return ResponseError(ctx.Response().Writer, e, http.StatusBadRequest, h.Logger)
	default:
		return ResponseError(ctx.Response().Writer, e, http.StatusInternalServerError, h.Logger)
	}

	return nil
}

func (h *CommunitygoaltrackerHandler) handleLogin(ctx echo.Context) error {
	var req achieverRequest

	// Decode the request.
	if err := ctx.Bind(&req); err != nil || req.Achiever == nil || req.Achiever.Email == nil || req.Achiever.Password == nil {
		return ResponseError(ctx.Response().Writer, ErrInvalidJSON, http.StatusBadRequest, h.Logger)
	}

	switch registeredAchiever, e := h.Communitygoaltracker.Login(*req.Achiever.Email, *req.Achiever.Password); e {
	case nil:
		tokenPair, e := middleware.GenerateTokenPair(*registeredAchiever.UUID, middleware.AccestokenLimit, middleware.RefreshtokenLimit)
		if e != nil {
			return ResponseError(ctx.Response().Writer, e, http.StatusInternalServerError, h.Logger)
		}
		// Do not send password to avoid exploits.
		registeredAchiever.Password = nil
		encodeJSON(ctx.Response().Writer, &achieverResponse{Achiever: registeredAchiever, Token: tokenPair}, h.Logger)
	case identity.ErrIncorrectCredentials:
		return ResponseError(ctx.Response().Writer, e, http.StatusUnauthorized, h.Logger)
	default:
		return ResponseError(ctx.Response().Writer, e, http.StatusInternalServerError, h.Logger)
	}

	return nil
}

func (h *CommunitygoaltrackerHandler) handleUpdateAchiever(ctx echo.Context) error {
	return nil
}

func (h *CommunitygoaltrackerHandler) handleUnRegister(ctx echo.Context) error {
	return nil
}

func (h *CommunitygoaltrackerHandler) handleCreateGoal(ctx echo.Context) error {
	return nil
}

func (h *CommunitygoaltrackerHandler) handleUpdateGoalProgress(ctx echo.Context) error {
	return nil
}

func (h *CommunitygoaltrackerHandler) handleAbandonGoal(ctx echo.Context) error {
	return nil
}

func (h *CommunitygoaltrackerHandler) handleDeleteGoal(ctx echo.Context) error {
	return nil
}
