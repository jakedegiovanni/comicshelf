package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
)

var ConfigCtxKey = &CtxKey{Key: "cfg"}

type CtxKey struct {
	Key string
}

type AppConfig struct {
	File    string        `mapstructure:"file"`
	Logging LoggingConfig `mapstructure:"logging"`
	Server  ServerConfig  `mapstructure:"server"`
}

type LoggingConfig struct {
	Level    string `mapstructure:"level"`
	Disabled bool   `mapstructure:"disabled"`
}

func (l LoggingConfig) SlogLevel() (slog.Level, error) {
	var lvl slog.Level
	err := lvl.UnmarshalText([]byte(l.Level))
	return lvl, err
}

func (l LoggingConfig) Writer() io.Writer {
	if l.Disabled {
		return io.Discard
	}

	return os.Stdout
}

func (l LoggingConfig) Logger() (*slog.Logger, error) {
	lvl, err := l.SlogLevel()
	if err != nil {
		return nil, fmt.Errorf("unknown slog level: %w", err)
	}

	w := l.Writer()

	logger := slog.New(slog.NewTextHandler(w, &slog.HandlerOptions{
		AddSource: true,
		Level:     lvl,
	}))

	return logger, nil
}

type ServerConfig struct {
	Address string `mapstructure:"address"`
}

func DefaultConfig() AppConfig {
	return AppConfig{
		Logging: LoggingConfig{
			Level:    "DEBUG",
			Disabled: false,
		},
		Server: ServerConfig{
			Address: "127.0.0.1:8080",
		},
	}
}

func GetConfigFromCtx(ctx context.Context) (*AppConfig, error) {
	cfg, ok := ctx.Value(ConfigCtxKey).(*AppConfig)
	if !ok {
		return nil, errors.New("could not extract config from context")
	}

	return cfg, nil
}
