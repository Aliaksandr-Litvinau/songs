package main

import (
	"context"
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
	cfg := config.GetConfig()

	// Initialize DB connection
	db, err := pg.Dial(cfg.DSN)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Failed to get database instance: %v", err)
	}

	// Waiting for Kafka to be ready with timeout
	waitKafkaCtx, waitKafkaCancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer waitKafkaCancel()

	log.Println("Waiting for Kafka to be ready...")
	for {
		select {
		case <-waitKafkaCtx.Done():
			log.Fatalf("Timeout waiting for Kafka: %v", waitKafkaCtx.Err())
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
		log.Fatalf("Failed to create Kafka service: %v", err)
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
		log.Printf("Error stopping HTTP server: %v", err)
	}

	grpcCtx, grpcCancel := context.WithTimeout(context.Background(), grpcShutdownTimeout)
	defer grpcCancel()
	if err := grpcServer.Shutdown(grpcCtx); err != nil {
		log.Printf("Error stopping gRPC server: %v", err)
	}

	// 2. Then Kafka service
	kafkaCtx, kafkaCancel := context.WithTimeout(context.Background(), kafkaShutdownTimeout)
	defer kafkaCancel()
	if err := kafkaService.Shutdown(kafkaCtx); err != nil {
		log.Printf("Error stopping Kafka service: %v", err)
	}

	// 3. At the end, close the connection to the database
	if err := sqlDB.Close(); err != nil {
		log.Printf("Error closing database connection: %v", err)
	}

	log.Println("Graceful shutdown completed")
}
