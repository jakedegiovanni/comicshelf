package main

import (
	"time"

	"github.com/jakedegiovanni/comicshelf/comicclient/marvel"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func marvelCmd(v *viper.Viper, svc *marvel.Client) *cobra.Command {

	marvel := &cobra.Command{
		Use: "marvel",
	}

	weekly := &cobra.Command{
		Use: "weekly",
	}

	today := &cobra.Command{
		Use: "today",
		RunE: func(cmd *cobra.Command, args []string) error {
			comics, err := svc.GetWeeklyComics(cmd.Context(), time.Now())
			if err != nil {
				return err
			}

			return prettyPrint(comics)
		},
	}

	weekly.AddCommand(today)

	marvel.AddCommand(weekly)

	return marvel
}
