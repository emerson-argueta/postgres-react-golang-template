package identity

// User model errors.
const (
	ErrUserNotFound = Error("user not found")
	ErrUserExists   = Error("user exists")
)

// Domain model errors
const (
	ErrDomainNotFound = Error("domain not found")
	ErrDomainExists   = Error("domain exists")
)

// Error represents a general identity error.
type Error string

// Error returns the error message.
func (e Error) Error() string { return string(e) }
