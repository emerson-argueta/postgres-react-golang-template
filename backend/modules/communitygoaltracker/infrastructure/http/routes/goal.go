package routes

import (
	"github.com/labstack/echo"
)

const (
	// GoalURL used for communitygoaltracker processes that modify a goal
	GoalURL = "/goal"
	// GoalAbandonURL used for communitygoaltracker goal abandon process
	GoalAbandonURL = GoalURL + "/abandon"
)

func (h *CommunityGoalTrackerHandler) handleCreateGoal(ctx echo.Context) error {
	return h.Controllers.CreateGoalController.Execute(ctx)
}

func (h *CommunityGoalTrackerHandler) handleRetrieveGoals(ctx echo.Context) error {
	return h.Controllers.RetrieveGoalsController.Execute(ctx)
}
func (h *CommunityGoalTrackerHandler) handleUpdateGoalProgress(ctx echo.Context) error {
	return h.Controllers.UpdateGoalProgressController.Execute(ctx)
}
func (h *CommunityGoalTrackerHandler) handleAbandonGoal(ctx echo.Context) error {
	return h.Controllers.AbandonGoalController.Execute(ctx)
}
func (h *CommunityGoalTrackerHandler) handleDeleteGoal(ctx echo.Context) error {
	return h.Controllers.DeleteGoalController.Execute(ctx)
}
