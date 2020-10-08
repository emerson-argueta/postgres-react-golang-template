package communitygoaltracker

// Goal model errors.
const (
	ErrGoalExists            = Error("goal exists")
	ErrGoalNotFound          = Error("goal not found")
	ErrIncompleteGoalDetails = Error("incomplete details for goal")
)

// Achiever moodel errors.
const (
	ErrAchieverExists            = Error("achiever exists")
	ErrAchieverNotFound          = Error("achiever not found")
	ErrIncompleteAchieverDetails = Error("incomplete details for achiever")
)

// Error represents a general domain error.
type Error string

// Error returns the error message.
func (e Error) Error() string { return string(e) }
