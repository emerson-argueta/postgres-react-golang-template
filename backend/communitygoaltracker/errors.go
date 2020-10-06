package communitygoaltracker

// Domain errors.
const (
	ErrInternal = Error("internal domain error")
)

// Error represents a general domain error.
type Error string

// Error returns the error message.
func (e Error) Error() string { return string(e) }
