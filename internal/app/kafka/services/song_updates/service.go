package song_updates

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"songs/internal/app/domain"
	"songs/internal/app/kafka"
)

// SongUpdateService controls the processing of song updates
type SongUpdateService struct {
	producer kafka.Producer
	consumer kafka.Consumer
	cancel   context.CancelFunc
	wg       sync.WaitGroup
	errChan  chan error
	done     chan struct{}
}

// SongUpdateHandler handles incoming song updates
type SongUpdateHandler struct{}

func (h *SongUpdateHandler) Handle(ctx context.Context, msg *kafka.Message) error {
	fmt.Printf("Got song update: ID=%d, GroupID=%d, Title=%s\n",
		msg.ID, msg.GroupID, msg.Title)
	return nil
}

// NewSongUpdateService creates a new SongUpdateService
func NewSongUpdateService(brokers []string, topic string) (*SongUpdateService, error) {
	cfg := &kafka.Config{
		Brokers: brokers,
		Topic:   topic,
	}

	producer, err := kafka.NewProducer(cfg)
	if err != nil {
		return nil, fmt.Errorf("error creating producer: %w", err)
	}

	handler := &SongUpdateHandler{}
	consumer, err := kafka.NewConsumer(cfg, handler)
	if err != nil {
		producer.Close()
		return nil, fmt.Errorf("error creating consumer: %w", err)
	}

	return &SongUpdateService{
		producer: producer,
		consumer: consumer,
		errChan:  make(chan error, 1),
		done:     make(chan struct{}),
	}, nil
}

func (s *SongUpdateService) Start(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	s.cancel = cancel

	// Run the update handler in a separate goroutine
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		if err := s.consumer.Start(ctx); err != nil {
			log.Printf("Error starting consumer: %v", err)
			s.errChan <- fmt.Errorf("consumer error: %w", err)
			cancel()
		}
	}()

	// Send test updates every 5 seconds
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()

		failureCount := 0
		const maxFailures = 3 // Maximum number of consecutive producer failures

		counter := 1
		for {
			select {
			case <-ticker.C:
				song := domain.Song{
					ID:          counter,
					GroupID:     counter%3 + 1,
					Title:       fmt.Sprintf("Song %d", counter),
					ReleaseDate: time.Now(),
					Text:        fmt.Sprintf("Text for song %d", counter),
					Link:        fmt.Sprintf("https://example.com/song/%d", counter),
				}

				msg := &kafka.Message{
					ID:          song.ID,
					GroupID:     song.GroupID,
					Title:       song.Title,
					ReleaseDate: song.ReleaseDate,
					Text:        song.Text,
					Link:        song.Link,
				}

				if err := s.producer.SendMessage(ctx, msg); err != nil {
					log.Printf("Error sending message: %v", err)
					failureCount++
					if failureCount >= maxFailures {
						log.Printf("Too many consecutive failures (%d), stopping service...", failureCount)
						s.errChan <- fmt.Errorf("too many producer failures: %w", err)
						cancel()
						return
					}
				} else {
					failureCount = 0
					log.Printf("Sent song update: ID=%d, GroupID=%d, Title=%s",
						msg.ID, msg.GroupID, msg.Title)
					counter++
				}
			case <-ctx.Done():
				log.Println("Context cancelled, stopping producer...")
				return
			}
		}
	}()

	// Wait for the service to stop
	go func() {
		defer close(s.done)
		select {
		case err := <-s.errChan:
			log.Printf("Service error occurred: %v", err)
			cancel()
		case <-ctx.Done():
			log.Println("Service context cancelled")
		}
		s.Stop() // Automatically stop the service
	}()

	return nil
}

func (s *SongUpdateService) Stop() {
	log.Println("Starting graceful shutdown of Kafka service...")
	if s.cancel != nil {
		s.cancel()
	}

	// Wait for all goroutines to complete
	done := make(chan struct{})
	go func() {
		s.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		log.Println("All goroutines completed successfully")
	case <-time.After(10 * time.Second):
		log.Println("Timeout waiting for goroutines to complete")
	}

	// Closing producer and consumer
	if err := s.producer.Close(); err != nil {
		log.Printf("Error closing producer: %v", err)
	}
	if err := s.consumer.Close(); err != nil {
		log.Printf("Error closing consumer: %v", err)
	}

	close(s.errChan)
	log.Println("Kafka service shutdown completed")
}
