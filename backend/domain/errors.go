package domain

// Domain errors.
const (
	ErrDomainInternal = Error("domain internal error")
)

// Error represents a general domain error.
type Error string

// Error returns the error message.
func (e Error) Error() string { return string(e) }
