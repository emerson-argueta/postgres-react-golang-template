package postgres

// postgres client errors.
const (
	ErrPostgresClientNoValuesForFieldInFilter = Error("field specified to create query filter has no associated values")
)

// Error represents a Church error.
type Error string

// Error returns the error message.
func (e Error) Error() string { return string(e) }
