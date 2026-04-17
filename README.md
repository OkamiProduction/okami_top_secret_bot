# 🤖 Okami Top Secret Bot

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://go.dev/)
[![Telegram Bot API](https://img.shields.io/badge/Telegram%20Bot%20API-v5.0-blue?style=flat&logo=telegram)](https://core.telegram.org/bots/api)
[![License](https://img.shields.io/badge/license-MIT-green)](./LICENSE)
[![CI/CD](https://github.com/OkamiProduction/okami_top_secret_bot/actions/workflows/deploy.yml/badge.svg)](https://github.com/OkamiProduction/okami_top_secret_bot/actions/workflows/deploy.yml)

Telegram-бот для автоматизации контрольных работ, написанный на Go с соблюдением принципов чистой архитектуры. Поддерживает структурированное логирование с ротацией, деплой на VPS через systemd и полностью автоматический CI/CD с уведомлениями в Telegram.

---

## 📁 Структура проекта

```text
.
├── .env.example               # Шаблон переменных окружения
├── .gitignore
├── go.mod
├── go.sum
├── main.go                    # Точка входа
├── internal/
│   ├── config/                # Загрузка конфигурации из .env
│   │   └── config.go
│   ├── logger/                # slog + ротация через lumberjack
│   │   └── logger.go
│   └── bot/                   # Основная логика бота
│       └── bot.go
├── deploy/
│   └── tgbot.service          # Systemd unit-файл
├── scripts/
│   └── clean_logs.sh          # Скрипт очистки логов
├── .github/
│   └── workflows/
│       └── deploy.yml         # Автоматический деплой на staging
└── README.md
```

---

## 📋 Требования

- Go **1.21** или новее
- Токен Telegram бота (получить у [@BotFather](https://t.me/BotFather))
- VPS с Linux (для staging/production) или локальная машина для разработки

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

## 🖥️ Настройка сервера (Staging)

### Системные требования
- Linux (Ubuntu/Debian)
- Доступ по SSH с правами root

### 1. Создание пользователей

На сервере должны быть созданы:

- **`tgbot`** — системный пользователь для запуска бота (без права входа).
- **`github-deploy`** — технический пользователь для CI/CD (без интерактивного входа).

Выполните под `root`:

```bash
# Пользователь для бота
useradd -r -s /bin/false -m -d /opt/tgbot tgbot

# Пользователь для GitHub Actions
useradd -r -s /bin/false -m -d /opt/github-deploy github-deploy
usermod -aG tgbot github-deploy
```

### 2. Настройка sudo для `github-deploy`

Создайте файл `/etc/sudoers.d/github-deploy`:

```bash
visudo -f /etc/sudoers.d/github-deploy
```

Содержимое:

```text
github-deploy ALL=(ALL) NOPASSWD: /bin/systemctl restart tgbot, /bin/systemctl status tgbot, /bin/chown tgbot\:tgbot /opt/tgbot/*, /bin/chmod +x /opt/tgbot/tgbot, /bin/chmod 600 /opt/tgbot/.env, /bin/cat > /opt/tgbot/.env
```

### 3. SSH-ключи для `github-deploy`

Сгенерируйте ключи и добавьте публичный в `authorized_keys`:

```bash
mkdir -p /opt/github-deploy/.ssh
chmod 700 /opt/github-deploy/.ssh
ssh-keygen -t ed25519 -C "github-actions-deploy" -f /opt/github-deploy/.ssh/id_ed25519 -N ""
cat /opt/github-deploy/.ssh/id_ed25519.pub >> /opt/github-deploy/.ssh/authorized_keys
chown -R github-deploy:github-deploy /opt/github-deploy/.ssh
chmod 600 /opt/github-deploy/.ssh/authorized_keys
```

**Важно:** Скопируйте **приватный ключ** (содержимое `/opt/github-deploy/.ssh/id_ed25519`) — он потребуется для секрета `STAGING_SSH_PRIVATE_KEY` в GitHub.

### 4. Размещение файлов проекта

Первоначально скопируйте бинарник и `.env` в `/opt/tgbot/`:

```bash
# Локально соберите под Linux
GOOS=linux GOARCH=amd64 go build -o tgbot .

# Скопируйте на сервер
scp tgbot .env root@<IP>:/opt/tgbot/
```

На сервере установите права:

```bash
chown -R tgbot:tgbot /opt/tgbot
chmod +x /opt/tgbot/tgbot
```

### 5. Установка systemd-сервиса

Скопируйте unit-файл из репозитория:

```bash
scp deploy/tgbot.service root@<IP>:/etc/systemd/system/
```

Активируйте и запустите:

```bash
systemctl daemon-reload
systemctl enable tgbot.service
systemctl start tgbot.service
```

### 6. Логирование

Создайте директорию для логов:

```bash
mkdir -p /var/log/tgbot
chown tgbot:tgbot /var/log/tgbot
```

---

## 🔐 Секреты GitHub Actions

Для работы CI/CD добавьте следующие секреты в репозиторий (**Settings → Secrets and variables → Actions**):

| Секрет | Описание |
|--------|----------|
| `STAGING_SSH_HOST` | IP-адрес или домен staging-сервера |
| `STAGING_SSH_USER` | `github-deploy` |
| `STAGING_SSH_PRIVATE_KEY` | **Приватный** SSH-ключ пользователя `github-deploy` |
| `STAGING_BOT_TOKEN` | Токен Telegram-бота для staging |
| `STAGING_BOT_DEBUG` | `true` или `false` |
| `STAGING_BOT_LOG_LEVEL` | `info` / `debug` / `warn` / `error` |
| `TELEGRAM_NOTIFY_BOT_TOKEN` | Токен бота, отправляющего уведомления |
| `TELEGRAM_NOTIFY_CHAT_ID` | ID чата/группы для уведомлений |

---

## 🔄 CI/CD (GitHub Actions)

При каждом пуше в ветку `main` (или `tg-bot` — зависит от настройки) workflow автоматически:

1. Клонирует репозиторий.
2. Устанавливает Go и собирает бинарник под Linux.
3. Копирует бинарник на сервер в `/opt/tgbot/`.
4. Генерирует `.env` файл из секретов.
5. Перезапускает systemd-сервис.
6. Отправляет уведомление в Telegram с именем автора, коммитом и статусом.

Файл workflow: `.github/workflows/deploy.yml`.

---

## 📊 Логирование и мониторинг

### Просмотр логов

```bash
# Логи systemd (текстовый вывод)
journalctl -u tgbot -f

# JSON-логи с ротацией
tail -f /var/log/tgbot/bot.log
```

### Параметры ротации

Встроены в `internal/logger/logger.go`:
- Максимальный размер файла: 10 МБ
- Хранится 5 старых файлов
- Максимальный возраст: 30 дней
- Сжатие gzip

### Очистка логов

Скрипт `scripts/clean_logs.sh` полностью очищает логи **без остановки сервиса**.

```bash
scp scripts/clean_logs.sh root@<IP>:/opt/
ssh root@<IP> "chmod +x /opt/clean_logs.sh && /opt/clean_logs.sh"
```

---

## 👥 Добавление нового разработчика

1. **GitHub:** пригласить в репозиторий (**Settings → Collaborators**).
2. **Telegram:** добавить в группу уведомлений.
3. **(Опционально)** Если нужен прямой SSH-доступ:
   - Создать пользователя на сервере и добавить в группу `tgbot`.
   - Настроить `sudo` по аналогии с `github-deploy`.
   - Добавить его публичный SSH-ключ.

---

## 🛠️ Разработка

### Добавление новых команд

Логика обработки сообщений находится в `internal/bot/bot.go` (метод `handleCommand`). Для добавления команды расширьте `switch`:

```go
case "/newfeature":
    return "✨ Новая возможность!"
```

### Уровни логирования

Управляются через `.env` или секреты CI/CD переменной `LOG_LEVEL`. Допустимые значения: `debug`, `info`, `warn`, `error`.

---

## 📈 Планы по развитию

- [x] Структурированное логирование
- [x] CI/CD с уведомлениями
- [ ] Мониторинг живости (Healthchecks.io)
- [ ] Метрики и алерты
- [ ] Production-окружение

---

## 📄 Лицензия

Проект распространяется под лицензией MIT. См. файл [LICENSE](./LICENSE).

---

**Happy hacking!** 🚀