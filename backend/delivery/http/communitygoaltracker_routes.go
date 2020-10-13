package http

import (
	"emersonargueta/m/v1/authorization"
	"emersonargueta/m/v1/authorization/jwt"
	"emersonargueta/m/v1/communitygoaltracker"
	"emersonargueta/m/v1/config"
	"emersonargueta/m/v1/delivery/middleware"
	"emersonargueta/m/v1/identity"
	"log"
	"net/http"
	"os"

	"github.com/labstack/echo"
)

const (
	// CommunitygoalTrackerURLPrefix used for communitygoaltracker routes
	CommunitygoalTrackerURLPrefix = "/communitygoaltracker"
	// AchieverURL used for communitygoaltracker processes that modify an achiever
	AchieverURL = CommunitygoalTrackerURLPrefix + "/achiever"
	// AchieverLoginURL used for communitygoaltracker login process
	AchieverLoginURL = CommunitygoalTrackerURLPrefix + "/achiever/login"
	// GoalURL used for communitygoaltracker processes that modify a goal
	GoalURL = CommunitygoalTrackerURLPrefix + "/goal"
	// GoalAbandonURL used for communitygoaltracker goal abandon process
	GoalAbandonURL = CommunitygoalTrackerURLPrefix + "/goal/abandon"
)

// CommunitygoaltrackerHandler represents an HTTP API handler.
type CommunitygoaltrackerHandler struct {
	*echo.Echo

	Communitygoaltracker communitygoaltracker.AllProcesses

	Authorization authorization.Processes

	// PaymentGateway stripe.Services

	Logger     *log.Logger
	Middleware *middleware.Middleware
}

// NewCommunitygoaltrackerHandler returns CommunitygoaltrackerHandler.
func NewCommunitygoaltrackerHandler(config *config.Config) *CommunitygoaltrackerHandler {
	h := &CommunitygoaltrackerHandler{
		Echo:          echo.New(),
		Logger:        log.New(os.Stderr, "", log.LstdFlags),
		Middleware:    middleware.New(config),
		Authorization: jwt.NewClient(config).Service(),
	}

	public := h.Group(RoutePrefix)
	public.POST(AchieverURL, h.handleRegister)
	public.POST(AchieverLoginURL, h.handleLogin)
	// TODO: post method to handle authorization for achiever

	restricted := h.Group(RoutePrefix)
	restricted.Use(h.Middleware.JwtMiddleware())
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
		tokenPair, e := h.Middleware.GenerateTokenPair(*registeredAchiever.UUID, middleware.AccestokenLimit, middleware.RefreshtokenLimit)
		if e != nil {
			return ResponseError(ctx.Response().Writer, e, http.StatusInternalServerError, h.Logger)
		}
		// Do not send password to avoid exploits.
		registeredAchiever.Password = nil
		encodeJSON(ctx.Response().Writer, &achieverResponse{Achiever: registeredAchiever, Token: tokenPair}, h.Logger)
	case communitygoaltracker.ErrAchieverIncompleteDetails:
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
		tokenPair, e := h.Middleware.GenerateTokenPair(*registeredAchiever.UUID, middleware.AccestokenLimit, middleware.RefreshtokenLimit)
		if e != nil {
			return ResponseError(ctx.Response().Writer, e, http.StatusInternalServerError, h.Logger)
		}
		// Do not send password to avoid exploits.
		registeredAchiever.Password = nil
		encodeJSON(ctx.Response().Writer, &achieverResponse{Achiever: registeredAchiever, Token: tokenPair}, h.Logger)
	case identity.ErrUserIncorrectCredentials:
		return ResponseError(ctx.Response().Writer, e, http.StatusUnauthorized, h.Logger)
	default:
		return ResponseError(ctx.Response().Writer, e, http.StatusInternalServerError, h.Logger)
	}

	return nil
}

