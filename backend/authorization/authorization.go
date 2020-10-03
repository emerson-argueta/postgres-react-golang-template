package authorization

import (
	"trustdonations.org/m/v2/delivery/middleware"
)

// Actions provides functions that can be used to manage
// authorization for administrator.
type Actions interface {
	Authorize(token *middleware.TokenPair) (*middleware.TokenPair, error)
}
