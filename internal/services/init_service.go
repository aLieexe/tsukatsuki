package services

import (
	"log"
	"os"

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

func (app *AppConfig) ExitCLI(teaProgram *tea.Program) {
	if app.Exit {
		err := teaProgram.ReleaseTerminal()
		if err != nil {
			log.Fatal(err)
		}

		os.Exit(1)
	}
}
