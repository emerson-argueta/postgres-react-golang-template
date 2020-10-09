package http_test

import (
	"bytes"
	mockcommunitygoaltracker "emersonargueta/m/v1/communitygoaltracker/mock"
	"emersonargueta/m/v1/delivery/http"
	"log"
)

// CommunitygoaltrackerHandler represents a test wrapper for http.AdminHandler.
type CommunitygoaltrackerHandler struct {
	*http.CommunitygoaltrackerHandler
	CommunitygoaltrackerService mockcommunitygoaltracker.Communitygoaltrackerservice
	LogOutput                   bytes.Buffer
}

// NewCommunitygoaltrackerHandler returns a CommunitygoaltrackerHandler.
func NewCommunitygoaltrackerHandler() *CommunitygoaltrackerHandler {
	h := &CommunitygoaltrackerHandler{CommunitygoaltrackerHandler: http.NewCommunitygoaltrackerHandler()}

	h.CommunitygoaltrackerHandler.Communitygoaltracker.Service = &h.CommunitygoaltrackerService

	h.Logger = log.New(VerboseWriter(&h.LogOutput), "", log.LstdFlags)
	return h
}
