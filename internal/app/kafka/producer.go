package kafka

import (
	"context"
	"encoding/json"

	"github.com/IBM/sarama"
)

type Producer interface {
	SendMessage(ctx context.Context, msg *Message) error
	Close() error
}

type producer struct {
	syncProducer sarama.SyncProducer
	topic        string
}

func NewProducer(cfg *Config) (Producer, error) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5

	syncProducer, err := sarama.NewSyncProducer(cfg.Brokers, config)
	if err != nil {
		return nil, err
	}

	return &producer{
		syncProducer: syncProducer,
		topic:        cfg.Topic,
	}, nil
}

func (p *producer) SendMessage(ctx context.Context, msg *Message) error {
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
