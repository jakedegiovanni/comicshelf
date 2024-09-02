package main

import (
	"time"

	"github.com/jakedegiovanni/comicshelf/marvel"
	"github.com/spf13/cobra"
)

func marvelCmd() *cobra.Command {
	// todo - configure strategy?

	today := &cobra.Command{
		Use: "today",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := getConfigFromCtx(cmd.Context())
			if err != nil {
				return err
			}

			svc := marvel.New(&cfg.Marvel)

			comics, err := svc.GetWeeklyComics(cmd.Context(), time.Now())
			if err != nil {
				return err
			}

			return prettyPrint(comics)
		},
	}

	weekly := &cobra.Command{
		Use: "weekly",
	}
	weekly.AddCommand(today)

	marvel := &cobra.Command{
		Use: "marvel",
	}
	marvel.AddCommand(weekly)

	return marvel
}
