package http_test

import (
	"emersonargueta/m/v1/delivery/http"
)

// Handler represents a test wrapper for http.Handler.
type Handler struct {
	*http.Handler

	*CommunitygoaltrackerHandler
}

// NewHandler returns a new instance of Handler.
func NewHandler() *Handler {
	h := &Handler{
		Handler:                     &http.Handler{},
		CommunitygoaltrackerHandler: NewCommunitygoaltrackerHandler(),
	}
	h.Handler.CommunitygoaltrackerHandler = h.CommunitygoaltrackerHandler.CommunitygoaltrackerHandler
	return h
}
