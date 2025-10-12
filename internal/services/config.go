package services

import (
	"fmt"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/aLieexe/tsukatsuki/internal/config"
	"github.com/aLieexe/tsukatsuki/internal/utils"
)

type AppConfig struct {
	ProjectName string
	AppPort     int
	Runtime     string
	MainPath    string

	ServerIP  string
	SetupUser string

	AppSiteAddress string
	Webserver      string

	Branch        string
	GithubActions string

	LocalPath  string
	RemotePath string
	OutputDir  string

	Exit bool
}

func NewAppConfig() *AppConfig {
	cfg := &AppConfig{
		ProjectName: "tsukatsuki",
		AppPort:     5050,
		Runtime:     "go",
		MainPath:    utils.GetMainFileLocation(),

		ServerIP:  "127.0.0.1",
		SetupUser: "user1",

		AppSiteAddress: "placeholder.com",
		Webserver:      "caddy",

		Branch:        "main",
		GithubActions: "none",

		LocalPath: utils.GetAbsolutePath(),
		OutputDir: "deploy",

		Exit: false,
	}

	cfg.RemotePath = fmt.Sprintf("/home/tsukatsuki/%s", cfg.ProjectName)

	return cfg
}

func (app *AppConfig) CreateConfigurationFile() error {
	var cfg config.AppConfigYaml

	cfg.Project.Name = app.ProjectName
	cfg.Project.Port = app.AppPort
	cfg.Project.Runtime = app.Runtime

	cfg.Server.IP = app.ServerIP
	cfg.Server.SetupUser = app.SetupUser

	cfg.Webserver.Domain = app.AppSiteAddress
	cfg.Webserver.Type = app.Webserver

	cfg.GithubActions.Mode = app.GithubActions
	cfg.GithubActions.Branch = app.Branch

	cfg.Path.LocalPath = app.LocalPath
	cfg.Path.RemotePath = app.ProjectName
	cfg.Path.OutputDir = app.OutputDir

	return config.CreateConfigFiles(cfg)
}

func NewAppConfigFromYaml(yamlConfig config.AppConfigYaml) *AppConfig {
	cfg := &AppConfig{
		ProjectName: yamlConfig.Project.Name,
		AppPort:     yamlConfig.Project.Port,
		Runtime:     yamlConfig.Project.Runtime,
		MainPath:    utils.GetMainFileLocation(),

		ServerIP:  yamlConfig.Server.IP,
		SetupUser: yamlConfig.Server.SetupUser,

		Webserver:      yamlConfig.Webserver.Type,
		AppSiteAddress: yamlConfig.Webserver.Domain,

		LocalPath:  yamlConfig.Path.LocalPath,
		RemotePath: yamlConfig.Path.RemotePath,
		OutputDir:  yamlConfig.Path.OutputDir,

		Branch:        yamlConfig.GithubActions.Branch,
		GithubActions: yamlConfig.GithubActions.Mode,
	}
	return cfg
}

func (app *AppConfig) ExitCLI(teaProgram *tea.Program) {
	if app.Exit {
		err := teaProgram.ReleaseTerminal()
		if err != nil {
			log.Fatal(err)
		}

		os.Exit(1)
	}
}
