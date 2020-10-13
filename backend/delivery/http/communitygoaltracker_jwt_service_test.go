package http_test

import (
	"bytes"
	mockcommunitygoaltracker "emersonargueta/m/v1/communitygoaltracker/mock"
	"emersonargueta/m/v1/config"
	"emersonargueta/m/v1/delivery/http"
	"emersonargueta/m/v1/delivery/middleware"
	"log"
	"testing"
)

// CommunitygoaltrackerHandler represents a test wrapper for http.AdminHandler.
type CommunitygoaltrackerHandler struct {
	*http.CommunitygoaltrackerHandler
	CommunitygoaltrackerProcesses mockcommunitygoaltracker.CommunitygoaltrackerProcesses
	LogOutput                     bytes.Buffer
}

// NewCommunitygoaltrackerHandler returns a CommunitygoaltrackerHandler using mock implementation of services.
func NewCommunitygoaltrackerHandler() *CommunitygoaltrackerHandler {
	h := &CommunitygoaltrackerHandler{
		CommunitygoaltrackerHandler: http.NewCommunitygoaltrackerHandler(
			middleware.NewClient(config.NewConfig()).JwtService(),
		),
	}

	h.CommunitygoaltrackerHandler.Communitygoaltracker.Processes = &h.CommunitygoaltrackerProcesses

	h.Logger = log.New(VerboseWriter(&h.LogOutput), "", log.LstdFlags)
	return h
}

func TestCommunitygoaltrackerService_Register(t *testing.T) {
	t.Run("OK", testCommunitygoaltrackerService_Register)
	t.Run("ErrAdminExists", testCommunitygoaltrackerService_Register_ErrUserExists)
}

func testCommunitygoaltrackerService_Register(t *testing.T) {

}
func testCommunitygoaltrackerService_Register_ErrUserExists(t *testing.T) {

}
