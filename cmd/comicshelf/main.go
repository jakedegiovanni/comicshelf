package main

import (
	"log/slog"
	"net/url"
	"os"
	"time"

	"github.com/jakedegiovanni/comicshelf/comicclient"
	"github.com/jakedegiovanni/comicshelf/comicclient/marvel"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func main() {
	v := viper.New()

	rootCmd := &cobra.Command{
		Use: "comicshelf",
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelDebug,
	}))

	marvelBaseUri, _ := url.Parse("https://gateway.marvel.com/v1/public")
	marvelSvc := marvel.New(&marvel.Config{
		Client: comicclient.Config{
			Timeout: 20 * time.Second,
			BaseURL: *marvelBaseUri,
		},
		DateLayout:    "2006-01-02",
		ReleaseOffset: -3,
	}, logger)

	rootCmd.AddCommand(serverCmd(v))
	rootCmd.AddCommand(marvelCmd(v, marvelSvc))

	cobra.CheckErr(rootCmd.Execute())
}
