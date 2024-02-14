package main

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func serverCmd(v *viper.Viper) *cobra.Command {
	server := &cobra.Command{
		Use: "server",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("server to be implemented")
			return nil
		},
	}

	return server
}
