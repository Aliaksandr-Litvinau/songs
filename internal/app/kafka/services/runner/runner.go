package runner

import (
	"context"
	"fmt"
	"log"
	"net"
	"songs/internal/app/kafka/services/song_updates"
	"time"
)

// waitForKafka try to connect to Kafka at the specified address
func waitForKafka(address string, timeout time.Duration) error {
	start := time.Now()
	for {
		conn, err := net.Dial("tcp", address)
		if err == nil {
			conn.Close()
			return nil
		}

		if time.Since(start) > timeout {
			return fmt.Errorf("timeout waiting for Kafka at %s: %v", address, err)
		}

		log.Printf("Waiting for Kafka to be available at %s...", address)
		time.Sleep(1 * time.Second)
	}
}

// RunKafkaService starts the Kafka service and returns a function for graceful shutdown
func RunKafkaService(ctx context.Context) (func(), error) {
	if err := waitForKafka("kafka:9092", 60*time.Second); err != nil {
		return nil, fmt.Errorf("kafka is not available: %v", err)
	}

	// Creating song updates service with Kafka
	updateService, err := song_updates.NewSongUpdateService(
		[]string{"kafka:9092"},
		"songs.updates",
	)
	if err != nil {
		return nil, fmt.Errorf("error creating song updates service: %v", err)
	}

	if err := updateService.Start(ctx); err != nil {
		return nil, fmt.Errorf("error starting song updates service: %v", err)
	}

	return func() {
		fmt.Println("Stopping Kafka service...")
		updateService.Stop()
	}, nil
}