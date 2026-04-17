package logger

import (
	"context"
	"log/slog"
	"os"

	"gopkg.in/natefinch/lumberjack.v2"
)

// New создаёт настроенный slog.Logger
func New(level string, logFile string) *slog.Logger {
	var levelVar slog.Level
	switch level {
	case "debug":
		levelVar = slog.LevelDebug
	case "warn":
		levelVar = slog.LevelWarn
	case "error":
		levelVar = slog.LevelError
	default:
		levelVar = slog.LevelInfo
	}

	// Ротация файла
	fileWriter := &lumberjack.Logger{
		Filename:   logFile,
		MaxSize:    10, // MB
		MaxBackups: 5,
		MaxAge:     30, // days
		Compress:   true,
	}

	// Два обработчика: консоль (текст) и файл (JSON)
	consoleHandler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: levelVar})
	fileHandler := slog.NewJSONHandler(fileWriter, &slog.HandlerOptions{Level: levelVar})

	multiHandler := &multiHandler{handlers: []slog.Handler{consoleHandler, fileHandler}}
	return slog.New(multiHandler)
}

// multiHandler пишет в несколько обработчиков
type multiHandler struct {
	handlers []slog.Handler
}

func (m *multiHandler) Enabled(ctx context.Context, level slog.Level) bool {
	for _, h := range m.handlers {
		if h.Enabled(ctx, level) {
			return true
		}
	}
	return false
}

func (m *multiHandler) Handle(ctx context.Context, r slog.Record) error {
	for _, h := range m.handlers {
		if err := h.Handle(ctx, r); err != nil {
			return err
		}
	}
	return nil
}

func (m *multiHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	handlers := make([]slog.Handler, len(m.handlers))
	for i, h := range m.handlers {
		handlers[i] = h.WithAttrs(attrs)
	}
	return &multiHandler{handlers: handlers}
}

func (m *multiHandler) WithGroup(name string) slog.Handler {
	handlers := make([]slog.Handler, len(m.handlers))
	for i, h := range m.handlers {
		handlers[i] = h.WithGroup(name)
	}
	return &multiHandler{handlers: handlers}
}
