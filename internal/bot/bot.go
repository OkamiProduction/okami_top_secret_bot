package bot

import (
	"context"
	"log/slog"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Bot struct {
	api    *tgbotapi.BotAPI
	logger *slog.Logger
}

func New(token string, debug bool, logger *slog.Logger) (*Bot, error) {
	api, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}
	api.Debug = debug
	return &Bot{api: api, logger: logger}, nil
}

func (b *Bot) Start(ctx context.Context) error {
	b.logger.Info("Бот запущен", "username", b.api.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := b.api.GetUpdatesChan(u)

	for {
		select {
		case update := <-updates:
			if update.Message == nil {
				continue
			}
			b.logger.Info("Получено сообщение",
				slog.String("user", update.Message.From.UserName),
				slog.Int64("chat_id", update.Message.Chat.ID),
				slog.String("text", update.Message.Text),
			)
			reply := b.handleCommand(update.Message.Text)
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, reply)
			if _, err := b.api.Send(msg); err != nil {
				b.logger.Error("Ошибка отправки", slog.Any("error", err))
			}
		case <-ctx.Done():
			b.logger.Info("Завершение работы бота")
			b.api.StopReceivingUpdates()
			return nil
		}
	}
}

func (b *Bot) handleCommand(text string) string {
	switch text {
	case "/start":
		return "👋 Привет! Я бот на Go."
	case "/help":
		return "ℹ️ Пока умею только /start и /help."
	default:
		return "Неизвестная команда. Попробуйте /help"
	}
}
