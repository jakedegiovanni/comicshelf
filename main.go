package main

import (
	"context"
	"errors"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func main() {
	cfg := DefaultConfig()
	v := viper.New()

	rootCmd := &cobra.Command{
		Use: "comicshelf",
	}

	rootCmd.AddCommand(Server(&cfg, v))
	rootCmd.AddCommand(NewMarvelCommand())

	rootCmd.PersistentFlags().StringVarP(&cfg.File, "config", "c", cfg.File, "if not present will check for existence of config.yml in current working directory. If none present will be ignored.")

	rootCmd.PersistentFlags().StringVarP(&cfg.Logging.Level, "loglevel", "l", cfg.Logging.Level, "DEBUG|INFO|WARN|ERROR")
	v.BindPFlag("logging.level", rootCmd.PersistentFlags().Lookup("loglevel"))

	rootCmd.PersistentFlags().BoolVarP(&cfg.Logging.Disabled, "silent", "", cfg.Logging.Disabled, "disable logging")
	v.BindPFlag("logging.disabled", rootCmd.PersistentFlags().Lookup("silent"))

	rootCmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		if cfg.File != "" {
			v.SetConfigFile(cfg.File)
		} else {
			v.SetConfigName("config")
			v.SetConfigType("yml")
			v.AddConfigPath(".")
		}

		err := v.ReadInConfig()
		if err != nil {
			var notFound viper.ConfigFileNotFoundError
			if !errors.As(err, &notFound) {
				return fmt.Errorf("could not read config: %w", err)
			}
		}

		var cfg AppConfig
		err = v.Unmarshal(&cfg)
		if err != nil {
			return err
		}

		ctx := cmd.Context()
		ctx = context.WithValue(ctx, ConfigCtxKey, &cfg)
		cmd.SetContext(ctx)
		return nil
	}

	cobra.CheckErr(rootCmd.ExecuteContext(context.Background()))
}
