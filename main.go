package main

import (
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
)

func main() {
	// 1. Загружаем .env файл
	err := godotenv.Load()
	if err != nil {
		log.Println("⚠️  .env файл не найден, используются системные переменные окружения")
	}

	// 2. Читаем обязательные переменные
	token := os.Getenv("TELEGRAM_TOKEN")
	if token == "" {
		log.Fatal("❌ TELEGRAM_TOKEN не задан ни в .env, ни в системном окружении")
	}

	debug := os.Getenv("DEBUG") == "true"

	// 3. Создаём клиент бота
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic("❌ Не удалось подключиться к Telegram API:", err)
	}

	// 4. Включаем отладку, если нужно
	bot.Debug = debug
	if debug {
		log.Printf("✅ Авторизован как @%s (режим отладки ВКЛ)", bot.Self.UserName)
	} else {
		log.Printf("✅ Бот @%s запущен", bot.Self.UserName)
	}

	// 5. Настройка получения обновлений (Long Polling)
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	// 6. Бесконечный цикл обработки сообщений
	for update := range updates {
		// Пропускаем не-сообщения (например, колбэки от кнопок)
		if update.Message == nil {
			continue
		}

		// Обработка команд
		var replyText string
		switch update.Message.Text {
		case "/start":
			replyText = "👋 Привет! Я бот на Go. Использую .env для конфигурации."
		case "/help":
			replyText = "ℹ️ Я пока умею только отвечать на /start и /help."
		default:
			replyText = "Я не понимаю эту команду. Попробуйте /help"
		}

		// Отправляем ответ
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, replyText)
		if _, err := bot.Send(msg); err != nil {
			log.Printf("❌ Ошибка отправки: %v", err)
		}
	}
}
