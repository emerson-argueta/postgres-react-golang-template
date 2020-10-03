package stripe

// Domain errors.
const (
	ErrStripeNewSubscription = Error("could not create stripe subscription")
)

// Error represents a general domain error.
type Error string

// Error returns the error message.
func (e Error) Error() string { return string(e) }
