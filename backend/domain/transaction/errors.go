package transaction

// Transaction errors.
const (
	ErrTransactionNotFound        = Error("transaction not found")
	ErrTransactionDonatorNotFound = Error("transactions not found because donator does not exist")
	ErrTransactionChurchNotFound  = Error("transactions not found because church does not exist")
	ErrTransactionInternal        = Error("transaction internal error")
)

// Error represents a transaction error.
type Error string

// Error returns the error message.
func (e Error) Error() string { return string(e) }
