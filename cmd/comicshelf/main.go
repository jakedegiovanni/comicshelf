package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/jakedegiovanni/comicshelf/cmd/internal/hooks"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

func main() {
	var cfgFile string

	rootCmd := &cobra.Command{
		Use: "comicshelf",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			toExpand := defaultConfig

			if cfgFile != "" {
				f, err := os.ReadFile(cfgFile)
				if err != nil {
					return fmt.Errorf("could not read config file: %w", err)
				}

				toExpand = string(f)
			}

			fi := os.Expand(toExpand, func(s string) string {
				sub := strings.SplitN(s, ":", 2)
				v := os.Getenv(sub[0])
				if v == "" && len(sub) > 1 {
					return sub[1]
				}

				return v
			})

			var in map[string]interface{}
			err := yaml.Unmarshal([]byte(fi), &in)
			if err != nil {
				return fmt.Errorf("could not unmarshal expanded values: %w", err)
			}

			var cfg config

			decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
				DecodeHook: mapstructure.ComposeDecodeHookFunc(
					mapstructure.StringToTimeDurationHookFunc(),
					mapstructure.StringToSliceHookFunc(","),
					hooks.UrlHook(),
					hooks.SlogLevelHook(),
				),
				Result: &cfg,
			})

			if err != nil {
				return fmt.Errorf("could not create config decoder: %w", err)
			}

			err = decoder.Decode(in)
			if err != nil {
				return fmt.Errorf("could not decode config file: %w", err)
			}

			ctx := context.WithValue(cmd.Context(), cfgCtxKey, &cfg)
			cmd.SetContext(ctx)
			slog.SetDefault(cfg.Logger.Slog())
			return nil
		},
	}

	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", cfgFile, "if not present will use default.yaml")

	rootCmd.AddCommand(serverCmd())
	rootCmd.AddCommand(marvelCmd())

	cobra.CheckErr(rootCmd.Execute())
}
