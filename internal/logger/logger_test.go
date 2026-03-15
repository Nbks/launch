package logger

import (
	"context"
	"log/slog"
	"testing"
)

func TestInit(t *testing.T) {
	Init(false)
	if Log == nil {
		t.Error("Expected Log to be initialized")
	}

	Init(true)
	if Log == nil {
		t.Error("Expected Log to be initialized with verbose")
	}
}

func TestInitialize(t *testing.T) {
	tests := []struct {
		name string
		cfg  Config
	}{
		{
			name: "default config",
			cfg:  Config{},
		},
		{
			name: "debug level",
			cfg:  Config{Level: LevelDebug},
		},
		{
			name: "info level",
			cfg:  Config{Level: LevelInfo},
		},
		{
			name: "warn level",
			cfg:  Config{Level: LevelWarn},
		},
		{
			name: "error level",
			cfg:  Config{Level: LevelError},
		},
		{
			name: "json format",
			cfg:  Config{Format: FormatJSON},
		},
		{
			name: "text format",
			cfg:  Config{Format: FormatText},
		},
		{
			name: "with source",
			cfg:  Config{AddSource: true},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Initialize(tt.cfg)
			if Log == nil {
				t.Error("Expected Log to be initialized")
			}
		})
	}
}

func TestWithContextAndFromContext(t *testing.T) {
	Initialize(Config{Level: LevelWarn})

	ctx := context.Background()

	ctxWithLogger := WithContext(ctx)
	if ctxWithLogger == nil {
		t.Error("Expected context to be returned")
	}

	logger := FromContext(ctxWithLogger)
	if logger == nil {
		t.Error("Expected logger from context")
	}
}

func TestFromContextWithoutLogger(t *testing.T) {
	Initialize(Config{Level: LevelWarn})

	ctx := context.Background()
	logger := FromContext(ctx)
	if logger == nil {
		t.Error("Expected fallback to default logger")
	}
}

func TestLevelConstants(t *testing.T) {
	tests := []struct {
		name string
		got  Level
		want Level
	}{
		{"LevelDebug", LevelDebug, Level(slog.LevelDebug)},
		{"LevelInfo", LevelInfo, Level(slog.LevelInfo)},
		{"LevelWarn", LevelWarn, Level(slog.LevelWarn)},
		{"LevelError", LevelError, Level(slog.LevelError)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.want {
				t.Errorf("got %v, want %v", tt.got, tt.want)
			}
		})
	}
}

func TestFormatConstants(t *testing.T) {
	if FormatText != "text" {
		t.Errorf("FormatText = %q, want %q", FormatText, "text")
	}
	if FormatJSON != "json" {
		t.Errorf("FormatJSON = %q, want %q", FormatJSON, "json")
	}
}
