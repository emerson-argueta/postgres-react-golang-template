package authorization

import (
	"emersonargueta/m/v1/delivery/middleware"
)

// Actions provides functions that can be used to manage
// authorization for administrator.
type Actions interface {
	Authorize(token *middleware.TokenPair) (*middleware.TokenPair, error)
}
