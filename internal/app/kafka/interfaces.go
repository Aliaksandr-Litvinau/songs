package kafka

import (
	"context"
	"songs/internal/app/kafka/models"
)

// Producer defines an interface for sending messages to Kafka
type Producer interface {
	SendMessage(ctx context.Context, msg *models.Message) error
	Close() error
}

// Consumer defines an interface for receiving messages from Kafka
type Consumer interface {
	Start(ctx context.Context) error
	Close() error
}

// MessageHandler defines an interface for handling messages
type MessageHandler interface {
	Handle(ctx context.Context, msg *models.Message) error
}
