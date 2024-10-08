package config

import (
	"github.com/spf13/viper"
	"log"
)

type Config struct {
	Port        string
	DatabaseURL string
	SwaggerURL  string
}

func LoadConfig() *Config {
	viper.SetConfigFile(".env")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}
	return &Config{
		Port:        viper.GetString("PORT"),
		DatabaseURL: viper.GetString("DATABASE_URL"),
		SwaggerURL:  viper.GetString("SWAGGER_URL"),
	}
}
