package goal

import "errors"

// State represents the state of the goal for a particular achiever.
type State int

const (
	// InProgress when a goal is below 100 in Progress.
	InProgress State = iota
	// Abondoned when a goal is no long inprogress.
	Abondoned
	// Completed when a goal is 100 in Progress.
	Completed
)

// Converts State to string representation.
func (s State) String() (res string, e error) {
	if s < InProgress || s > Completed {
		return res, ErrGoalInvalidState
	}

	return [...]string{"inprogress", "abandoned", "completed"}[s], e

}

// ToState Converts string to Role enum
func ToState(s string) (res State, e error) {
	res, ok := map[string]State{"inprogress": InProgress, "abandoned": Abondoned, "completed": Completed}[s]
	if !ok {
		e = errors.New("could convert to State")
	}
	return res, e

}

// NewState creates a new goal state, goal state inprogress by default. Returns
// error if state is invalid.
func NewState(s *string) (res State, e error) {
	if s == nil {
		return InProgress, nil
	}
	return ToState(*s)
}
