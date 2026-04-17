# 🤖 Okami Top Secret Bot

[![Go Version](https://img.shields.io/badge/Go-1.22+-00ADD8?style=flat&logo=go)](https://go.dev/)
[![Telegram Bot API](https://img.shields.io/badge/Telegram%20Bot%20API-v5.0-blue?style=flat&logo=telegram)](https://core.telegram.org/bots/api)
[![CI](https://github.com/OkamiProduction/okami_top_secret_bot/actions/workflows/ci.yml/badge.svg)](https://github.com/OkamiProduction/okami_top_secret_bot/actions/workflows/ci.yml)
[![Deploy](https://github.com/OkamiProduction/okami_top_secret_bot/actions/workflows/deploy.yml/badge.svg)](https://github.com/OkamiProduction/okami_top_secret_bot/actions/workflows/deploy.yml)
[![License](https://img.shields.io/badge/license-MIT-green)](./LICENSE)

**Telegram-бот для автоматизации контрольных работ**, написанный на Go с соблюдением принципов чистой архитектуры. Поддерживает структурированное логирование с ротацией, полностью автоматический CI/CD, деплой на VPS через systemd и прозрачные уведомления в Telegram.

> ⚠️ **Важно:** Репозиторий сделан публичным **исключительно для включения обязательной защиты веток** (бесплатный тариф GitHub Team не использовался). Код не предназначен для внешнего использования или контрибьютинга.

---

## 📌 Текущее состояние проекта

- ✅ **Основа бота** — обработка команд `/start`, `/help`, легко расширяется.
- ✅ **Чистая архитектура** — разделение на слои (`config`, `logger`, `bot`), независимость от фреймворков.
- ✅ **Структурированное логирование** — `log/slog` + ротация файлов (`lumberjack`).
- ✅ **Деплой на VPS** — systemd-сервис, отдельные пользователи для запуска и CI/CD.
- ✅ **CI/CD (GitHub Actions)**:
  - `ci.yml` — сборка и тесты при каждом PR в `main`.
  - `deploy.yml` — автоматический деплой на staging-сервер при пуше в `main`.
- ✅ **Защита ветки `main`** — обязательные PR и успешный статус `test`.
- ✅ **Уведомления в Telegram** — о результатах CI/CD, ревью и мёрджах.
- 🚧 **В планах**:
  - Мониторинг живости бота (Healthchecks.io / Uptime Robot).
  - Добавление функционала контрольных работ.
  - Production-окружение (отдельный сервер и workflow).

---

## 🚀 Быстрый старт (локальная разработка)

### Требования
- Go 1.22+
- Токен Telegram-бота (получить у [@BotFather](https://t.me/BotFather))

### 1. Клонирование репозитория

```bash
git clone https://github.com/OkamiProduction/okami_top_secret_bot.git
cd okami_top_secret_bot
```

### 2. Настройка переменных окружения

Скопируйте шаблон и заполните `.env`:

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

### 3. Установка зависимостей и запуск

```bash
go mod tidy
go run main.go
```

Бот начнёт принимать сообщения (long polling). Логи выводятся в консоль и в файл `bot.log`.

---

## 📚 Подробная документация

Вся детальная информация вынесена в отдельные документы в папке `docs/`. Рекомендуется изучить их **перед началом работы с проектом**.

| Документ | Содержание |
|----------|------------|
| [📖 Локальная разработка](docs/local-dev.md) | Структура проекта, добавление новых команд, конфигурация |
| [🚀 Деплой и CI/CD](docs/deployment.md) | Настройка сервера, systemd, GitHub Actions, защита веток |
| [📊 Логирование](docs/logging.md) | Конфигурация `slog`, ротация файлов, очистка логов |
| [🔐 Безопасность и пользователи](docs/security.md) | Разделение привилегий, создание пользователей, sudoers, секреты |
| [📨 Уведомления](docs/notifications.md) | Настройка Telegram-бота для уведомлений, типы сообщений |

---

## 🛠️ Используемые технологии

- **Язык:** Go 1.22
- **Telegram API:** [go-telegram-bot-api/v5](https://github.com/go-telegram-bot-api/telegram-bot-api)
- **Логирование:** `log/slog` + [lumberjack](https://github.com/natefinch/lumberjack)
- **Конфигурация:** [godotenv](https://github.com/joho/godotenv)
- **CI/CD:** GitHub Actions (сборка, тесты, деплой по SSH)
- **Инфраструктура:** Linux VPS, systemd, разделение пользователей (`tgbot`, `github-deploy`)

---

## 📄 Лицензия

MIT. См. файл [LICENSE](./LICENSE).