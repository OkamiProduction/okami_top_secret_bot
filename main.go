package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"okami_top_secret_bot/internal/bot"
	"okami_top_secret_bot/internal/config"
	"okami_top_secret_bot/internal/logger"
)

func main() {
	// 1. Конфигурация
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Ошибка загрузки конфигурации: %v", err)
	}

	// 2. Инициализация логгера
	logger := logger.New(cfg.LogLevel, cfg.LogFile)

	// 3. Создание бота
	b, err := bot.New(cfg.TelegramToken, cfg.Debug, logger)
	if err != nil {
		logger.Error("Не удалось создать бота", "error", err)
		os.Exit(1)
	}

	// 4. Контекст для graceful shutdown
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// 5. Запуск бота
	if err := b.Start(ctx); err != nil {
		logger.Error("Бот завершился с ошибкой", "error", err)
		os.Exit(1)
	}
}
