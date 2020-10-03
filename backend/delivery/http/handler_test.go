package http_test

import (
	"trustdonations.org/m/v2/delivery/http"
)

// Handler represents a test wrapper for http.Handler.
type Handler struct {
	*http.Handler

	AdministratorHandler *AdministratorHandler
}

// NewHandler returns a new instance of Handler.
func NewHandler() *Handler {
	h := &Handler{
		Handler:              &http.Handler{},
		AdministratorHandler: NewAdministratorHandler(),
	}
	h.Handler.AdministratorHandler = h.AdministratorHandler.AdministratorHandler
	return h
}
