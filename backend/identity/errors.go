package identity

// User errors.
const (
	ErrUserNotFound = Error("user not found")
	ErrUserExists   = Error("user exists")
)

// Error represents a general user error.
type Error string

// Error returns the error message.
func (e Error) Error() string { return string(e) }
