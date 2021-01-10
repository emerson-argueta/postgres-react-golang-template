package goal

// Achiever represents the achiever's status for a goal.
type Achiever struct {
	State    *State
	Progress *Progress
	Messages *Messages
}

// UpdateProgress for an achiever
func (a *Achiever) UpdateProgress(progress Progress) error {
	if a.Progress.ToInt() == 100 {
		return ErrGoalCompleteCannotUpdateProgress
	}
	if progress.ToInt() == 100 {
		state := Completed
		a.State = &state
	}
	a.Progress = &progress
	return nil
}

// AbandonGoal for an achiever. Goal cannot be abaondoned if it is completed.
func (a *Achiever) AbandonGoal() error {
	if a == nil {
		return nil
	}

	if *a.State == Completed {
		return ErrCannotAbandon
	}
	abandonedState := Abondoned
	a.State = &abandonedState
	return nil
}

// AddMessagesToGoal for an achiever
func (a *Achiever) AddMessagesToGoal(message Message, timestamp Timestamp) {
	(*a.Messages)[timestamp] = message
}
