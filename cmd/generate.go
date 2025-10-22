/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/aLieexe/tsukatsuki/internal/config"
	"github.com/aLieexe/tsukatsuki/internal/log"
	"github.com/aLieexe/tsukatsuki/internal/services"
)

// generateCmd represents the generate command
var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		log := log.InitLogger(cmd)

		yamlConfig, err := config.GetConfigFromFiles()
		if err != nil {
			log.Error(fmt.Sprintf("failed getting configuration files: %s", err))
			log.Warn("Please run `tsukatsuki init` before trying this command again")
			os.Exit(1)
		}

		cfg := services.NewAppConfigFromYaml(yamlConfig)
		err = cfg.GenerateDeploymentFiles()
		if err != nil {
			log.Error(fmt.Sprintf("failed generating configuration files: %s", err))
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(generateCmd)
}