func (h *CommunitygoaltrackerHandler) handleUpdateAchiever(ctx echo.Context) error {
	var req achieverRequest

	// Decode the request.
	if err := ctx.Bind(&req); err != nil || req.Achiever == nil {
		return ResponseError(ctx.Response().Writer, ErrInvalidJSON, http.StatusBadRequest, h.Logger)
	}

	// extract administrator uuid from token stored by JwtMiddleware handler func
	token := ctx.Get("user")
	uuid, err := h.Authorization.Authorize(token)
	if err != nil {
		return ResponseError(ctx.Response().Writer, err, http.StatusInternalServerError, h.Logger)
	}

	a := req.Achiever
	a.UUID = uuid

	switch e := h.Communitygoaltracker.UpdateAchiever(a); e {
	case nil:
		// Do not send password to avoid exploits.
		a.Password = nil
		encodeJSON(ctx.Response().Writer, &achieverResponse{Achiever: a}, h.Logger)
	case communitygoaltracker.ErrAchieverIncompleteDetails:
		return ResponseError(ctx.Response().Writer, e, http.StatusBadRequest, h.Logger)
	case communitygoaltracker.ErrAchieverNotFound, identity.ErrUserNotFound:
		return ResponseError(ctx.Response().Writer, e, http.StatusNotFound, h.Logger)
	default:
		return ResponseError(ctx.Response().Writer, e, http.StatusInternalServerError, h.Logger)
	}

	return nil
}

func (h *CommunitygoaltrackerHandler) handleUnRegister(ctx echo.Context) error {
	var req achieverRequest

	// Decode the request.
	if err := ctx.Bind(&req); err != nil || req.Achiever == nil {
		return ResponseError(ctx.Response().Writer, ErrInvalidJSON, http.StatusBadRequest, h.Logger)
	}

	// extract administrator uuid from token stored by JwtMiddleware handler func
	token := ctx.Get("user")
	uuid, err := h.Authorization.Authorize(token)
	if err != nil {
		return ResponseError(ctx.Response().Writer, err, http.StatusInternalServerError, h.Logger)
	}

	a := req.Achiever
	a.UUID = uuid

	switch e := h.Communitygoaltracker.UnRegister(a); e {
	case nil:
		// Do not send password to avoid exploits.
		a.Password = nil
		encodeJSON(ctx.Response().Writer, &achieverResponse{Achiever: a}, h.Logger)
	case communitygoaltracker.ErrAchieverNotFound, communitygoaltracker.ErrGoalNotFound:
		return ResponseError(ctx.Response().Writer, e, http.StatusNotFound, h.Logger)
	default:
		return ResponseError(ctx.Response().Writer, e, http.StatusInternalServerError, h.Logger)
	}

	return nil
}

func (h *CommunitygoaltrackerHandler) handleCreateGoal(ctx echo.Context) error {
	var req goalRequest

	// Decode the request.
	if err := ctx.Bind(&req); err != nil || req.Name == nil {
		return ResponseError(ctx.Response().Writer, ErrInvalidJSON, http.StatusBadRequest, h.Logger)
	}

	// extract administrator uuid from token stored by JwtMiddleware handler func
	token := ctx.Get("user")
	uuid, err := h.Authorization.Authorize(token)
	if err != nil {
		return ResponseError(ctx.Response().Writer, err, http.StatusInternalServerError, h.Logger)
	}

	g := requestToGoal(&req, *uuid)

	switch newGoal, e := h.Communitygoaltracker.CreateGoal(g); e {
	case nil:
		encodeJSON(ctx.Response().Writer, goalToResponse(newGoal), h.Logger)
	case communitygoaltracker.ErrGoalIncompleteDetails:
		return ResponseError(ctx.Response().Writer, e, http.StatusBadRequest, h.Logger)
	case communitygoaltracker.ErrGoalExists:
		return ResponseError(ctx.Response().Writer, e, http.StatusConflict, h.Logger)
	default:
		return ResponseError(ctx.Response().Writer, e, http.StatusInternalServerError, h.Logger)
	}

	return nil
}

