package config

import (
	"os"
	"strconv"
	"sync"
	"time"
)

var (
	instance *Config
	once     sync.Once
)

type Config struct {
	HTTPAddr       string
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
func GetConfig() *Config {
	once.Do(func() {
		instance = &Config{
			HTTPAddr:       getEnv("HTTP_ADDR", ":8080"),
			DSN:            getEnv("DSN", "postgres://user:password@localhost:5432/music_library?sslmode=disable"),
			MigrationsPath: getEnv("MIGRATIONS_PATH", "file:///songs/internal/app/migrations"),
			Kafka: KafkaConfig{
				Brokers:        []string{getEnv("KAFKA_BROKERS", "kafka:9092")},
				Topic:          getEnv("KAFKA_TOPIC", "songs.updates"),
				GroupID:        getEnv("KAFKA_GROUP_ID", "songs_consumer_group"),
				SessionTimeout: time.Duration(getEnvAsInt("KAFKA_SESSION_TIMEOUT", 10)) * time.Second,
			},
		}
	})
	return instance
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists {
		if i, err := strconv.Atoi(value); err == nil {
			return i
		}
	}
	return defaultValue
}
