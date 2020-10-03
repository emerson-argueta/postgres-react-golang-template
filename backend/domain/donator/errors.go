package donator

// Donator errors.
const (
	ErrDonatorNotFound = Error("donator not found")
	ErrDonatorExists   = Error("donator already exists")
	ErrDonatorInternal = Error("donator internal error")

	ErrDonatorFieldNotEditable = Error("administrator cannot edit donator field")
)

// Error represents a donator error.
type Error string

// Error returns the error message.
func (e Error) Error() string { return string(e) }