func (h *CommunitygoaltrackerHandler) handleUpdateGoalProgress(ctx echo.Context) error {
	var req goalRequest

	// Decode the request.
	if err := ctx.Bind(&req); err != nil || req.ID == nil || req.Progress == nil {
		return ResponseError(ctx.Response().Writer, ErrInvalidJSON, http.StatusBadRequest, h.Logger)
	}

	// extract administrator uuid from token stored by JwtMiddleware handler func
	token := ctx.Get("user")
	uuid, err := h.Authorization.Authorize(token)
	if err != nil {
		return ResponseError(ctx.Response().Writer, err, http.StatusInternalServerError, h.Logger)
	}

	switch updatedGoal, e := h.Communitygoaltracker.UpdateGoalProgress(*uuid, *req.ID, *req.Progress); e {
	case nil:
		encodeJSON(ctx.Response().Writer, goalToResponse(updatedGoal), h.Logger)
	case communitygoaltracker.ErrGoalInvalidProgress, communitygoaltracker.ErrGoalWithNoAchievers, communitygoaltracker.ErrGoalNotFound:
		return ResponseError(ctx.Response().Writer, e, http.StatusBadRequest, h.Logger)
	case communitygoaltracker.ErrGoalExists:
		return ResponseError(ctx.Response().Writer, e, http.StatusConflict, h.Logger)
	default:
		return ResponseError(ctx.Response().Writer, e, http.StatusInternalServerError, h.Logger)
	}

	return nil
}

func (h *CommunitygoaltrackerHandler) handleAbandonGoal(ctx echo.Context) error {
	var req goalRequest

	// Decode the request.
	if err := ctx.Bind(&req); err != nil || req.ID == nil {
		return ResponseError(ctx.Response().Writer, ErrInvalidJSON, http.StatusBadRequest, h.Logger)
	}

	// extract administrator uuid from token stored by JwtMiddleware handler func
	token := ctx.Get("user")
	uuid, err := h.Authorization.Authorize(token)
	if err != nil {
		return ResponseError(ctx.Response().Writer, err, http.StatusInternalServerError, h.Logger)
	}

	switch e := h.Communitygoaltracker.AbandonGoal(*uuid, *req.ID); e {
	case nil:
		encodeJSON(ctx.Response().Writer, req, h.Logger)
	case communitygoaltracker.ErrGoalWithNoAchievers, communitygoaltracker.ErrGoalNotFound:
		return ResponseError(ctx.Response().Writer, e, http.StatusBadRequest, h.Logger)
	case communitygoaltracker.ErrGoalExists:
		return ResponseError(ctx.Response().Writer, e, http.StatusConflict, h.Logger)
	default:
		return ResponseError(ctx.Response().Writer, e, http.StatusInternalServerError, h.Logger)
	}

	return nil
}

func (h *CommunitygoaltrackerHandler) handleDeleteGoal(ctx echo.Context) error {
	var req goalRequest

	// Decode the request.
	if err := ctx.Bind(&req); err != nil || req.Name == nil || req.ID == nil {
		return ResponseError(ctx.Response().Writer, ErrInvalidJSON, http.StatusBadRequest, h.Logger)
	}

	// extract administrator uuid from token stored by JwtMiddleware handler func
	token := ctx.Get("user")
	uuid, err := h.Authorization.Authorize(token)
	if err != nil {
		return ResponseError(ctx.Response().Writer, err, http.StatusInternalServerError, h.Logger)
	}

	switch e := h.Communitygoaltracker.DeleteGoal(*uuid, *req.ID); e {
	case nil:
		encodeJSON(ctx.Response().Writer, req, h.Logger)
	case communitygoaltracker.ErrGoalWithNoAchievers:
		return ResponseError(ctx.Response().Writer, e, http.StatusBadRequest, h.Logger)
	case communitygoaltracker.ErrGoalCannotDelete:
		return ResponseError(ctx.Response().Writer, e, http.StatusUnauthorized, h.Logger)
	case communitygoaltracker.ErrGoalNotFound:
		return ResponseError(ctx.Response().Writer, e, http.StatusNotFound, h.Logger)
	default:
		return ResponseError(ctx.Response().Writer, e, http.StatusInternalServerError, h.Logger)
	}

	return nil
}
