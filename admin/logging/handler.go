package logging

import (
	"context"
	"log/slog"
	"os"
)

type customHandler struct {
	baseHandler slog.Handler
}

func (c customHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return c.baseHandler.Enabled(ctx, level)
}

func (c customHandler) Handle(ctx context.Context, record slog.Record) error {
	return c.baseHandler.Handle(ctx, record)
}

func (c customHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return c.baseHandler.WithAttrs(attrs)
}

func (c customHandler) WithGroup(name string) slog.Handler {
	return c.baseHandler.WithGroup(name)
}

func SetupLogging() {
	slogLevel := slog.LevelInfo
	switch os.Getenv("LOG_LEVEL") {
	case "debug":
		slogLevel = slog.LevelDebug
	case "warn":
		slogLevel = slog.LevelWarn
	case "error":
		slogLevel = slog.LevelError
	}

	slog.SetDefault(slog.New(customHandler{
		baseHandler: slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level:     slogLevel,
			AddSource: true,
		}),
	}))
}
