package config

import (
	"os"
	"time"
)

type AppConfig struct {
	Environment    string
	ServerPort     string
	WebhookSecret  string
	IStarConfigVar IStarConfig
}

type IStarConfig struct {
	APIKey     string
	BaseURL    string
	Timeout    time.Duration
	MaxRetries int
}

func Load() *AppConfig {
	return &AppConfig{
		Environment:   os.Getenv("ENV"),
		ServerPort:    os.Getenv("PORT"),
		WebhookSecret: os.Getenv("WEBHOOK_SECRET"),
		IStarConfigVar: IStarConfig{
			APIKey:     os.Getenv("ISTAR_API_KEY"),
			BaseURL:    os.Getenv("ISTAR_BASE_URL"),
			Timeout:    10 * time.Second,
			MaxRetries: 3,
		},
	}
}
