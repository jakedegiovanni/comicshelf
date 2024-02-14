package main

import (
	"errors"
	"fmt"

	"github.com/jakedegiovanni/comicshelf/cmd/hooks"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func main() {
	cfg := defaultConfig()
	v := viper.New()

	rootCmd := &cobra.Command{
		Use: "comicshelf",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
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

			var cfg config
			err = v.Unmarshal(&cfg, viper.DecodeHook(
				mapstructure.ComposeDecodeHookFunc(
					mapstructure.StringToTimeDurationHookFunc(),
					mapstructure.StringToSliceHookFunc(","),
					hooks.UrlHook(),
					hooks.SlogLevelHook(),
				),
			))
			if err != nil {
				return err
			}

			putConfigIntoCtx(cmd, &cfg)
			return nil
		},
	}

	rootCmd.PersistentFlags().StringVarP(&cfg.File, "config", "c", cfg.File, "if not present will check for existence of config.yml in current directory. If none present will be ignored.")

	var logLvl string
	rootCmd.PersistentFlags().StringVarP(&logLvl, "loglevel", "l", cfg.Logger.Level.String(), "DEBUG|INFO|WARN|ERROR")
	_ = v.BindPFlag("logger.level", rootCmd.PersistentFlags().Lookup("loglevel"))

	rootCmd.PersistentFlags().BoolVarP(&cfg.Logger.Disabled, "silent", "", cfg.Logger.Disabled, "disable logging")
	_ = v.BindPFlag("logging.disabled", rootCmd.PersistentFlags().Lookup("silent"))

	rootCmd.AddCommand(serverCmd(v))
	rootCmd.AddCommand(marvelCmd(v))

	cobra.CheckErr(rootCmd.Execute())
}
