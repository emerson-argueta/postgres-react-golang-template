package goal

// Fields used to create goal
type Fields struct {
	Name      *Name
	Achievers *Achievers
}

// Goal represents a goal that an achiever is trying to complete.
type Goal interface {
	GetID() int64
	GetName() Name
	GetAchievers() Achievers
	Delete(achieverUserID string) error
}

type goal struct {
	ID        int64
	Name      Name
	Achievers Achievers
}

// Create a goal
func Create(goalFields *Fields, id *int64) (res Goal, e error) {
	if goalFields.Name == nil {
		return nil, ErrGoalIncompleteDetails
	}
	goal := &goal{
		Name:      *goalFields.Name,
		Achievers: *goalFields.Achievers,
	}
	isNewGoal := (id == nil)
	if !isNewGoal {
		goal.ID = *id
	}

	return goal, nil
}

func (g *goal) GetID() int64 {
	return g.ID
}
func (g *goal) GetName() Name {
	return g.Name
}
func (g *goal) GetAchievers() Achievers {
	return g.Achievers
}

// Delete goal by achiever given achiever user id
func (g *goal) Delete(achieverUserID string) error {
	userIDs := g.Achievers.UserIDs()
	if len(userIDs) == 1 && userIDs[0] == achieverUserID {
		return nil
	}
	return ErrGoalCannotDelete
}
