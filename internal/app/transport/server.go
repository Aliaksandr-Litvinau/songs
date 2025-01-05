package transport

import (
	"context"
	"fmt"
	"net/http"
)

type Server struct {
	httpServer  *http.Server
	songService SongService
}

func NewServer(addr string, songService SongService) *Server {
	server := &Server{
		songService: songService,
	}

	// Initialize router with service
	router := SetupRouter(songService)

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
