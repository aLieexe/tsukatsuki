/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/aLieexe/tsukatsuki/internal/config"
	"github.com/aLieexe/tsukatsuki/internal/services"
	"github.com/aLieexe/tsukatsuki/internal/ui/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

// deployCmd represents the deploy command
var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Create a deployment based on configuration file",
	Long:  `Deploy the application based on the configuration detailed on tsukatsuki.yaml generated via tsukatsuki init`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg := &services.AppConfig{}

		if !config.ConfigFileExist() {
			log.Println("please generate a config file with tsukatsuki init before deploying")
			os.Exit(1)
		}

		// ask root password, or we cooked
		var password textinput.Output

		teaProgram := tea.NewProgram(textinput.InitializePasswordInputModel(&password, "what is the root password of your server", "", cfg))
		_, err := teaProgram.Run()
		if err != nil {
			log.Println(err)
		}

		fmt.Println(password)
	},
}

func init() {
	rootCmd.AddCommand(deployCmd)
}
