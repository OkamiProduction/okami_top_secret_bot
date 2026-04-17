#!/bin/bash
# Очистка всех логов Telegram-бота без остановки сервиса

LOG_DIR="/var/log/tgbot"
SERVICE_NAME="tgbot"

echo "🧹 Очистка файловых логов в $LOG_DIR"
rm -rf "$LOG_DIR"/*
mkdir -p "$LOG_DIR"

echo "🧹 Очистка journald-логов для сервиса $SERVICE_NAME"
# Удаляем архивные journal-файлы этого юнита (текущий буфер останется, но это мелочи)
journalctl --vacuum-time=1s --unit="$SERVICE_NAME" 2>/dev/null || true

echo "✅ Логи очищены. Сам сервис не трогали."