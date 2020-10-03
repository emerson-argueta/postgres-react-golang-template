package http_test

import (
	"io"
	"os"
	"testing"
	"time"

	"trustdonations.org/m/v2/delivery/http"
)

// Now represents the mocked current time.
var Now = time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC)

// Server represents a test wrapper for http.Server.
type Server struct {
	*http.Server

	Handler *Handler
}

// NewServer returns a new instance of Server.
func NewServer() *Server {
	s := &Server{
		Server:  http.NewServer(),
		Handler: NewHandler(),
	}
	s.Server.Handler = s.Handler.Handler

	// Use random port.
	s.Addr = ":0"

	return s
}

// VerboseWriter returns a multi-writer to STDERR and w if the "-v" flag is set.
func VerboseWriter(w io.Writer) io.Writer {
	if testing.Verbose() {
		return io.MultiWriter(w, os.Stderr)
	}
	return w
}
