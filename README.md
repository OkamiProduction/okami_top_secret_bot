# 🤖 Okami Top Secret Bot

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://go.dev/)
[![Telegram Bot API](https://img.shields.io/badge/Telegram%20Bot%20API-v5.0-blue?style=flat&logo=telegram)](https://core.telegram.org/bots/api)
[![License](https://img.shields.io/badge/license-MIT-green)](./LICENSE)

Telegram-бот для автоматизации контрольных работ, написанный на Go с соблюдением принципов чистой архитектуры. Поддерживает структурированное логирование с ротацией, удобный деплой на VPS и мониторинг.

---

## 📁 Структура проекта

```text
.
├── .env.example          # Пример конфигурации
├── .gitignore
├── go.mod
├── go.sum
├── main.go               # Точка входа
├── internal/
│   ├── config/           # Загрузка конфигурации из .env
│   │   └── config.go
│   ├── logger/           # slog + ротация через lumberjack
│   │   └── logger.go
│   └── bot/              # Основная логика бота
│       └── bot.go
├── scripts/
│   └── clean_logs.sh     # Скрипт очистки логов
└── README.md
```

---

## 📋 Требования

- Go **1.21** или новее
- Токен Telegram бота (получить у [@BotFather](https://t.me/BotFather))
- VPS с Linux (для production) или локальная машина

---

## 🚀 Быстрый старт (локально)

### 1. Клонирование

```bash
git clone https://github.com/OkamiProduction/okami_top_secret_bot.git
cd okami_top_secret_bot
```

### 2. Настройка переменных

Скопируйте шаблон и отредактируйте `.env`:

```bash
cp .env.example .env
```

Пример `.env`:

```env
TELEGRAM_TOKEN=1234567890:ABCdef...
DEBUG=true
LOG_LEVEL=info
LOG_FILE=./bot.log
```

### 3. Установка зависимостей

```bash
go mod tidy
```

### 4. Запуск

```bash
go run main.go
```

---

## 📦 Сборка и деплой на VPS

### Локальная сборка под Linux

```bash
GOOS=linux GOARCH=amd64 go build -o tgbot .
```

### Копирование на сервер

```bash
scp tgbot .env root@<IP>:/root/tgbot/
```

### Настройка systemd сервиса (один раз)

На сервере создайте файл `/etc/systemd/system/tgbot.service`:

```ini
[Unit]
Description=Telegram Bot (tgbot)
After=network.target

[Service]
Type=simple
User=root
WorkingDirectory=/root/tgbot
EnvironmentFile=/root/tgbot/.env
ExecStart=/root/tgbot/tgbot
Restart=always
RestartSec=5
StandardOutput=journal
StandardError=journal

[Install]
WantedBy=multi-user.target
```

Активируйте и запустите:

```bash
systemctl daemon-reload
systemctl enable tgbot.service
systemctl start tgbot.service
```

### Обновление бота

```bash
# Локально
GOOS=linux GOARCH=amd64 go build -o tgbot .
scp tgbot root@<IP>:/root/tgbot/

# На сервере
systemctl restart tgbot
```

---

## 📊 Логирование

Проект использует стандартный пакет `log/slog` с двумя обработчиками:
- **Консоль** (текстовый формат) — попадает в `journald`, виден через `journalctl`.
- **Файл** (JSON) — записывается по пути `LOG_FILE` с автоматической ротацией.

### Параметры ротации (в коде `internal/logger/logger.go`)

- **MaxSize:** 10 МБ
- **MaxBackups:** 5 файлов
- **MaxAge:** 30 дней
- **Compress:** true (gzip)

Общий объём логов не превышает ~50 МБ.

### Просмотр логов

```bash
# journald (текстовые логи)
journalctl -u tgbot -f

# Файловые JSON-логи
tail -f /var/log/tgbot/bot.log
```

### Очистка логов

В папке `scripts/` лежит `clean_logs.sh` для полной очистки логов **без остановки сервиса**.

```bash
#!/bin/bash
# Очистка всех логов Telegram-бота без остановки сервиса

LOG_DIR="/var/log/tgbot"
SERVICE_NAME="tgbot"

echo "🧹 Очистка файловых логов в $LOG_DIR"
rm -rf "$LOG_DIR"/*
mkdir -p "$LOG_DIR"

echo "🧹 Очистка journald-логов для сервиса $SERVICE_NAME"
journalctl --vacuum-time=1s --unit="$SERVICE_NAME" 2>/dev/null || true

echo "✅ Логи очищены. Сам сервис не трогали."
```

Запуск на сервере:

```bash
chmod +x scripts/clean_logs.sh
./scripts/clean_logs.sh
```

---

## 🔄 Управление сервисом

| Команда | Действие |
|---------|----------|
| `systemctl start tgbot` | Запустить |
| `systemctl stop tgbot` | Остановить |
| `systemctl restart tgbot` | Перезапустить |
| `systemctl status tgbot` | Проверить статус |
| `journalctl -u tgbot -f` | Смотреть логи |

---

## 🛠️ Разработка

### Добавление новых команд

Логика обработки сообщений находится в `internal/bot/bot.go` в методе `handleCommand`. Для добавления новой команды расширьте `switch`:

```go
func (b *Bot) handleCommand(text string) string {
    switch text {
    case "/start":
        return "👋 Привет!"
    case "/newfeature":
        return "✨ Новая возможность!"
    default:
        return "Неизвестная команда"
    }
}
```

### Уровни логирования

Управляются через `.env` переменной `LOG_LEVEL`. Допустимые значения: `debug`, `info`, `warn`, `error`.

---

## 📄 Лицензия

Проект распространяется под лицензией MIT. См. файл [LICENSE](./LICENSE).

---

**Happy hacking!** 🚀