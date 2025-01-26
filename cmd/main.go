package main

import (
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"log"
	"os"
	"os/signal"
	"songs/internal/app/config"
	"songs/internal/app/repository/pgrepo"
	"songs/internal/app/service"
	"songs/internal/app/transport/grpc"
	"songs/internal/app/transport/http"
	pg "songs/internal/pkg"
	"sync"
	"syscall"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
	os.Exit(0)
}

func run() error {
	// read config
	cfg := config.Read()

	pgDB, err := pg.Dial(cfg.DSN)
	if err != nil {
		return fmt.Errorf("pg.Dial failed: %w", err)
	}

	if pgDB != nil {
		if err := runPgMigrations(cfg.MigrationsPath, cfg.DSN); err != nil {
			return fmt.Errorf("runPgMigrations failed: %w", err)
		}
	}

	// Initialize repo
	songRepo := pgrepo.NewSongRepo(pgDB)
	// Initialize the song service
	songService := service.NewSongService(songRepo)

	// Create servers
	httpServer := http.NewServer(cfg.HTTPAddr, songService)
	grpcServer := grpc.NewServer(cfg.GRPCAddr, songService)

	// Channel for graceful shutdown
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	// WaitGroup for tracking running servers
	wg := &sync.WaitGroup{}
	wg.Add(2)

	// Start HTTP server
	go func() {
		defer wg.Done()
		if err := httpServer.Run(); err != nil {
			log.Printf("HTTP server error: %v", err)
		}
	}()

	// Start gRPC server
	go func() {
		defer wg.Done()
		if err := grpcServer.Run(); err != nil {
			log.Printf("gRPC server error: %v", err)
		}
	}()

	// Wait for shutdown signal
	<-shutdown
	log.Println("Shutting down servers...")

	// Graceful shutdown
	//httpServer.Stop()
	grpcServer.Stop()

	wg.Wait()
	log.Println("Servers stopped")

	return nil
}

func runPgMigrations(path, dsn string) error {
	if path == "" {
		return errors.New("no migrations path provided")
	}
	if dsn == "" {
		return errors.New("no DSN provided")
	}

	log.Println("Initializing migrations")
	m, err := migrate.New(path, dsn)
	if err != nil {
		return err
	}

	log.Println("Running migrations")
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}

	log.Println("Migrations completed successfully")
	return nil
}
