package services

import (
	"log"
	"os"

	"github.com/aLieexe/tsukatsuki/internal/config"
	"github.com/aLieexe/tsukatsuki/internal/utils"
	tea "github.com/charmbracelet/bubbletea"
)

type AppConfig struct {
	ProjectName    string
	AppSiteAddress string
	AppPort        int
	Runtime        string
	MainPath       string

	ServerIP string

	Webserver string
	SSLEmail  string

	Branch        string
	GithubActions string

	Exit bool
}

func setValue[T comparable](value, defaultValue T) T {
	var zero T
	if value != zero {
		return value
	}
	return defaultValue
}

func (app *AppConfig) CreateConfigurationFile() error {
	var cfg config.AppConfigYaml

	cfg.Project.Name = setValue(app.ProjectName, utils.GetProjectDirectory())
	cfg.Project.Domain = setValue(app.AppSiteAddress, "placeholder.com")
	cfg.Project.Port = setValue(app.AppPort, 6969)
	cfg.Project.Runtime = setValue(app.Runtime, "go")

	cfg.Server.IP = setValue(app.ServerIP, "127.0.0.1")

	cfg.Webserver.Type = setValue(app.Webserver, "caddy")

	cfg.GithubActions.Mode = setValue(app.GithubActions, "ci")

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
		ProjectName:    yamlConfig.Project.Name,
		AppSiteAddress: yamlConfig.Project.Domain,
		AppPort:        yamlConfig.Project.Port,
		Runtime:        yamlConfig.Project.Runtime,
		MainPath:       utils.GetMainFileLocation(),

		ServerIP: yamlConfig.Server.IP,

		Webserver: yamlConfig.Webserver.Type,
		SSLEmail:  yamlConfig.Webserver.SSLEmail,

		Branch:        yamlConfig.GithubActions.Branch,
		GithubActions: yamlConfig.GithubActions.Mode,
	}
	return cfg
}
