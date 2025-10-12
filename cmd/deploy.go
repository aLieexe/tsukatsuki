/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"

	"github.com/aLieexe/tsukatsuki/internal/config"
	"github.com/aLieexe/tsukatsuki/internal/services"
	"github.com/aLieexe/tsukatsuki/internal/ui/textinput"
)

// deployCmd represents the deploy command
var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Create a deployment based on configuration file",
	Long:  `Deploy the application based on the configuration detailed on tsukatsuki.yaml generated via tsukatsuki init`,
	Run: func(cmd *cobra.Command, args []string) {
		if !config.ConfigFileExist() {
			fmt.Println("please generate a config file with tsukatsuki init before deploying")
			os.Exit(1)
		}

		yamlConfig, err := config.GetConfigFromFiles()
		if err != nil {
			fmt.Println(err)
		}

		cfg := services.NewAppConfigFromYaml(yamlConfig)

		fmt.Println("Running server setup with inventory.ini")

		err = services.ExecAnsible(filepath.Join(cfg.OutputDir, "ansible"), "setup.yaml")
		if err != nil {
			fmt.Println(err)
		}

		fmt.Println("In order to continue you must provide us with a user with an admin priviliges")
		var password textinput.Output

		teaProgram := tea.NewProgram(textinput.InitializePasswordInputModel(&password, "what is the root password of your server", "12345678", cfg))
		_, err = teaProgram.Run()
		if err != nil {
			fmt.Println(err)
		}

		err = services.ExecAnsibleWithPassword(filepath.Join(cfg.OutputDir, "ansible"), "setup.yaml", password.Value)
		if err != nil {
			fmt.Println(err)
		}

		err = services.ExecAnsible(filepath.Join(cfg.OutputDir, "ansible"), "deploy.yaml")
		if err != nil {
			fmt.Println(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(deployCmd)
}
