package http_test

import (
	"fmt"
	"net/url"

	"emersonargueta/m/v1/delivery/http"
)

// MustOpenServerClient returns a running server and associated client. Panic on error.
func MustOpenServerClient() (*Server, *http.Client) {
	// Create and open test server.
	s := NewServer()
	if err := s.Open(); err != nil {
		panic(err)
	}

	// Create a client pointing to the server.
	c := http.NewClient()
	c.URL = &url.URL{Scheme: "http", Host: fmt.Sprintf("localhost:%d", s.Port())}

	return s, c
}
