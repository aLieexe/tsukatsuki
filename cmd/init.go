/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/aLieexe/tsukatsuki/internal/config"
	"github.com/aLieexe/tsukatsuki/internal/prompts"
	"github.com/aLieexe/tsukatsuki/internal/services"
	"github.com/aLieexe/tsukatsuki/internal/ui/singleselect"
	"github.com/aLieexe/tsukatsuki/internal/ui/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

type UserInput struct {
	AppName       *textinput.Output
	ServerIP      *textinput.Output
	AppDomain     *textinput.Output
	AppPort       *textinput.Output
	Webserver     *singleselect.Output
	Runtime       *singleselect.Output
	GithubActions *singleselect.Output
}

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Create a configuration file",
	Long:  `This should create a configuration file, that later can be used to deploy using tsukatsuki deploy`,
	Run: func(cmd *cobra.Command, args []string) {

		if config.ConfigFileExist() {
			// TODO: should ask if they want to reinitialize, remind the fact that config file already existed, if want to then continue, else quit

		}
		

		// early init part
		userInput := &UserInput{
			AppName:   &textinput.Output{},
			AppPort:   &textinput.Output{},
			ServerIP:  &textinput.Output{},
			AppDomain: &textinput.Output{},

			Webserver:     &singleselect.Output{},
			Runtime:       &singleselect.Output{},
			GithubActions: &singleselect.Output{},
		}

		appConfig := &services.AppConfig{}

		selectionSchema := prompts.InitializeSelectionsSchema()

		// AppName Question
		teaProgram := tea.NewProgram(textinput.InitializeTextinputModel(userInput.AppName, "What is your app name", "tsukatsuki-app", appConfig))
		if _, err := teaProgram.Run(); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		appConfig.ProjectName = userInput.AppName.Value
		appConfig.ExitCLI(teaProgram)

		// TODO: uncommment everything
		// //AppPort Question
		// teaProgram = tea.NewProgram(textinput.InitializeTextinputModel(userInput.AppPort, "What is your app Port", "4000", appConfig))
		// if _, err := teaProgram.Run(); err != nil {
		// 	fmt.Println(err)
		// 	os.Exit(1)
		// }

		// converted, err := strconv.Atoi(userInput.AppPort.Value)
		// if err != nil {
		// 	fmt.Println("Invalid port duh")
		// }
		// appConfig.AppPort = converted
		// appConfig.ExitCLI(teaProgram)

		// //ServerIP Question
		// teaProgram = tea.NewProgram(textinput.InitializeTextinputModel(userInput.ServerIP, "What is your server IP", "127.0.0.1", appConfig))
		// if _, err := teaProgram.Run(); err != nil {
		// 	fmt.Println(err)
		// 	os.Exit(1)
		// }

		// appConfig.ServerIP = userInput.ServerIP.Value
		// appConfig.ExitCLI(teaProgram)

		// // AppDomain
		// teaProgram = tea.NewProgram(textinput.InitializeTextinputModel(userInput.AppDomain, "What is the endpoint that will be used for this App (enter to use ip)", "placeholder.com", appConfig))
		// if _, err := teaProgram.Run(); err != nil {
		// 	fmt.Println(err)
		// 	os.Exit(1)
		// }

		// appConfig.AppDomain = userInput.AppDomain.Value
		// appConfig.ExitCLI(teaProgram)

		// // webserver single select question
		// teaProgram = tea.NewProgram(singleselect.InitializeSingleSelectModel(userInput.Webserver, selectionSchema.Flow["webserver"], appConfig))
		// if _, err := teaProgram.Run(); err != nil {
		// 	fmt.Println(err)
		// 	os.Exit(1)
		// }
		// appConfig.Webserver = userInput.Webserver.Value
		// appConfig.ExitCLI(teaProgram)

		// //run time question
		// teaProgram = tea.NewProgram(singleselect.InitializeSingleSelectModel(userInput.Runtime, selectionSchema.Flow["runtime"], appConfig))
		// if _, err := teaProgram.Run(); err != nil {
		// 	fmt.Println(err)
		// 	os.Exit(1)
		// }
		// appConfig.Runtime = userInput.Runtime.Value
		// appConfig.ExitCLI(teaProgram)

		// actions question
		teaProgram = tea.NewProgram(singleselect.InitializeSingleSelectModel(userInput.GithubActions, selectionSchema.Flow["actions"], appConfig))
		if _, err := teaProgram.Run(); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		appConfig.GithubActions = userInput.GithubActions.Value
		appConfig.ExitCLI(teaProgram)


		
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
