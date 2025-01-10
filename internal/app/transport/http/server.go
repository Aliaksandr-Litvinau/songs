package http

import (
	"context"
	"fmt"
	"net/http"
	"songs/internal/app/transport"
)

type Server struct {
	httpServer  *http.Server
	songService transport.SongService
}

func NewServer(addr string, songService transport.SongService) *Server {
	server := &Server{
		songService: songService,
	}

	// Initialize router with service
	router := transport.SetupRouter(songService)

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
