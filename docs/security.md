# 🔐 Безопасность и пользователи

В этом документе описана модель безопасности сервера, на котором развёрнут бот. Все настройки выполнены с соблюдением принципа наименьших привилегий: каждый компонент системы имеет доступ только к тому, что ему действительно необходимо.

---

## 👥 Пользователи и их роли

На сервере созданы следующие учётные записи:

| Пользователь | UID | Назначение | Возможность входа по SSH | Привилегии |
|--------------|-----|------------|--------------------------|------------|
| `root` | 0 | Администратор сервера | Да (по ключу) | Полный доступ |
| `tgbot` | системный | Запуск сервиса бота | **Нет** (`/bin/false`) | Чтение и исполнение файлов в `/opt/tgbot`. Запись только в `/var/log/tgbot` (через приложение). |
| `github-deploy` | системный | Деплой из GitHub Actions | Только по SSH-ключу | Запись в `/opt/tgbot` (через группу `tgbot`), перезапуск сервиса (`sudo` на ограниченный набор команд). |
| `okami` (опционально) | обычный | Администратор для ручных операций | Да (по ключу) | Полный `sudo` (используется для настройки). |

---

## 🛠️ Создание пользователей (команды для `root`)

### 1. Пользователь `tgbot`

```bash
useradd -r -s /bin/false -m -d /opt/tgbot tgbot
```

- `-r` — системный пользователь.
- `-s /bin/false` — запрет интерактивного входа.
- `-m -d /opt/tgbot` — создание домашней директории `/opt/tgbot`.

### 2. Пользователь `github-deploy`

```bash
useradd -r -s /bin/false -m -d /opt/github-deploy github-deploy
usermod -aG tgbot github-deploy
```

Добавление в группу `tgbot` даёт `github-deploy` права на запись в `/opt/tgbot` (при условии, что директория имеет права `775` и группу `tgbot`).

### 3. Административный пользователь (пример)

```bash
useradd -m -s /bin/bash okami
usermod -aG sudo okami
mkdir -p /home/okami/.ssh
# ... добавить публичный ключ в authorized_keys
```

---

## 📁 Права на файлы и директории

### `/opt/tgbot` (файлы бота)

```bash
chown -R tgbot:tgbot /opt/tgbot
chmod 775 /opt/tgbot
```

> **Важно:** Бинарный файл `/opt/tgbot/tgbot` создаётся и обновляется от имени `tgbot` через CI/CD, поэтому его владельцем всегда остаётся `tgbot:tgbot`. Права `755` устанавливаются автоматически.

- Владелец: `tgbot`
- Группа: `tgbot`
- Права: `rwxrwxr-x` (владелец и группа могут писать, остальные — только читать и исполнять).

**Почему группа имеет запись?** Чтобы пользователь `github-deploy` (член группы `tgbot`) мог копировать новые версии бинарника и создавать `.env`.

### `/var/log/tgbot` (логи)

```bash
mkdir -p /var/log/tgbot
chown tgbot:tgbot /var/log/tgbot
chmod 755 /var/log/tgbot
```

Логи пишутся от имени `tgbot` (приложение работает под этим пользователем), поэтому владелец — `tgbot`. Группе запись не нужна, так как `github-deploy` не пишет логи напрямую.

### `/opt/github-deploy/.ssh` (ключи для CI/CD)

```bash
chmod 700 /opt/github-deploy/.ssh
chmod 600 /opt/github-deploy/.ssh/authorized_keys
chmod 600 /opt/github-deploy/.ssh/id_ed25519
chown -R github-deploy:github-deploy /opt/github-deploy/.ssh
```

---

## 🔑 SSH-ключи

### Для `github-deploy` (используется GitHub Actions)

1. Генерация ключа **на сервере** от имени `github-deploy` (или от `root` с последующей сменой владельца):

   ```bash
   mkdir -p /opt/github-deploy/.ssh
   ssh-keygen -t ed25519 -C "github-actions-deploy" -f /opt/github-deploy/.ssh/id_ed25519 -N ""
   cat /opt/github-deploy/.ssh/id_ed25519.pub >> /opt/github-deploy/.ssh/authorized_keys
   chown -R github-deploy:github-deploy /opt/github-deploy/.ssh
   ```

