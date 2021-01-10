package achiever

// Achiever moodel errors.
const (
	ErrAchieverExists            = Error("achiever exists")
	ErrAchieverNotFound          = Error("achiever not found")
	ErrAchieverIncompleteDetails = Error("incomplete details for achiever")
	ErrAchieverInvalidRole       = Error("invalid role for achiever")
	ErrAchieverGoalNotFound      = Error("Achiever's goal not found")
)

// Error represents a general domain error.
type Error string

// Error returns the error message.
func (e Error) Error() string { return string(e) }
