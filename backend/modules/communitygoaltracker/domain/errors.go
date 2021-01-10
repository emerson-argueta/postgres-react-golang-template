package communitygoaltracker

// Goal model errors.
const (
	ErrGoalExists            = Error("goal exists")
	ErrGoalNotFound          = Error("goal not found")
	ErrGoalIncompleteDetails = Error("incomplete details for goal")
	ErrGoalInvalidProgress   = Error("invalid progress value for goal")
	ErrGoalWithNoAchievers   = Error("goal does not have achievers")
	ErrGoalCannotDelete      = Error("cannot delete goal which contains more than 1 achiever")
)

// Achiever moodel errors.
const (
	ErrAchieverExists            = Error("achiever exists")
	ErrAchieverNotFound          = Error("achiever not found")
	ErrAchieverIncompleteDetails = Error("incomplete details for achiever")
)

// Error represents a general domain error.
type Error string

// Error returns the error message.
func (e Error) Error() string { return string(e) }
