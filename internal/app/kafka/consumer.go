package kafka

import (
	"context"
	"encoding/json"
	"log"

	"github.com/IBM/sarama"
)

type MessageHandler interface {
	Handle(ctx context.Context, msg *Message) error
}

type Consumer interface {
	Start(ctx context.Context) error
	Close() error
}

type kafkaConsumer struct {
	consumer sarama.Consumer
	topic    string
	handler  MessageHandler
}

func NewConsumer(cfg *Config, handler MessageHandler) (Consumer, error) {
	config := sarama.NewConfig()
	config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin
	config.Consumer.Offsets.Initial = sarama.OffsetNewest

	consumer, err := sarama.NewConsumer(cfg.Brokers, config)
	if err != nil {
		return nil, err
	}

	return &kafkaConsumer{
		consumer: consumer,
		topic:    cfg.Topic,
		handler:  handler,
	}, nil
}

func (c *kafkaConsumer) Start(ctx context.Context) error {
	partitions, err := c.consumer.Partitions(c.topic)
	if err != nil {
		return err
	}

	for _, partition := range partitions {
		pc, err := c.consumer.ConsumePartition(c.topic, partition, sarama.OffsetNewest)
		if err != nil {
			return err
		}

		go func(pc sarama.PartitionConsumer) {
			defer pc.Close()
			for {
				select {
				case msg := <-pc.Messages():
					var kafkaMsg Message
					if err := json.Unmarshal(msg.Value, &kafkaMsg); err != nil {
						log.Printf("Failed to unmarshal message: %v", err)
						continue
					}

					if err := c.handler.Handle(ctx, &kafkaMsg); err != nil {
						log.Printf("Failed to handle message: %v", err)
					}
				case <-ctx.Done():
					return
				}
			}
		}(pc)
	}

	return nil
}

func (c *kafkaConsumer) Close() error {
	return c.consumer.Close()
}
