package routes

import (
	"emersonargueta/m/v1/authorization"
	"emersonargueta/m/v1/modules/communitygoaltracker/usecase/controller"
	"emersonargueta/m/v1/shared/infrastructure/http/middleware"
	"log"
	"os"

	"github.com/labstack/echo"
)

const (
	// CommunitygoalTrackerURLPrefix used for communitygoaltracker routes
	CommunitygoalTrackerURLPrefix = "/communitygoaltracker"
)

// CommunityGoalTrackerHandler represents an HTTP API handler.
type CommunityGoalTrackerHandler struct {
	*echo.Echo
	*controller.Controllers
	Logger *log.Logger
}

// NewCommunityGoalTrackerHandler uses the labstack echo router.
func NewCommunityGoalTrackerHandler(apiBaseURL string) *CommunityGoalTrackerHandler {
	h := new(CommunityGoalTrackerHandler)

	echoRouter := echo.New()
	logger := log.New(os.Stderr, "", log.LstdFlags)

	authorizationService := authorization.AuthorizationService
	controllers := controller.New(authorizationService, logger)

	h.Echo = echoRouter
	h.Logger = logger
	h.Controllers = controllers

	restricted := h.Group(apiBaseURL + CommunitygoalTrackerURLPrefix)
	restricted.Use(middleware.JwtMiddleware)
	restricted.GET(AchieverURL, h.handleRetrieveAchiever)
	restricted.GET(AchieversURL, h.handleRetrieveAchievers)
	restricted.PATCH(AchieverURL, h.handleUpdateAchiever)

	restricted.Use(middleware.JwtMiddleware)
	restricted.POST(GoalURL, h.handleCreateGoal)
	restricted.GET(GoalURL, h.handleRetrieveGoals)
	restricted.PATCH(GoalURL, h.handleUpdateGoalProgress)
	restricted.DELETE(GoalAbandonURL, h.handleAbandonGoal)
	restricted.DELETE(GoalURL, h.handleDeleteGoal)

	return h
}
