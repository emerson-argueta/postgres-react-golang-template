package identity

// User model errors.
const (
	ErrUserNotFound             = Error("user not found")
	ErrUserExists               = Error("user exists")
	ErrUserIncompleteDetails    = Error("incomplete details for user")
	ErrUserIncorrectCredentials = Error("incorrect credentials")
)

// Error represents a general identity error.
type Error string

// Error returns the error message.
func (e Error) Error() string { return string(e) }
