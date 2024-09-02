package main

import (
	"context"
	_ "embed"
	"errors"
	"io"
	"log/slog"
	"os"

	"github.com/jakedegiovanni/comicshelf/internal/filedb"
	"github.com/jakedegiovanni/comicshelf/internal/server"
	"github.com/jakedegiovanni/comicshelf/marvel"
)

//go:embed default.yaml
var defaultConfig string

type ctxKey struct{ key string }

var cfgCtxKey = &ctxKey{"cfg"}

type config struct {
	File   string        `mapstructure:"file"`
	Marvel marvel.Config `mapstructure:"marvel"`
	Server server.Config `mapstructure:"server"`
	FileDB filedb.Config `mapstructure:"filedb"`
	Logger LoggingConfig `mapstructure:"logger"`
}

type LoggingConfig struct {
	Level    slog.Level `mapstructure:"level"`
	Disabled bool       `mapstructure:"disabled"`
}

func (l LoggingConfig) Writer() io.Writer {
	if l.Disabled {
		return io.Discard
	}

	return os.Stdout
}

func (l LoggingConfig) Slog() *slog.Logger {
	w := l.Writer()

	logger := slog.New(slog.NewTextHandler(w, &slog.HandlerOptions{
		AddSource: true,
		Level:     l.Level,
	}))

	return logger
}

func getConfigFromCtx(ctx context.Context) (*config, error) {
	cfg, ok := ctx.Value(cfgCtxKey).(*config)
	if !ok {
		return nil, errors.New("could not extract config from context")
	}

	return cfg, nil
}
