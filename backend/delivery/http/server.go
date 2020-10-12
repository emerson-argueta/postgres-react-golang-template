package http

import (
	"io"
	"net"
	"net/http"
	"os"
	"testing"

	"emersonargueta/m/v1/config"
)

// Server represents an HTTP server.
type Server struct {
	ln net.Listener

	// Handler to serve.
	Handler *Handler

	// Bind address to open.
	Addr string
}

// NewServer returns a new instance of Server.
func NewServer(config *config.Config) *Server {

	return &Server{Addr: config.Host + ":" + config.Port}
}

// Open opens a socket and serves the HTTP server.
func (s *Server) Open() error {
	// Open socket.
	ln, err := net.Listen("tcp", s.Addr)
	if err != nil {
		return err
	}
	s.ln = ln

	// Start HTTP server.
	go func() { http.Serve(s.ln, s.Handler) }()

	return nil
}

// Close closes the socket.
func (s *Server) Close() error {
	if s.ln != nil {
		s.ln.Close()
	}
	return nil
}

// Port returns the port that the server is open on. Only valid after open.
func (s *Server) Port() int {
	return s.ln.Addr().(*net.TCPAddr).Port
}

// VerboseWriter returns a multi-writer to STDERR and w if the "-v" flag is set.
func VerboseWriter(w io.Writer) io.Writer {
	if testing.Verbose() {
		return io.MultiWriter(w, os.Stderr)
	}
	return w
}
