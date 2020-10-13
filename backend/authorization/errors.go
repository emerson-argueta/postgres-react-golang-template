package authorization

// jwt authorization errors.
const (
	ErrAuthorizationFailed = Error("could not authorize")
)

// Error represents a general middleware error.
type Error string

// Error returns the error message.
func (e Error) Error() string { return string(e) }
