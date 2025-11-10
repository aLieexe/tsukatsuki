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
	"github.com/aLieexe/tsukatsuki/internal/log"
	"github.com/aLieexe/tsukatsuki/internal/services"
	"github.com/aLieexe/tsukatsuki/internal/ui/textinput"
)

// deployCmd represents the deploy command
var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Create a deployment based on configuration file",
	Long:  `Deploy the application based on the configuration detailed on tsukatsuki.yaml generated via tsukatsuki init`,
	Run: func(cmd *cobra.Command, args []string) {
		logger := log.InitLogger(cmd)

		if !config.ConfigFileExist() {
			logger.Warn("please generate a config file with tsukatsuki init before deploying")
			os.Exit(1)
		}

		yamlConfig, err := config.GetConfigFromFiles()
		if err != nil {
			if err != nil {
				logger.Error(fmt.Sprintf("failed to read config from files: %s", err))
				os.Exit(1)
			}
		}

		app := services.NewAppConfigFromYaml(yamlConfig)

		logger.Info("Attempting to run ansible server setup with inventory.ini")

		port, err := services.ProbeSSH(app.ServerIP, []int{app.SSHPort, 22})
		if err != nil {
			logger.Warn(fmt.Sprintf("Failed to find open SSH Port: %s", err))
		}

		logger.Info(fmt.Sprintf("Using SSH port %d for connection", port))

		err = services.ExecAnsible(logger, filepath.Join(app.OutputDir, "ansible"), "setup.yaml", port)
		if err != nil {
			logger.Warn(fmt.Sprintf("Failed executing using inventory.ini: %s", err))
		}

		logger.Info("In order to continue you must provide us with a user with an admin priviliges")
		var password textinput.Output

		teaProgram := tea.NewProgram(textinput.InitializePasswordInputModel(&password, "what is the root password of your server", "12345678", app))
		_, err = teaProgram.Run()
		if err != nil {
			logger.Error(fmt.Sprintf("failed to receive input %s: ", err))
			os.Exit(1)
		}

		err = services.ExecAnsibleWithPassword(logger, filepath.Join(app.OutputDir, "ansible"), "setup.yaml", password.Value, port)
		if err != nil {
			logger.Error(fmt.Sprintf("failed to execute: %s", err))
			os.Exit(1)
		}

		logger.Info("Deploying application")

		port, err = services.ProbeSSH(app.ServerIP, []int{app.SSHPort, 22})
		if err != nil {
			logger.Warn(fmt.Sprintf("Failed to find open SSH Port: %s", err))
		}

		logger.Info(fmt.Sprintf("Using SSH port %d for connection", port))
		err = services.ExecAnsible(logger, filepath.Join(app.OutputDir, "ansible"), "deploy.yaml", port)
		if err != nil {
			logger.Error(fmt.Sprintf("failed to execute: %s", err))
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(deployCmd)
}
