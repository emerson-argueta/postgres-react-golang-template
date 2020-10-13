package http

import (
	"net"
	"net/http"

	"emersonargueta/m/v1/config"
)

const (
	// Network type for use in server
	Network = "tcp"
)

// Server represents an HTTP server.
type Server struct {
	listener net.Listener
	Handler  *Handler
	Address  string
}

// NewServer returns a new instance of Server.
func NewServer(config *config.Config) *Server {
	return &Server{Address: config.Host + ":" + config.Port}
}

// Open a socket and serve the HTTP server.
func (s *Server) Open() error {
	openSocket, err := net.Listen(Network, s.Address)
	if err != nil {
		return err
	}

	s.listener = openSocket

	// Start HTTP server.
	go func() { http.Serve(s.listener, s.Handler) }()

	return nil
}

// Close the socket.
func (s *Server) Close() error {
	if s.listener != nil {
		s.listener.Close()
	}
	return nil
}

// Port that the server is open on. Only valid after open.
func (s *Server) Port() int {
	return s.listener.Addr().(*net.TCPAddr).Port
}
