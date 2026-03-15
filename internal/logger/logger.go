package logger

import (
	"context"
	"io"
	"log/slog"
	"os"
)

type Config struct {
	Level     Level // Log level (debug, info, warn, error)
	Output    io.Writer
	Format    Format // text or json
	AddSource bool   // include source file:line in logs
}

type Level slog.Level

const (
	LevelDebug = Level(slog.LevelDebug)
	LevelInfo  = Level(slog.LevelInfo)
	LevelWarn  = Level(slog.LevelWarn)
	LevelError = Level(slog.LevelError)
)

type Format string

const (
	FormatText Format = "text"
	FormatJSON Format = "json"
)

var Log *slog.Logger

func Init(verbose bool) {
	cfg := Config{
		Level:  LevelWarn,
		Output: os.Stderr,
		Format: FormatText,
	}
	if verbose {
		cfg.Level = LevelDebug
	}
	Initialize(cfg)
}

func Initialize(cfg Config) {
	if cfg.Output == nil {
		cfg.Output = os.Stderr
	}
	if cfg.Format == "" {
		cfg.Format = FormatText
	}

	opts := &slog.HandlerOptions{
		Level:     slog.Level(cfg.Level),
		AddSource: cfg.AddSource,
	}

	var handler slog.Handler
	switch cfg.Format {
	case FormatJSON:
		handler = slog.NewJSONHandler(cfg.Output, opts)
	default:
		handler = slog.NewTextHandler(cfg.Output, opts)
	}

	Log = slog.New(handler)
}

func WithContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, contextKey{}, Log)
}

func FromContext(ctx context.Context) *slog.Logger {
	if logger, ok := ctx.Value(contextKey{}).(*slog.Logger); ok {
		return logger
	}
	return Log
}

type contextKey struct{}

func Debug(msg string, args ...any) {
	Log.Debug(msg, args...)
}

func Info(msg string, args ...any) {
	Log.Info(msg, args...)
}

func Warn(msg string, args ...any) {
	Log.Warn(msg, args...)
}

func Error(msg string, args ...any) {
	Log.Error(msg, args...)
}
