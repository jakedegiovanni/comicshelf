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
	"github.com/spf13/viper"
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

func defaultConfig(v *viper.Viper) config {
	marvelBaseUri, _ := url.Parse("https://gateway.marvel.com/v1/public")

	v.Set("logger.level", slog.LevelDebug.String())
	v.Set("logger.disabled", false)

	v.Set("filedb.filename", "db.json")

	v.Set("server.address", "127.0.0.1:8080")

	v.Set("marvel.client.timeout", "20s")
	v.Set("marvel.client.base_url", marvelBaseUri.String())
	v.Set("marvel.date_layout", "2006-01-02")
	v.Set("marvel.release_offset", -3)

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

func getConfigFromCtx(ctx context.Context) (*config, error) {
	cfg, ok := ctx.Value(cfgCtxKey).(*config)
	if !ok {
		return nil, errors.New("could not extract config from context")
	}

	return cfg, nil
}
