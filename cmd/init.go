/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"strconv"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"

	"github.com/aLieexe/tsukatsuki/internal/config"
	"github.com/aLieexe/tsukatsuki/internal/prompts"
	"github.com/aLieexe/tsukatsuki/internal/services"
	"github.com/aLieexe/tsukatsuki/internal/ui/singleselect"
	"github.com/aLieexe/tsukatsuki/internal/ui/textinput"
	"github.com/aLieexe/tsukatsuki/internal/utils"
)

type UserInput struct {
	AppName        *textinput.Output
	ServerIP       *textinput.Output
	AppSiteAddress *textinput.Output
	AppPort        *textinput.Output
	Branch         *textinput.Output
	SetupUser      *textinput.Output

	Webserver     *singleselect.Output
	Runtime       *singleselect.Output
	GithubActions *singleselect.Output
}

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Create a configuration file",
	Long:  `Create a configuration file named tsukatsuki.yaml, that later can be used to deploy using tsukatsuki deploy`,
	Run: func(cmd *cobra.Command, args []string) {
		if config.ConfigFileExist() {
			// TODO: should ask if they want to reinitialize, remind the fact that config file already existed, if want to then continue, else quit
		}

		cfg := services.NewAppConfig()

		// early init part
		userInput := &UserInput{
			AppName:        &textinput.Output{},
			AppPort:        &textinput.Output{},
			ServerIP:       &textinput.Output{},
			AppSiteAddress: &textinput.Output{},
			Branch:         &textinput.Output{},
			SetupUser:      &textinput.Output{},

			Webserver:     &singleselect.Output{},
			Runtime:       &singleselect.Output{},
			GithubActions: &singleselect.Output{},
		}

		selectionSchema := prompts.InitializeSelectionsSchema()

		// AppName Question
		teaProgram := tea.NewProgram(textinput.InitializeTextinputModel(userInput.AppName, "What is your app name", utils.GetProjectDirectory(), cfg, nil))
		if _, err := teaProgram.Run(); err != nil {
			log.Println(err)
			os.Exit(1)
		}

		cfg.ProjectName = userInput.AppName.Value
		cfg.ExitCLI(teaProgram)

		// AppPort Question
		teaProgram = tea.NewProgram(textinput.InitializeTextinputModel(userInput.AppPort, "What is your app Port", "6969", cfg, utils.PortValidator))
		if _, err := teaProgram.Run(); err != nil {
			log.Println(err)
			os.Exit(1)
		}

		converted, err := strconv.Atoi(userInput.AppPort.Value)
		if err != nil {
			log.Println("port is invalid, defaulted to 6969")
		}
		cfg.AppPort = converted
		cfg.ExitCLI(teaProgram)

		// ServerIP Question
		teaProgram = tea.NewProgram(textinput.InitializeTextinputModel(userInput.ServerIP, "What is your server IP", "127.0.0.1", cfg, utils.IpValidator))
		if _, err := teaProgram.Run(); err != nil {
			log.Println(err)
			os.Exit(1)
		}

		cfg.ServerIP = userInput.ServerIP.Value
		cfg.ExitCLI(teaProgram)

		// Setup User
		teaProgram = tea.NewProgram(textinput.InitializeTextinputModel(userInput.SetupUser, "Please provide a sudo user that is not root", "user1", cfg, nil))
		if _, err := teaProgram.Run(); err != nil {
			log.Println(err)
			os.Exit(1)
		}

		cfg.SetupUser = userInput.SetupUser.Value
		cfg.ExitCLI(teaProgram)

		// AppSiteAddress
		teaProgram = tea.NewProgram(textinput.InitializeTextinputModel(userInput.AppSiteAddress, "What is the endpoint that will be used for this App (enter to use ip)", "placeholder.com", cfg, utils.SiteAddressValidator))
		if _, err := teaProgram.Run(); err != nil {
			log.Println(err)
			os.Exit(1)
		}

		if !strings.HasPrefix(userInput.AppSiteAddress.Value, "http://") && !strings.HasPrefix(userInput.AppSiteAddress.Value, "https://") {
			userInput.AppSiteAddress.Value = "http://" + userInput.AppSiteAddress.Value
		}
		u, _ := url.Parse(userInput.AppSiteAddress.Value)
		cfg.AppSiteAddress = u.Host
		cfg.ExitCLI(teaProgram)

		// webserver single select question
		teaProgram = tea.NewProgram(singleselect.InitializeSingleSelectModel(userInput.Webserver, selectionSchema.Flow["webserver"], cfg))
		if _, err := teaProgram.Run(); err != nil {
			log.Println(err)
			os.Exit(1)
		}
		cfg.Webserver = userInput.Webserver.Value
		cfg.ExitCLI(teaProgram)

		// run time question
		teaProgram = tea.NewProgram(singleselect.InitializeSingleSelectModel(userInput.Runtime, selectionSchema.Flow["runtime"], cfg))
		if _, err := teaProgram.Run(); err != nil {
			log.Println(err)
			os.Exit(1)
		}
		cfg.Runtime = userInput.Runtime.Value
		cfg.ExitCLI(teaProgram)

		// actions question
		teaProgram = tea.NewProgram(singleselect.InitializeSingleSelectModel(userInput.GithubActions, selectionSchema.Flow["actions"], cfg))
		if _, err := teaProgram.Run(); err != nil {
			log.Println(err)
			os.Exit(1)
		}
		cfg.GithubActions = userInput.GithubActions.Value
		cfg.ExitCLI(teaProgram)

		err = cfg.CreateConfigurationFile()
		if err != nil {
			log.Println(err)
		}

		if cfg.GithubActions != "none" {
			teaProgram = tea.NewProgram(textinput.InitializeTextinputModel(userInput.Branch, "What branch do you want to use to trigger Github Actions", "main", cfg, nil))
			if _, err := teaProgram.Run(); err != nil {
				log.Println(err)
				os.Exit(1)
			}

			cfg.Branch = userInput.Branch.Value
			cfg.ExitCLI(teaProgram)
		}

		err = cfg.GenerateDeploymentFiles()
		if err != nil {
			fmt.Println(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
