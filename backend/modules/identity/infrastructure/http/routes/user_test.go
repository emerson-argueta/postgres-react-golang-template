package routes_test

import (
	"emersonargueta/m/v1/modules/identity/infrastructure/http/routes"
	"testing"
)

// NewIdentityHandler uses the labstack echo router.
func NewTestHandler() *routes.IdentityHandler {
	return routes.NewIdentityHandler("/test")
}
func TestLoginRoute(t *testing.T) {
	tHandler := NewTestHandler()

	tHandler.HandleRegisterUser(tHandler.Echo.AcquireContext())
}
