package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"

	"songs/internal/app/config"
	"songs/internal/app/kafka/models"

	"github.com/IBM/sarama"
)

type kafkaConsumer struct {
	group   sarama.ConsumerGroup
	topics  []string
	handler MessageHandler
	wg      sync.WaitGroup
}

type consumerGroupHandler struct {
	handler MessageHandler
	wg      *sync.WaitGroup
	errChan chan error
}

func (h *consumerGroupHandler) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (h *consumerGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }

func (h *consumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for {
		select {
		case msg, ok := <-claim.Messages():
			if !ok {
				return nil
			}

			var kafkaMsg models.Message
			if err := json.Unmarshal(msg.Value, &kafkaMsg); err != nil {
				log.Printf("Failed to unmarshal message from partition %d: %v", claim.Partition(), err)
				continue
			}

			if err := h.handler.Handle(session.Context(), &kafkaMsg); err != nil {
				log.Printf("Failed to handle message from partition %d: %v", claim.Partition(), err)
				select {
				case h.errChan <- fmt.Errorf("critical error handling message from partition %d: %w", claim.Partition(), err):
				default:
					log.Printf("Additional error occurred on partition %d: %v", claim.Partition(), err)
				}
			}

			session.MarkMessage(msg, "")

		case <-session.Context().Done():
			return nil
		}
	}
}

func NewConsumer(cfg *config.KafkaConfig, handler MessageHandler) (Consumer, error) {
	saramaConfig := sarama.NewConfig()
	saramaConfig.Consumer.Group.Session.Timeout = cfg.SessionTimeout
	saramaConfig.Consumer.Offsets.Initial = sarama.OffsetNewest
	saramaConfig.Version = sarama.V2_8_1_0

	group, err := sarama.NewConsumerGroup(cfg.Brokers, cfg.GroupID, saramaConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create consumer group: %w", err)
	}

	return &kafkaConsumer{
		group:   group,
		topics:  []string{cfg.Topic},
		handler: handler,
	}, nil
}

func (c *kafkaConsumer) Start(ctx context.Context) error {
	errChan := make(chan error, 1)

	c.wg.Add(1)
	go func() {
		defer c.wg.Done()

		handler := &consumerGroupHandler{
			handler: c.handler,
			wg:      &c.wg,
			errChan: errChan,
		}

		for {
			if err := c.group.Consume(ctx, c.topics, handler); err != nil {
				if err == sarama.ErrClosedConsumerGroup {
					return
				}
				log.Printf("Error from consumer group: %v", err)
				select {
				case errChan <- fmt.Errorf("consumer group error: %w", err):
				default:
					log.Printf("Additional consumer group error: %v", err)
				}
			}

			if ctx.Err() != nil {
				return
			}
		}
	}()

	select {
	case err := <-errChan:
		return fmt.Errorf("consumer error: %w", err)
	case <-ctx.Done():
		return nil
	}
}

func (c *kafkaConsumer) Close() error {
	if err := c.group.Close(); err != nil {
		return fmt.Errorf("failed to close consumer group: %w", err)
	}

	c.wg.Wait()
	return nil
}
