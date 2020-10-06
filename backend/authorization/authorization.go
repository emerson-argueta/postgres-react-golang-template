package authorization

import (
	"emersonargueta/m/v1/delivery/middleware"
)

// Service provides processes used in authorization.
type Service interface {
	Authorize(token *middleware.TokenPair, uuid string) (*middleware.TokenPair, error)
}
