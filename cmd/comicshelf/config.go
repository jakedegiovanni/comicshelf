package main

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"net/url"
	"os"
	"time"

	"github.com/jakedegiovanni/comicshelf/comicclient"
	"github.com/jakedegiovanni/comicshelf/comicclient/marvel"
	"github.com/jakedegiovanni/comicshelf/filedb"
	"github.com/jakedegiovanni/comicshelf/server"
	"github.com/spf13/cobra"
)

type ctxKey struct{ key string }

var cfgCtxKey = &ctxKey{"cfg"}

type config struct {
	File   string        `mapstructure:"file"`
	Marvel marvel.Config `mapstructure:"marvel"`
	Server server.Config `mapstructure:"server"`
	FileDB filedb.Config `mapstructure:"filedb"`
	Logger LoggingConfig `mapstructure:"logger"`
}

func defaultConfig() config {
	marvelBaseUri, _ := url.Parse("https://gateway.marvel.com/v1/public")

	return config{
		Marvel: marvel.Config{
			Client: comicclient.Config{
				Timeout: 20 * time.Second,
				BaseURL: marvelBaseUri,
			},
			DateLayout:    "2006-01-02",
			ReleaseOffset: -3,
		},
		Server: server.Config{
			Address: "127.0.0.1:8080",
		},
		FileDB: filedb.Config{
			Filename: "db.json",
		},
		Logger: LoggingConfig{
			Level:    slog.LevelDebug,
			Disabled: false,
		},
	}
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

func putConfigIntoCtx(cmd *cobra.Command, cfg *config) {
	ctx := context.WithValue(cmd.Context(), cfgCtxKey, &cfg)
	cmd.SetContext(ctx)
}

func getConfigFromCtx(ctx context.Context) (*config, error) {
	cfg, ok := ctx.Value(cfgCtxKey).(*config)
	if !ok {
		return nil, errors.New("could not extract config from context")
	}

	return cfg, nil
}
