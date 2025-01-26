package song_updates

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"songs/internal/app/config"
	"songs/internal/app/domain"
	"songs/internal/app/kafka"
	"songs/internal/app/kafka/interfaces"
	"songs/internal/app/kafka/models"

	"golang.org/x/sync/errgroup"
)

// SongUpdateService управляет обработкой обновлений песен
type SongUpdateService struct {
	producer interfaces.Producer
	consumer interfaces.Consumer
	cancel   context.CancelFunc
	wg       sync.WaitGroup
	done     chan struct{}
}

// SongUpdateHandler обрабатывает входящие обновления песен
type SongUpdateHandler struct{}

func (h *SongUpdateHandler) Handle(ctx context.Context, msg *models.Message) error {
	fmt.Printf("Got song update: ID=%d, GroupID=%d, Title=%s\n",
		msg.ID, msg.GroupID, msg.Title)
	return nil
}

// NewSongUpdateService создает новый SongUpdateService
func NewSongUpdateService(cfg *config.KafkaConfig) (*SongUpdateService, error) {
	producer, err := kafka.NewProducer(cfg)
	if err != nil {
		return nil, fmt.Errorf("error creating producer: %w", err)
	}

	handler := &SongUpdateHandler{}
	consumer, err := kafka.NewConsumer(cfg, handler)
	if err != nil {
		if closeErr := producer.Close(); closeErr != nil {
			return nil, fmt.Errorf("error creating consumer: %w, failed to close producer: %v", err, closeErr)
		}
		return nil, fmt.Errorf("error creating consumer: %w", err)
	}

	return &SongUpdateService{
		producer: producer,
		consumer: consumer,
		done:     make(chan struct{}),
	}, nil
}

func (s *SongUpdateService) Start(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	s.cancel = cancel

	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		if err := s.consumer.Start(ctx); err != nil {
			log.Printf("Error starting consumer: %v", err)
			return fmt.Errorf("consumer error: %w", err)
		}
		return nil
	})

	g.Go(func() error {
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()

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

				msg := &models.Message{
					ID:          song.ID,
					GroupID:     song.GroupID,
					Title:       song.Title,
					ReleaseDate: song.ReleaseDate,
					Text:        song.Text,
					Link:        song.Link,
				}

				if err := s.producer.SendMessage(ctx, msg); err != nil {
					log.Printf("Error sending message: %v", err)
					return fmt.Errorf("producer error: %w", err)
				}

				log.Printf("Sent song update: ID=%d, GroupID=%d, Title=%s",
					msg.ID, msg.GroupID, msg.Title)
				counter++

			case <-ctx.Done():
				return nil
			}
		}
	})

	go func() {
		defer close(s.done)
		if err := g.Wait(); err != nil {
			log.Printf("Service error occurred: %v", err)
		}
	}()

	return nil
}

func (s *SongUpdateService) Shutdown(ctx context.Context) error {
	log.Println("Starting graceful shutdown of Kafka service...")

	if s.cancel != nil {
		s.cancel()
	}

	// Wait for all goroutines to complete with a timeout
	select {
	case <-s.done:
	case <-ctx.Done():
		return fmt.Errorf("kafka service shutdown timeout exceeded")
	}

	var errs []error
	if err := s.producer.Close(); err != nil {
		errs = append(errs, fmt.Errorf("error closing producer: %w", err))
	}
	if err := s.consumer.Close(); err != nil {
		errs = append(errs, fmt.Errorf("error closing consumer: %w", err))
	}

	log.Println("Kafka service shutdown completed")

	if len(errs) > 0 {
		return fmt.Errorf("kafka service shutdown errors: %v", errs)
	}
	return nil
}
