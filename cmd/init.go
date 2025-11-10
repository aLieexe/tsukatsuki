/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"strconv"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"

	"github.com/aLieexe/tsukatsuki/internal/config"
	"github.com/aLieexe/tsukatsuki/internal/log"
	"github.com/aLieexe/tsukatsuki/internal/prompts"
	"github.com/aLieexe/tsukatsuki/internal/services"
	"github.com/aLieexe/tsukatsuki/internal/ui/multiselect"
	"github.com/aLieexe/tsukatsuki/internal/ui/singleselect"
	"github.com/aLieexe/tsukatsuki/internal/ui/textinput"
	"github.com/aLieexe/tsukatsuki/internal/utils"
)

type UserInput struct {
	AppName        *textinput.Output
	ServerIP       *textinput.Output
	AppSiteAddress *textinput.Output
	AppPort        *textinput.Output
	SetupUser      *textinput.Output
	SSHPort        *textinput.Output

	Webserver *singleselect.Output
	Runtime   *singleselect.Output
	Security  *singleselect.Output

	Services      *multiselect.Output
	GithubActions *multiselect.Output
}

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Create a configuration file",
	Long:  `Create a configuration file named tsukatsuki.yaml, that later can be used to deploy using tsukatsuki deploy`,
	Run: func(cmd *cobra.Command, args []string) {
		runInitCommand(cmd)
	},
}

