package http

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
)

type Server struct {
	server *http.Server
	addr   string
}

func NewServer(addr string, handler http.Handler) *Server {
	return &Server{
		server: &http.Server{
			Addr:    addr,
			Handler: handler,
		},
		addr: addr,
	}
}

func (s *Server) Start(_ context.Context) error {
	log.Printf("Starting HTTP server on %s", s.addr)
	if err := s.server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("failed to start HTTP server: %w", err)
	}
	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	log.Println("Stopping HTTP server...")
	if err := s.server.Shutdown(ctx); err != nil {
		return fmt.Errorf("failed to stop HTTP server: %w", err)
	}
	log.Println("HTTP server stopped")
	return nil
}
