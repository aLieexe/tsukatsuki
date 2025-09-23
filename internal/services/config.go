package services

import (
	"fmt"
	"log"
	"os"

	"github.com/aLieexe/tsukatsuki/internal/config"
	"github.com/aLieexe/tsukatsuki/internal/utils"
	tea "github.com/charmbracelet/bubbletea"
)

type AppConfig struct {
	ProjectName string
	AppPort     int
	Runtime     string
	MainPath    string

	ServerIP string

	AppSiteAddress string
	Webserver      string
	SSLEmail       string

	Branch        string
	GithubActions string

	LocalPath  string
	RemotePath string

	Exit bool
}

func setValue[T comparable](value, defaultValue T) T {
	var zero T
	if value != zero {
		return value
	}
	return defaultValue
}

func NewAppConfig() *AppConfig {
	cfg := &AppConfig{
		ProjectName: "tsukatsuki",
		AppPort:     5050,
		Runtime:     "go",
		MainPath:    utils.GetMainFileLocation(),

		ServerIP: "127.0.0.1",

		AppSiteAddress: "placeholder.com",
		Webserver:      "caddy",
		SSLEmail:       "hello@gmail.com",

		Branch:        "main",
		GithubActions: "none",

		Exit: false,
	}

	return cfg
}

func (app *AppConfig) CreateConfigurationFile() error {
	var cfg config.AppConfigYaml

	cfg.Project.Name = setValue(app.ProjectName, utils.GetProjectDirectory())
	cfg.Project.Port = setValue(app.AppPort, 6969)
	cfg.Project.Runtime = setValue(app.Runtime, "go")

	cfg.Server.IP = setValue(app.ServerIP, "127.0.0.1")

	cfg.Webserver.Domain = setValue(app.AppSiteAddress, "placeholder.com")
	cfg.Webserver.Type = setValue(app.Webserver, "caddy")
	cfg.Webserver.SSLEmail = setValue(app.SSLEmail, "hello@gmail.com")

	cfg.GithubActions.Mode = setValue(app.GithubActions, "none")
	cfg.GithubActions.Branch = setValue(app.Branch, "main")

	cfg.Path.LocalPath = utils.GetAbsolutePath()
	cfg.Path.RemotePath = fmt.Sprintf("/home/tsukatsuki/%s", cfg.Project.Name)

	return config.CreateConfigFiles(cfg)
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

func NewAppConfigFromYaml(yamlConfig config.AppConfigYaml) *AppConfig {
	cfg := &AppConfig{
		ProjectName: yamlConfig.Project.Name,
		AppPort:     yamlConfig.Project.Port,
		Runtime:     yamlConfig.Project.Runtime,
		MainPath:    utils.GetMainFileLocation(),

		ServerIP: yamlConfig.Server.IP,

		Webserver:      yamlConfig.Webserver.Type,
		SSLEmail:       yamlConfig.Webserver.SSLEmail,
		AppSiteAddress: yamlConfig.Webserver.Domain,

		LocalPath:  yamlConfig.Path.LocalPath,
		RemotePath: yamlConfig.Path.RemotePath,

		Branch:        yamlConfig.GithubActions.Branch,
		GithubActions: yamlConfig.GithubActions.Mode,
	}
	return cfg
}