2. **Приватный ключ** (содержимое `id_ed25519`) копируется в секрет GitHub `STAGING_SSH_PRIVATE_KEY`.

3. Публичный ключ уже находится в `authorized_keys`, что разрешает вход для соответствующего приватного ключа.

### Для администратора (`okami`)

Аналогично создаётся пара ключей (обычно на локальной машине), публичный ключ добавляется в `/home/okami/.ssh/authorized_keys`.

---

## 🛡️ Настройка `sudo` для `github-deploy`

Чтобы `github-deploy` мог перезапускать сервис и изменять права на файлы, ему делегированы конкретные команды через `sudo` **без запроса пароля**.

Файл `/etc/sudoers.d/github-deploy` (редактировать через `visudo -f /etc/sudoers.d/github-deploy`):

```text
github-deploy ALL=(ALL) NOPASSWD: /bin/systemctl stop tgbot, /bin/systemctl restart tgbot, /bin/systemctl status tgbot, /usr/bin/tee /opt/tgbot/.env, /bin/chown tgbot\:tgbot /opt/tgbot/.env, /bin/chmod 600 /opt/tgbot/.env
github-deploy ALL=(tgbot) NOPASSWD: ALL
```

**Пояснение:**
- Первая строка разрешает управление сервисом и запись `.env` от `root`.
- Вторая строка разрешает выполнять **любые** команды от имени пользователя `tgbot` без пароля. Это необходимо для шага деплоя, где используется `sudo -u tgbot bash -c '...'`.

Так как `github-deploy` не имеет интерактивного входа, такое делегирование безопасно.

Убедитесь, что пути к исполняемым файлам совпадают с реальными (проверьте через `which systemctl`, `which tee`).

---

## 🔐 Секреты GitHub Actions

Все чувствительные данные (токены, ключи) хранятся в **Secrets** репозитория и передаются в workflow как переменные окружения.

| Секрет | Описание | Где используется |
|--------|----------|------------------|
| `STAGING_SSH_HOST` | IP-адрес сервера | `deploy.yml` (SSH/SCP) |
| `STAGING_SSH_USER` | `github-deploy` | `deploy.yml` |
| `STAGING_SSH_PRIVATE_KEY` | Приватный ключ пользователя `github-deploy` | `deploy.yml` |
| `STAGING_BOT_TOKEN` | Токен Telegram-бота для staging | Генерация `.env` на сервере |
| `STAGING_BOT_DEBUG` | `true` / `false` | Генерация `.env` |
| `STAGING_BOT_LOG_LEVEL` | Уровень логирования (`info`, `debug`) | Генерация `.env` |
| `TELEGRAM_NOTIFY_BOT_TOKEN` | Токен бота-уведомителя | Все workflow с уведомлениями |
| `TELEGRAM_NOTIFY_CHAT_ID` | ID чата для уведомлений | Все workflow с уведомлениями |

**Важно:** Эти секреты никогда не выводятся в логи GitHub Actions (автоматически маскируются).

---

## 🚫 Что запрещено и почему

- **Прямой вход под `tgbot`** — оболочка `/bin/false` предотвращает интерактивный доступ.
- **Использование `root` в CI/CD** — вместо этого создан `github-deploy` с ограниченными правами.
- **Хранение `.env` в репозитории** — файл добавлен в `.gitignore`, все секреты передаются через CI.
- **Права `777` на директории** — всегда используются минимально необходимые права (`755`, `775`, `700`).

---

## 🔄 Что делать при компрометации ключа

1. **Сгенерировать новую пару ключей** для `github-deploy` (см. раздел выше).
2. **Заменить секрет `STAGING_SSH_PRIVATE_KEY`** в GitHub.
3. **Удалить скомпрометированный публичный ключ** из `/opt/github-deploy/.ssh/authorized_keys`.
4. **Перезапустить упавший workflow** — деплой должен пройти успешно с новым ключом.

Если скомпрометирован токен бота — пересоздать его через @BotFather и обновить секрет `STAGING_BOT_TOKEN`.

---

## 🔗 Связанные документы

- [🚀 Деплой и CI/CD](deployment.md) — как эти пользователи задействованы в автоматизации.
- [📨 Уведомления](notifications.md) — безопасность Telegram-ботов.