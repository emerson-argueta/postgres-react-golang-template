package routes

import (
	"github.com/labstack/echo"
)

const (
	// AchieverURL used for communitygoaltracker processes that modify an achiever
	AchieverURL = "/achiever"
	// AchieversURL used for communitygoaltracker processes that modify an achiever goal achievers
	AchieversURL = AchieverURL + "/achievers"
)

func (h *CommunityGoalTrackerHandler) handleUpdateAchiever(ctx echo.Context) error {
	return h.Controllers.UpdateAchieverController.Execute(ctx)
}

func (h *CommunityGoalTrackerHandler) handleRetrieveAchiever(ctx echo.Context) error {
	return h.Controllers.RetrieveAchieverController.Execute(ctx)
}
func (h *CommunityGoalTrackerHandler) handleRetrieveAchievers(ctx echo.Context) error {
	return h.Controllers.RetrieveAchieversController.Execute(ctx)
}
