package repository

import "emersonargueta/m/v1/modules/communitygoaltracker/domain/goal"

// GoalRepo used to modify the goal model.
type GoalRepo interface {
	// CreateGoal implementation must return ErrGoalExists if the goal
	// exists.
	CreateGoal(goal.Goal) error
	// RetrieveGoalByID implementation must return ErrGoalNotFound if the
	// goal is not found.
	RetrieveGoalByID(id int64) (goal.Goal, error)
	// RetrieveGoalByName implementation must return ErrGoalNotFound if the
	// goal is not found.
	RetrieveGoalByName(name string) (goal.Goal, error)
	// UpdateDomain implementation must search goal by id and return
	// ErrGoalNotFound if goal not found. Must return ErrGoalExists if
	// the update name conflicts with another goal.
	UpdateGoal(goal.Goal) error
	// DeleteGoal implementation must search goal by id and return
	// ErrGoalNotFound if goal not found.
	DeleteGoal(id int64) error
	// RetrieveGoalsByIDs implementation must return ErrGoalNotFound if the
	// none of the goals are found.
	RetrieveGoalsByIDs(ids []int64) ([]goal.Goal, error)
	// RetrieveGoalsByNames implementation must return ErrGoalNotFound if the
	// none of the goals are found.
	RetrieveGoalsByNames(names []string) ([]goal.Goal, error)
	// UpdateDomain implementation must search goals by id and return
	// ErrGoalNotFound if none of the goals not found. Must return ErrGoalExists if
	// any of the update names conflicts with another goal.
	UpdateGoals(goals []goal.Goal) error
}
