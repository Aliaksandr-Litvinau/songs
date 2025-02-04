package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type Config struct {
	HTTPAddr       string
	GRPCAddr       string
	DSN            string
	MigrationsPath string
	SwaggerURL     string
	Kafka          KafkaConfig
}

type KafkaConfig struct {
	Brokers        []string
	Topic          string
	GroupID        string
	SessionTimeout time.Duration
}

// GetConfig returns singleton instance of Config
func GetConfig() (Config, error) {
	timeout, err := getDurationFromEnv("KAFKA_SESSION_TIMEOUT", 10)
	if err != nil {
		return Config{}, fmt.Errorf("failed to get kafka timeout: %w", err)
	}

	return Config{
		HTTPAddr:       getEnv("HTTP_ADDR", ":8080"),
		GRPCAddr:       getEnv("GRPC_ADDR", ":50051"),
		DSN:            getEnv("DSN", "postgres://user:password@postgres:5432/music_library?sslmode=disable"),
		MigrationsPath: getEnv("MIGRATIONS_PATH", "file:///songs/internal/app/migrations"),
		Kafka: KafkaConfig{
			Brokers:        []string{getEnv("KAFKA_BROKERS", "kafka:9092")},
			Topic:          getEnv("KAFKA_TOPIC", "songs.updates"),
			GroupID:        getEnv("KAFKA_GROUP_ID", "songs_consumer_group"),
			SessionTimeout: timeout,
		},
	}, nil
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getDurationFromEnv(key string, defaultSeconds int) (time.Duration, error) {
	value := os.Getenv(key)
	if value == "" {
		return time.Duration(defaultSeconds) * time.Second, nil
	}

	seconds, err := strconv.Atoi(value)
	if err != nil {
		return 0, fmt.Errorf("invalid duration for %s: %w", key, err)
	}

	return time.Duration(seconds) * time.Second, nil
}
