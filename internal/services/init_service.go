package services

import (
	"log"
	"os"

	"github.com/aLieexe/tsukatsuki/internal/config"
	"github.com/aLieexe/tsukatsuki/internal/utils"
	tea "github.com/charmbracelet/bubbletea"
)

type AppConfig struct {
	ProjectName string

	AppDomain string
	AppPort   int

	ServerIP string

	Webserver     string
	Runtime       string
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
	cfg.Project.Domain = setValue(app.AppDomain, "placeholder.com")
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
