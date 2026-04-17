# 🚀 Деплой и CI/CD

В этом документе описан процесс развёртывания бота на staging-сервере, настройка непрерывной интеграции и доставки (CI/CD) через GitHub Actions, а также правила защиты основной ветки.

---

## 🖥️ Серверная инфраструктура

### Используемые пользователи

| Пользователь | Назначение | Доступ |
|--------------|------------|--------|
| `tgbot` | Запуск сервиса бота | Только чтение/исполнение своих файлов. **Не может** писать в `/opt/tgbot`. |
| `github-deploy` | Доставка кода из GitHub Actions | Может записывать файлы в `/opt/tgbot` (через группу) и перезапускать сервис (`sudo`). |
| `okami` (опционально) | Администратор сервера | Полный `sudo`, используется для настройки и отладки. |

Такое разделение минимизирует ущерб при компрометации одного из компонентов.

### Расположение файлов

- **Бинарник и `.env`**: `/opt/tgbot/`
- **Логи приложения**: `/var/log/tgbot/`
- **Systemd unit**: `/etc/systemd/system/tgbot.service`
- **Домашняя директория `github-deploy`**: `/opt/github-deploy/` (содержит SSH-ключи)

---

## 🔧 Первоначальная настройка сервера (выполняется один раз)

Все команды выполняются от `root`.

### 1. Создание пользователей

```bash
# Пользователь для запуска бота
useradd -r -s /bin/false -m -d /opt/tgbot tgbot

# Технический пользователь для GitHub Actions
useradd -r -s /bin/false -m -d /opt/github-deploy github-deploy

# Добавляем github-deploy в группу tgbot, чтобы он мог писать в /opt/tgbot
usermod -aG tgbot github-deploy
```

### 2. Настройка прав на директории

```bash
# Директория бота
chown -R tgbot:tgbot /opt/tgbot
chmod 775 /opt/tgbot   # группа tgbot имеет запись

# Директория для логов
mkdir -p /var/log/tgbot
chown tgbot:tgbot /var/log/tgbot
chmod 755 /var/log/tgbot
```

### 3. Настройка `sudo` для `github-deploy`

Создайте файл `/etc/sudoers.d/github-deploy`:

```bash
visudo -f /etc/sudoers.d/github-deploy
```

Содержимое (проверьте пути через `which`):

```text
github-deploy ALL=(ALL) NOPASSWD: /bin/systemctl restart tgbot, /bin/systemctl status tgbot, /bin/chown tgbot\:tgbot /opt/tgbot/*, /bin/chmod +x /opt/tgbot/tgbot, /bin/chmod 600 /opt/tgbot/.env, /usr/bin/tee /opt/tgbot/.env
```

### 4. SSH-ключ для GitHub Actions

Сгенерируйте пару ключей **на сервере** от имени `github-deploy`:

```bash
mkdir -p /opt/github-deploy/.ssh
chmod 700 /opt/github-deploy/.ssh
ssh-keygen -t ed25519 -C "github-actions-deploy" -f /opt/github-deploy/.ssh/id_ed25519 -N ""
cat /opt/github-deploy/.ssh/id_ed25519.pub >> /opt/github-deploy/.ssh/authorized_keys
chown -R github-deploy:github-deploy /opt/github-deploy/.ssh
chmod 600 /opt/github-deploy/.ssh/authorized_keys
```

**Важно:** Скопируйте **приватный ключ** (содержимое `/opt/github-deploy/.ssh/id_ed25519`) — он потребуется для секрета `STAGING_SSH_PRIVATE_KEY` в GitHub.

### 5. Systemd-сервис

Скопируйте `deploy/tgbot.service` из репозитория в `/etc/systemd/system/` и выполните:

```bash
systemctl daemon-reload
systemctl enable tgbot.service
systemctl start tgbot.service
```

---

## 🔐 Секреты GitHub Actions

Добавьте в репозиторий (**Settings → Secrets and variables → Actions**):

| Секрет | Описание |
|--------|----------|
| `STAGING_SSH_HOST` | IP-адрес или домен сервера |
| `STAGING_SSH_USER` | `github-deploy` |
| `STAGING_SSH_PRIVATE_KEY` | Приватный SSH-ключ (из шага 4) |
| `STAGING_BOT_TOKEN` | Токен бота (тестового) |
| `STAGING_BOT_DEBUG` | `true` или `false` |
| `STAGING_BOT_LOG_LEVEL` | `info` / `debug` / `warn` / `error` |
| `TELEGRAM_NOTIFY_BOT_TOKEN` | Токен бота-уведомителя |
| `TELEGRAM_NOTIFY_CHAT_ID` | ID чата для уведомлений |

---

## 🤖 Workflow GitHub Actions

### 📄 `.github/workflows/ci.yml` — непрерывная интеграция

- **Триггер:** Pull Request в ветку `main`.
- **Действия:** сборка (`go build`), запуск тестов (`go test`), уведомление в Telegram.

### 📄 `.github/workflows/deploy.yml` — непрерывная доставка

- **Триггер:** Push в ветку `main` (после мёрджа PR).
- **Действия:**
  1. Сборка бинарника под Linux.
  2. Копирование на сервер через SCP.
  3. Генерация `.env` из секретов.
  4. Перезапуск systemd-сервиса.
  5. Уведомление в Telegram.

### 📄 `.github/workflows/notify-merge.yml` — уведомление о мёрдже

- **Триггер:** Закрытие PR с флагом `merged`.
- **Действия:** Отправка сообщения в Telegram с информацией о PR.

### 📄 `.github/workflows/notify-review.yml` — уведомления о ревью

- **Триггер:** Отправка ревью (`submitted`).
- **Действия:** Уведомление об approve, request changes или комментарии.

---

## 🛡️ Защита ветки `main`

В настройках репозитория (**Settings → Branches**) добавлено правило для `main`:

- ✅ **Require a pull request before merging** — прямые пуши запрещены.
- ✅ **Require status checks to pass before merging** — обязателен успех `test` (job из `ci.yml`).
- ✅ **Include administrators** — правило действует на всех.

> ⚠️ Поскольку репозиторий публичный, эти правила **строго enforced**. В приватных репозиториях на бесплатном тарифе они не работают.

---

## 🔄 Типичный процесс разработки

1. Разработчик создаёт feature-ветку и делает изменения.
2. Открывает Pull Request в `main`.
3. GitHub Actions запускает `ci.yml` (сборка, тесты).
4. После успешного прохождения CI ревьюер проверяет код и ставит approve.
5. Разработчик (или ревьюер) нажимает **Merge pull request**.
6. После слияния автоматически запускается `deploy.yml`, который доставляет код на staging-сервер.
7. На каждом этапе в Telegram-чат приходят соответствующие уведомления.

---

## 🧹 Обслуживание

### Обновление бинарника вручную (если CI/CD недоступен)

```bash
# Локально
GOOS=linux GOARCH=amd64 go build -o tgbot .
scp tgbot .env github-deploy@<IP>:/opt/tgbot/
ssh github-deploy@<IP> "sudo systemctl restart tgbot"
```

### Просмотр логов на сервере

```bash
journalctl -u tgbot -f           # текстовые
tail -f /var/log/tgbot/bot.log   # JSON
```

---

## 🔗 Связанные документы

- [🔐 Безопасность и пользователи](security.md) — подробности о правах и sudoers.
- [📨 Уведомления](notifications.md) — настройка Telegram-бота.