package main

import (
	"time"

	"github.com/jakedegiovanni/comicshelf/marvel"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func marvelCmd(v *viper.Viper) *cobra.Command {

	// todo - configure through viper

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
