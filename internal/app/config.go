package app

import (
	"os"
)

// Config представляет конфигурацию приложения
type Config struct {
	Host string
	Port string
	DB   struct {
		DSN string
	}
}

// NewConfig создает конфигурацию приложения из переменных окружения
func NewConfig(configPath string) (*Config, error) {
	return &Config{
		Host: "0.0.0.0",
		Port: os.Getenv("SERVER_PORT"),
		DB: struct {
			DSN string
		}{
			DSN: os.Getenv("DATABASE_URL"),
		},
	}, nil
}
