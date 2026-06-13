package configs

import (
	"errors"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	BotToken   string
}

func Load() (Config, error) {
	err := godotenv.Load()
	if err != nil {
		return Config{}, err
	}

	cfg := Config{
		DBHost:     os.Getenv("DB_HOST"),
		DBPort:     os.Getenv("DB_PORT"),
		DBUser:     os.Getenv("DB_USER"),
		DBPassword: os.Getenv("DB_PASSWORD"),
		DBName:     os.Getenv("DB_NAME"),
		BotToken:   os.Getenv("BOT_TOKEN"),
	}

	if cfg.DBHost == "" {
		return Config{}, errors.New("DB_HOST is required")
	}
	if cfg.DBPort == "" {
		return Config{}, errors.New("DB_PORT is required")
	}
	if cfg.DBUser == "" {
		return Config{}, errors.New("DB_USER is required")
	}
	if cfg.DBPassword == "" {
		return Config{}, errors.New("DB_PASSWORD is required")
	}
	if cfg.DBName == "" {
		return Config{}, errors.New("DB_NAME is required")
	}
	if cfg.BotToken == "" {
		return Config{}, errors.New("BOT_TOKEN is required")
	}

	return cfg, nil
}
