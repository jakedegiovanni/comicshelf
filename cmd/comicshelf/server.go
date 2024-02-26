package main

import (
	"log/slog"

	"github.com/jakedegiovanni/comicshelf/comicclient/marvel"
	"github.com/jakedegiovanni/comicshelf/filedb"
	"github.com/jakedegiovanni/comicshelf/server"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func serverCmd(v *viper.Viper) *cobra.Command {
	// todo viper config

	server := &cobra.Command{
		Use: "server",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := getConfigFromCtx(cmd.Context())
			if err != nil {
				return err
			}

			slog.SetDefault(cfg.Logger.Slog())

			userSvc, err := filedb.New(&cfg.FileDB)
			if err != nil {
				return err
			}
			defer userSvc.Shutdown()

			marvelSvc := marvel.New(&cfg.Marvel)

			svc := server.New(&cfg.Server, marvelSvc, marvelSvc, userSvc)
			return svc.Run(cmd.Context())
		},
	}

	return server
}
