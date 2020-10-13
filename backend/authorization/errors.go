package authorization

// jwt authorization errors.
const (
	ErrAuthorizationFailed        = Error("could not authorize")
	ErrAuthorizationKeyNotCreated = Error("could not create a the authorization key")
)

// Error represents a general middleware error.
type Error string

// Error returns the error message.
func (e Error) Error() string { return string(e) }
