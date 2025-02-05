package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"songs/internal/app/common"
	"songs/internal/app/config"
	"songs/internal/app/kafka/services/song_updates"
	"songs/internal/app/repository/pgrepo"
	"songs/internal/app/service"
	"songs/internal/app/transport"
	"songs/internal/app/transport/grpc"
	"songs/internal/app/transport/http"
	pg "songs/internal/pkg"
)

const (
	httpShutdownTimeout  = 10 * time.Second
	grpcShutdownTimeout  = 10 * time.Second
	kafkaShutdownTimeout = 10 * time.Second
	// dbShutdownTimeout    = 5 * time.Second
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
	os.Exit(0)
}

func run() error {
	cfg, err := config.GetConfig()
	if err != nil {
		return fmt.Errorf("config.GetConfig failed: %w", err)
	}

	// Initialize DB connection
	db, err := pg.Dial(cfg.DSN)
	if err != nil {
		return fmt.Errorf("pg.Dial failed: %w", err)
	}

	// Access the underlying sql.DB instance and call Close
	// https://forum.golangbridge.org/t/cant-close-db-connection-with-db-close/34657/2
	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("Error closing database connection: %v", err)
		}
	}()

	// Waiting for Kafka to be ready with timeout
	waitKafkaCtx, waitKafkaCancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer waitKafkaCancel()

	log.Println("Waiting for Kafka to be ready...")
	for {
		select {
		case <-waitKafkaCtx.Done():
			return waitKafkaCtx.Err()
		default:
			conn, err := net.Dial("tcp", cfg.Kafka.Brokers[0])
			if err == nil {
				if closeErr := conn.Close(); closeErr != nil {
					log.Printf("Warning: error closing test connection to Kafka: %v", closeErr)
				}
				log.Println("Kafka is ready")
				goto kafkaReady
			}
			log.Printf("Kafka is not ready yet: %v", err)
			time.Sleep(1 * time.Second)
		}
	}
kafkaReady:

	// Initialize repo and service
	repo := pgrepo.NewSongRepo(db)
	songService := service.NewSongService(repo)

	// Create HTTP server
	router := transport.SetupRouter(songService)
	httpServer := http.NewServer(cfg.HTTPAddr, router)

	// Create gRPC server
	grpcServer := grpc.NewServer(cfg.GRPCAddr, songService)

	// Create Kafka service
	kafkaService, err := song_updates.NewSongUpdateService(&cfg.Kafka)
	if err != nil {
		return fmt.Errorf("failed to create new song Kafka service: %w", err)
	}

	// Create a context with cancellation
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Run all services
	services := []struct {
		name    string
		service common.Lifecycle
	}{
		{"HTTP", httpServer},
		{"gRPC", grpcServer},
		{"Kafka", kafkaService},
	}

	for _, svc := range services {
		go func(name string, service common.Lifecycle) {
			log.Printf("Starting %s service...", name)
			if err := service.Start(ctx); err != nil {
				log.Printf("Error starting %s service: %v", name, err)
				cancel()
			}
		}(svc.name, svc.service)
	}

	// Waiting for a signal for a graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	log.Println("Received shutdown signal. Starting graceful shutdown...")

	// Cancel context to start shutdown process
	cancel()

	// stop services sequentially in the correct order
	// 1. HTTP and gRPC first (they accept incoming requests)
	httpCtx, httpCancel := context.WithTimeout(context.Background(), httpShutdownTimeout)
	defer httpCancel()
	if err := httpServer.Shutdown(httpCtx); err != nil {
		return fmt.Errorf("error stopping HTTP server: %w", err)
	}

	grpcCtx, grpcCancel := context.WithTimeout(context.Background(), grpcShutdownTimeout)
	defer grpcCancel()
	if err := grpcServer.Shutdown(grpcCtx); err != nil {
		return fmt.Errorf("error stopping gRPC server: %w", err)
	}

	// 2. Then Kafka service
	kafkaCtx, kafkaCancel := context.WithTimeout(context.Background(), kafkaShutdownTimeout)
	defer kafkaCancel()
	if err := kafkaService.Shutdown(kafkaCtx); err != nil {
		return fmt.Errorf("error stopping Kafka service: %w", err)
	}

	log.Println("Graceful shutdown completed")
	return nil
}
