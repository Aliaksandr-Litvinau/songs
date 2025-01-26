package kafka

import (
	"context"
	"encoding/json"

	"songs/internal/app/config"
	"songs/internal/app/kafka/models"

	"github.com/IBM/sarama"
)

type Producer interface {
	SendMessage(ctx context.Context, msg *models.Message) error
	Close() error
}

type producer struct {
	syncProducer sarama.SyncProducer
	topic        string
}

func NewProducer(cfg *config.KafkaConfig) (Producer, error) {
	saramaConfig := sarama.NewConfig()
	saramaConfig.Producer.Return.Successes = true
	saramaConfig.Producer.RequiredAcks = sarama.WaitForAll
	saramaConfig.Producer.Retry.Max = 5

	syncProducer, err := sarama.NewSyncProducer(cfg.Brokers, saramaConfig)
	if err != nil {
		return nil, err
	}

	return &producer{
		syncProducer: syncProducer,
		topic:        cfg.Topic,
	}, nil
}

func (p *producer) SendMessage(ctx context.Context, msg *models.Message) error {
	value, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	_, _, err = p.syncProducer.SendMessage(&sarama.ProducerMessage{
		Topic: p.topic,
		Value: sarama.ByteEncoder(value),
	})

	return err
}

func (p *producer) Close() error {
	return p.syncProducer.Close()
}
