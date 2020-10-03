package administrator

// Admin errors.
const (
	ErrAdministratorNotAuthorized = Error("administrator not authorized to perform action")

	ErrAdministratorNotFound             = Error("administrator not found")
	ErrAdministratorExists               = Error("administrator already exists")
	ErrAdministratorInternal             = Error("administrator internal error")
	ErrAdministratorIncorrectCredentials = Error("incorrect credentials")
	ErrAdministratorIncompleteDetails    = Error("could not register administrator, insufficient details to register administrator")
	ErrAdministratorFieldNotEditable     = Error("administrator cannot edit administrator field")

	ErrAdministratorRegister = Error("error when registering administrator")
	ErrAdministratorLogin    = Error("error for with administrator login")

	ErrAdministratorSubscriptionPaymentGateway = Error("could not retrieve payment gateway details")

	ErrAdministratorChurchExists               = Error("could not add church, administrator is already part of the church")
	ErrAdministratorChurchIncompleteDetails    = Error("could not add or create church, insufficient details to add or create church")
	ErrAdministratorChurchIncorrectCredentilas = Error("could not add church, incorrect credential")
	ErrAdministratorNoChurches                 = Error("could not read churches, administrator has not created or belong to any church")
	ErrAdministratorDoesNotBelongToChurch      = Error("administrator does not belong this church")

	ErrAdministratorDonatorIncompleteDetails = Error("could not add donator, insufficient details to add donator")

	ErrAdministratorFreeUsageLimitReached = Error("Administrator has reached free usage limit, please upgrade plan to continue use")
)

// Error represents a admin error.
type Error string

// Error returns the error message.
func (e Error) Error() string { return string(e) }
