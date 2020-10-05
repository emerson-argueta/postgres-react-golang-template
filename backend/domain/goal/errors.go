package goal

// Church errors.
const (
	ErrChurchNotFound             = Error("church not found")
	ErrChurchExists               = Error("church already exists")
	ErrChurchInternal             = Error("church internal error")
	ErrChurchIncorrectCredentials = Error("incorrect credentials")

	ErrChurchDonatorExists = Error("could not add donator, donator already exists for church")

	ErrChurchAdministrators = Error("could not find administrators for church")

	ErrChurchCreatorDoesNotExists       = Error("could not find creator administrator for church")
	ErrChurchAdministratorDoesNotBelong = Error("administrator does not belong to church")

	ErrChurchFieldNotEditable       = Error("administrator cannot edit church field")
	ErrChurchFieldEditUnPriveledged = Error("administrator does not have priveledge to edit field")

	ErrChurchDonatorDoesNotExists = Error("could not find donator in church")
)

// Error represents a Church error.
type Error string

// Error returns the error message.
func (e Error) Error() string { return string(e) }
