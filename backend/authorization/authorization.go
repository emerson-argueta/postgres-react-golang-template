package authorization

// Processes used to authenticate.
type Processes interface {
	// Authorize a user using a key. Return's the user's uuid if succesful,or ErrAuthorizationFailed
	Authorize(key interface{}) (*string, error)
}
