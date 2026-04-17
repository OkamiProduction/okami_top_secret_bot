# Telegram Bot on Go 🤖

[![Go Version](https://img.shields.io/badge/Go-1.20+-00ADD8?style=flat&logo=go)](https://go.dev/)
[![Telegram Bot API](https://img.shields.io/badge/Telegram%20Bot%20API-v5.0-blue?style=flat&logo=telegram)](https://core.telegram.org/bots/api)
[![License](https://img.shields.io/badge/license-MIT-green)](./LICENSE)

Простой, но расширяемый Telegram-бот, написанный на Go с использованием библиотеки `go-telegram-bot-api/v5` и конфигурацией через `.env`.

**Основные возможности:**
- Обработка команд `/start` и `/help`
- Лёгкое добавление новых команд
- Работа в режиме Long Polling
- Готов к деплою на VPS с systemd

---

## 📋 Требования

- Go **1.20** или новее
- Токен Telegram бота (получить у [@BotFather](https://t.me/BotFather))
- VPS или локальная машина для запуска (Linux / macOS / Windows)

---

## 🚀 Быстрый старт (локально)

### 1. Клонирование репозитория

```bash
git clone <url-репозитория>
cd my-telegram-bot
```

### 2. Настройка переменных окружения

Скопируйте пример конфигурации:

```bash
cp .env.example .env
```

Отредактируйте `.env`, вставив токен вашего бота:

```env
TELEGRAM_TOKEN=1234567890:ABCdefGHIjklMNOpqrsTUVwxyz
DEBUG=true
```

> 🔒 Файл `.env` добавлен в `.gitignore` и не должен попадать в репозиторий.

### 3. Установка зависимостей

```bash
go mod download
# или
go mod tidy
```

### 4. Запуск

```bash
go run main.go
```

После запуска вы увидите сообщение об успешной авторизации. Откройте Telegram, найдите своего бота и отправьте команду `/start`.

---

## 📦 Сборка бинарного файла

Для создания исполняемого файла:

```bash
go build -o mybot .
```

Для кросс-компиляции под Linux (если вы на Windows/macOS):

```bash
GOOS=linux GOARCH=amd64 go build -o mybot .
```

---

## 🖥️ Деплой на VPS

### 1. Копирование файлов на сервер

```bash
scp mybot .env root@<IP-сервера>:/root/tgbot/
```

### 2. Настройка systemd-сервиса (рекомендуется)

Создайте файл `/etc/systemd/system/tgbot.service`:

```ini
[Unit]
Description=Telegram Bot on Go
After=network.target

[Service]
Type=simple
User=root
WorkingDirectory=/root/tgbot
EnvironmentFile=/root/tgbot/.env
ExecStart=/root/tgbot/mybot
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
```

Активируйте и запустите сервис:

```bash
systemctl daemon-reload
systemctl enable tgbot.service
systemctl start tgbot.service
```

Проверьте статус:

```bash
systemctl status tgbot.service
```

Логи бота:

```bash
journalctl -u tgbot -f
```

---

## ⚙️ Переменные окружения

| Переменная      | Описание                                      | По умолчанию |
|-----------------|-----------------------------------------------|--------------|
| `TELEGRAM_TOKEN`| **Обязательно**. Токен Telegram бота          | —            |
| `DEBUG`         | Включает подробное логирование запросов к API | `false`      |

---

## 📁 Структура проекта

```text
.
├── .env.example       # Шаблон переменных окружения
├── .gitignore         # Исключения для Git
├── go.mod             # Модуль Go
├── go.sum             # Контрольные суммы зависимостей
├── LICENSE            # Лицензия MIT
├── main.go            # Исходный код бота
└── README.md          # Документация
```

---

## 🛠️ Разработка

### Добавление новых команд

Обработка сообщений находится в функции `main()` внутри цикла `for update := range updates`. Чтобы добавить новую команду, расширьте конструкцию `switch`:

```go
switch update.Message.Text {
case "/start":
    replyText = "👋 Привет!"
case "/newcommand":
    replyText = "Вы вызвали новую команду!"
default:
    replyText = "Неизвестная команда. Попробуйте /help"
}
```

### Отладка

Установите `DEBUG=true` в `.env`, чтобы видеть все запросы и ответы к Telegram API в логах.

---

## 📄 Лицензия

Этот проект распространяется под лицензией MIT. Подробнее см. в файле [LICENSE](./LICENSE).

---

## 🤝 Вклад в проект

Если вы хотите предложить улучшения или нашли ошибку:

1. Создайте issue с описанием проблемы или идеи.
2. Создайте форк репозитория и ветку с изменениями.
3. Отправьте pull request.

---

**Приятной разработки!** 🚀