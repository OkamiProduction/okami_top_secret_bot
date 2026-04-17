package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	TelegramToken string
	Debug         bool
	LogLevel      string
	LogFile       string
}

func Load() (*Config, error) {
	_ = godotenv.Load() // ошибка не критична

	cfg := &Config{
		TelegramToken: os.Getenv("TELEGRAM_TOKEN"),
		Debug:         strings.ToLower(os.Getenv("DEBUG")) == "true",
		LogLevel:      strings.ToLower(os.Getenv("LOG_LEVEL")),
		LogFile:       os.Getenv("LOG_FILE"),
	}

	// значения по умолчанию
	if cfg.LogLevel == "" {
		cfg.LogLevel = "info"
	}
	if cfg.LogFile == "" {
		cfg.LogFile = "bot.log" // или /var/log/tgbot/bot.log для сервера
	}

	if cfg.TelegramToken == "" {
		// ошибку вернём, но без log.Fatal — пусть вызывающий решает
		return nil, fmt.Errorf("TELEGRAM_TOKEN is required")
	}

	return cfg, nil
}
