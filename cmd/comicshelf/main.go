package main

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use: "comicshelf",
}

var viperCfg = viper.New()

func main() {
	cobra.CheckErr(rootCmd.Execute())
}
