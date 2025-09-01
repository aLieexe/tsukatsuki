package services

import (
	"log"
	"os"

	"github.com/aLieexe/tsukatsuki/internal/config"
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

func (app *AppConfig) CreateConfigurationFile() error {

	var cfg config.AppConfigYaml

	cfg.Project.Name = app.ProjectName
	cfg.Project.Domain = app.AppDomain
	cfg.Project.Port = app.AppPort
	cfg.Project.Runtime = app.Runtime

	cfg.Server.IP = app.ServerIP

	cfg.Webserver.Type = app.Webserver

	cfg.GithubActions.Mode = app.GithubActions

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
