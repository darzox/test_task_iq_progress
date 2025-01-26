package config

import (
	"log"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type Config struct {
	HttpPort int `env:"HTTP_PORT,required"`

	PostgresHost     string `env:"POSTGRES_HOST,required"`
	PostgresPort     int    `env:"POSTGRES_PORT,required"`
	PostgresUser     string `env:"POSTGRES_USER,required"`
	PostgresPassword string `env:"POSTGRES_PASSWORD,required"`
	PostgresDb       string `env:"POSTGRES_DB,required"`
	PostgresSslmode  string `env:"POSTGRES_SSLMODE" envDefault:"require"`
}

func NewConfig() (*Config, error) {
	// Load .env file for local development
	err := godotenv.Load()
	if err != nil {
		log.Printf("Error loading .env file: %v", err)
	}

	cfg := new(Config)

	// Parse environment variables into the Config struct
	if err := env.Parse(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
