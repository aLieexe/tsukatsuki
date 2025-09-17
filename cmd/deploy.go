/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/aLieexe/tsukatsuki/internal/config"
	"github.com/aLieexe/tsukatsuki/internal/services"
	"github.com/spf13/cobra"
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

		res := make([]string, 1)
		res = append(res, cfg.Webserver)
		cfg.GenerateCompose(res, "out")

		// fmt.Println("In order to continue you must provide us with a user with an admin priviliges")
		// // TODO: Guide, make sure it can ssh aswell

		// // ask root password, or we cooked
		// var password textinput.Output
		// var username textinput.Output

		// teaProgram := tea.NewProgram(textinput.InitializeTextinputModel(&username, "Username to use", "user1", cfg, nil))
		// _, err := teaProgram.Run()
		// if err != nil {
		// 	fmt.Println(err)
		// }

		// teaProgram = tea.NewProgram(textinput.InitializePasswordInputModel(&password, "what is the root password of your server", "0812083", cfg))
		// _, err = teaProgram.Run()
		// if err != nil {
		// 	fmt.Println(err)
		// }

		// // TODO: Execute stuff that are generated on init prompts,

	},
}

func init() {
	rootCmd.AddCommand(deployCmd)
}
