package authorization

// Processes used to authenticate.
type Processes interface {
	// Authorize a user using a key. Return's the user's uuid if succesful,or
	// ErrAuthorizationFailed
	Authorize(key interface{}) (*string, error)
	// NewKey for a user is created. Returns a map containing key information if
	// successful, or ErrAuthorizationKeyNotCreated.
	NewKey(uuid string) (map[string]string, error)
	// ReAuthorize a user that has a valid but expired key. Returns a map
	// containing key information if successful, or
	// ErrAuthorizationKeyNotCreated.
	ReAuthorize(key interface{}) (map[string]string, error)
}
