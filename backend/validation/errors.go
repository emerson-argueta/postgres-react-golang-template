package validation

// Validation errors.
const (
	ErrValidationPassword    = Error("password could not be validated, must be at least 8 characters")
	ErrValidationUserEmail   = Error("user email could not be validated")
	ErrValidationPhoneNumber = Error("phone number could not be validated")
)

// Error represents a admin error.
type Error string

// Error returns the error message.
func (e Error) Error() string { return string(e) }
