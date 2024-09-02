package main

import (
	"github.com/jakedegiovanni/comicshelf/internal/filedb"
	"github.com/jakedegiovanni/comicshelf/internal/server"
	"github.com/jakedegiovanni/comicshelf/marvel"
	"github.com/spf13/cobra"
)

func serverCmd() *cobra.Command {
	// todo configure strategy?

	server := &cobra.Command{
		Use: "server",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := getConfigFromCtx(cmd.Context())
			if err != nil {
				return err
			}

			userSvc, err := filedb.New(&cfg.FileDB)
			if err != nil {
				return err
			}
			defer userSvc.Shutdown()

			marvelSvc := marvel.New(&cfg.Marvel)

			svc, err := server.New(&cfg.Server, marvelSvc, marvelSvc, userSvc)
			if err != nil {
				return err
			}

			return svc.Run(cmd.Context())
		},
	}

	return server
}
