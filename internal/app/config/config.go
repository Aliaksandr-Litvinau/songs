package config

import (
	"github.com/joho/godotenv"
	"log"
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
	var config Config
	// Find .env file
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}

	// Getting and using a value from .env
	httpAddr := os.Getenv("HTTP_ADDR")
	if httpAddr != "" {
		config.HTTPAddr = httpAddr
	}
	dsn := os.Getenv("DSN")
	if dsn != "" {
		config.DSN = dsn
	}
	migrationsPath := os.Getenv("MIGRATIONS_PATH")
	if migrationsPath != "" {
		config.MigrationsPath = migrationsPath
	}
	return config
}