func runInitCommand(cmd *cobra.Command) {
	logger := log.InitLogger(cmd)

	if config.ConfigFileExist() {
		logger.Warn("Continuing will create a new tsukatsuki.yaml. You can quit by using `Ctrl + C`")
	}

	cfg := services.NewAppConfig()

	outputDir, err := cmd.Flags().GetString("output")
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to read flags: %s", err))
	}

	cfg.OutputDir = outputDir

	userInput := &UserInput{
		AppName:        &textinput.Output{},
		AppPort:        &textinput.Output{},
		ServerIP:       &textinput.Output{},
		AppSiteAddress: &textinput.Output{},
		SetupUser:      &textinput.Output{},
		SSHPort:        &textinput.Output{},

		Webserver: &singleselect.Output{},
		Runtime:   &singleselect.Output{},
		Security:  &singleselect.Output{},

		Services:      &multiselect.Output{},
		GithubActions: &multiselect.Output{},
	}

	selectionSchema := prompts.NewSelectionsSchema()
	questionSchema := prompts.NewQuestionSchema()
	imageMap := services.GetDefaultImageMap()

	// AppName Question
	teaProgram := tea.NewProgram(textinput.InitializeTextinputModel(userInput.AppName, questionSchema.Questions["app-name"], cfg, nil))
	if _, err := teaProgram.Run(); err != nil {
		logger.Error(fmt.Sprintf("error receiving input: %s", err))
		os.Exit(1)
	}

	cfg.ProjectName = userInput.AppName.Value
	cfg.ExitCLI(teaProgram)

	// AppPort Question
	teaProgram = tea.NewProgram(textinput.InitializeTextinputModel(userInput.AppPort, questionSchema.Questions["app-port"], cfg, utils.PortValidator))
	if _, err := teaProgram.Run(); err != nil {
		logger.Error(fmt.Sprintf("error receiving input: %s", err))
		os.Exit(1)
	}

	converted, err := strconv.Atoi(userInput.AppPort.Value)
	if err != nil {
		logger.Warn("port is invalid, defaulted to 6969")
	}
	cfg.AppPort = converted
	cfg.ExitCLI(teaProgram)

	// run time question
	teaProgram = tea.NewProgram(singleselect.InitializeSingleSelectModel(userInput.Runtime, selectionSchema.Questions["runtime"], cfg))
	if _, err := teaProgram.Run(); err != nil {
		logger.Error(fmt.Sprintf("error receiving input: %s", err))
		os.Exit(1)
	}
	cfg.Runtime = userInput.Runtime.Value
	cfg.ExitCLI(teaProgram)

	appImg, exists := imageMap[userInput.Runtime.Value]
	if !exists {
		logger.Error(fmt.Sprintf("failed to map %s", userInput.Runtime.Value))
		os.Exit(1)
	}
	cfg.BuildImage = appImg

	// ServerIP Question
	teaProgram = tea.NewProgram(textinput.InitializeTextinputModel(userInput.ServerIP, questionSchema.Questions["server-ip"], cfg, utils.IpValidator))
	if _, err := teaProgram.Run(); err != nil {
		logger.Error(fmt.Sprintf("error receiving input: %s", err))
		os.Exit(1)
	}

	cfg.ServerIP = userInput.ServerIP.Value
	cfg.ExitCLI(teaProgram)

	// Setup User
	teaProgram = tea.NewProgram(textinput.InitializeTextinputModel(userInput.SetupUser, questionSchema.Questions["server-user"], cfg, nil))
	if _, err := teaProgram.Run(); err != nil {
		logger.Error(fmt.Sprintf("error receiving input: %s", err))
		os.Exit(1)
	}

	cfg.SetupUser = userInput.SetupUser.Value
	cfg.ExitCLI(teaProgram)

	// Security
	teaProgram = tea.NewProgram(singleselect.InitializeSingleSelectModel(userInput.Security, selectionSchema.Questions["security"], cfg))
	if _, err := teaProgram.Run(); err != nil {
		logger.Error(fmt.Sprintf("error receiving input: %s", err))
		os.Exit(1)
	}
	cfg.ExitCLI(teaProgram)

	if userInput.Security.Value == "true" {
		cfg.Security = true
		teaProgram = tea.NewProgram(textinput.InitializeTextinputModel(userInput.SSHPort, questionSchema.Questions["server-port"], cfg, utils.PortValidator))
		if _, err := teaProgram.Run(); err != nil {
			logger.Error(fmt.Sprintf("error receiving input: %s", err))
			os.Exit(1)
		}

		converted, err := strconv.Atoi(userInput.SSHPort.Value)
		if err != nil {
			logger.Warn("port is invalid, defaulted to 22")
		}
		cfg.SSHPort = converted
		cfg.ExitCLI(teaProgram)
	}

	// AppSiteAddress
	teaProgram = tea.NewProgram(textinput.InitializeTextinputModel(userInput.AppSiteAddress, questionSchema.Questions["webserver-endpoint"], cfg, utils.SiteAddressValidator))
	if _, err := teaProgram.Run(); err != nil {
		logger.Error(fmt.Sprintf("error receiving input: %s", err))
		os.Exit(1)
	}

	cfg.AppSiteAddress = userInput.AppSiteAddress.Value
	cfg.ExitCLI(teaProgram)

	// webserver single select question
	teaProgram = tea.NewProgram(singleselect.InitializeSingleSelectModel(userInput.Webserver, selectionSchema.Questions["webserver"], cfg))
	if _, err := teaProgram.Run(); err != nil {
		logger.Error(fmt.Sprintf("error receiving input: %s", err))
		os.Exit(1)
	}
	cfg.Webserver = userInput.Webserver.Value
	cfg.ExitCLI(teaProgram)

	cfg.WebserverImage, exists = imageMap[userInput.Webserver.Value]
	if !exists {
		logger.Error(fmt.Sprintf("failed to map %s", userInput.Runtime.Value))
		os.Exit(1)
	}

	// Services Multi-choice question
	teaProgram = tea.NewProgram(multiselect.InitializeMultiSelectModel(userInput.Services, selectionSchema.Questions["services"], cfg))
	if _, err := teaProgram.Run(); err != nil {
		logger.Error(fmt.Sprintf("error receiving input: %s", err))
		os.Exit(1)
	}
	cfg.ExitCLI(teaProgram)

	for _, svc := range userInput.Services.Value {
		img, exists := imageMap[svc]
		if !exists {
			logger.Error(fmt.Sprintf("failed to map %s", userInput.Runtime.Value))
			os.Exit(1)
		}

		service := services.Service{
			Name:        svc,
			DockerImage: img,
		}
		cfg.Services = append(cfg.Services, service)
	}

	// actions question
	teaProgram = tea.NewProgram(multiselect.InitializeMultiSelectModel(userInput.GithubActions, selectionSchema.Questions["actions"], cfg))
	if _, err := teaProgram.Run(); err != nil {
		logger.Error(fmt.Sprintf("error receiving input: %s", err))
		os.Exit(1)
	}
	for _, actions := range userInput.GithubActions.Value {

		action := services.GithubActions{
			Type: actions,
		}
		cfg.GithubActions = append(cfg.GithubActions, action)
	}

	err = cfg.SaveConfigToFile()
	if err != nil {
		logger.Error(fmt.Sprintf("Failed creating configuration file: %s", err))
	}

	err = cfg.GenerateDeploymentFiles()
	if err != nil {
		logger.Error(fmt.Sprintf("Failed generating deployment files: %s", err))
	}
}

func init() {
	initCmd.Flags().String("output", "deploy", "folder where the configuration will be generated")
	rootCmd.AddCommand(initCmd)
}
