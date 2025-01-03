package config

import (
	"os"
)

type Config struct {
	HTTPAddr       string
	DSN            string
	MigrationsPath string
	SwaggerURL     string
}

// Read reads config from environment.
func Read() Config {
	return Config{
		HTTPAddr:       getEnv("HTTP_ADDR", ":8080"),
		DSN:            getEnv("DSN", "postgres://user:password@localhost:5432/music_library?sslmode=disable"),
		MigrationsPath: getEnv("MIGRATIONS_PATH", "file:///songs/internal/app/migrations"),
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
