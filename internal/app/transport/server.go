package transport

import (
	"context"
	"fmt"
	"net/http"
	"songs/internal/app/service"
)

type Server struct {
	httpServer *http.Server
	service    service.SongService
}

func NewServer(addr string, service service.SongService) *Server {
	server := &Server{
		service: service,
	}

	// Initialize router with service
	router := SetupRouter(service)

	// Setup http server
	server.httpServer = &http.Server{
		Addr:    addr,
		Handler: router,
	}

	return server
}

func (s *Server) Run() error {
	if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("failed to start server: %w", err)
	}
	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	if err := s.httpServer.Shutdown(ctx); err != nil {
		return fmt.Errorf("failed to shutdown server: %w", err)
	}
	return nil
}
