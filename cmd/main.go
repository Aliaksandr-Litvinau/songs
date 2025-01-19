package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"songs/internal/app/config"
	"songs/internal/app/kafka/services/runner"
	"songs/internal/app/repository/pgrepo"
	"songs/internal/app/service"
	"songs/internal/app/transport"
	pg "songs/internal/pkg"
	"syscall"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
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

	// Create context with cancellation
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start Kafka service
	stopKafka, err := runner.RunKafkaService(ctx)
	if err != nil {
		return fmt.Errorf("failed to start kafka service: %w", err)
	}
	defer stopKafka()

	// Create HTTP server
	server := transport.NewServer(cfg.HTTPAddr, &songService)

	// Start HTTP server in a goroutine
	go func() {
		if err := server.Run(); err != nil {
			log.Printf("failed to run server: %v", err)
			cancel()
		}
	}()

	// Wait for interrupt signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	sig := <-sigChan
	log.Printf("Received signal %v, shutting down...", sig)

	// Trigger graceful shutdown
	cancel()
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
