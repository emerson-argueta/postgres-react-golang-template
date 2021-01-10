package http

import (
	"context"
	"emersonargueta/m/v1/shared/infrastructure"
	"emersonargueta/m/v1/shared/infrastructure/http/api"

	"net/http"
	"os"
	"os/signal"
	"time"
)

// Server for an api.
type Server struct {
	http.Server
	baseHandler *api.BaseHandler
	Addr        string
	config      *infrastructure.Config
}

// NewServer returns a new instance of Server.
func NewServer() *Server {
	config := infrastructure.NewConfig

	return &Server{Addr: ":" + config.HTTPServer.Port, config: config}
}

// Open opens a socket and serves the HTTP server.
func (s *Server) Open() error {

	baseHandler := api.NewBaseHandler(s.config.HTTPServer.APIBaseURL)
	s.Handler = baseHandler

	// Start HTTP server.
	go func() { http.ListenAndServe(s.Addr, baseHandler.ServeMux) }()

	// Block until an OS interrupt is received.
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	sig := <-stop
	println("Got signal:", sig)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := s.Shutdown(ctx); err != nil {
		return err
	}

	return nil
}
