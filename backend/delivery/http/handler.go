package http

import (
	"net/http"
	"strings"
)

// Handler is a collection of all the service handlers.
type Handler struct {
	*CommunitygoaltrackerHandler
}

// ServeHTTP delegates a request to the appropriate subhandler.
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if strings.HasPrefix(r.URL.Path, RoutePrefix+CommunitygoalTrackerURLPrefix) {
		h.CommunitygoaltrackerHandler.ServeHTTP(w, r)
	} else {
		http.NotFound(w, r)
	}
}
