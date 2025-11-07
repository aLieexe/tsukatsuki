package services

import (
	"fmt"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/aLieexe/tsukatsuki/internal/config"
	"github.com/aLieexe/tsukatsuki/internal/utils"
)

type Service struct {
	Name        string
	DockerImage string
}

type GithubActions struct {
	Type string
}

func GetDefaultImageMap() map[string]string {
	imageMap := map[string]string{
		"go": "golang:1.24.4-bookworm",

		"caddy": "caddy:2.10.2-alpine",

		"postgresql": "postgres:18.0-alpine",
		"redis":      "redis:8.2-alpine3.22",
	}

	return imageMap
}

type AppConfig struct {
	ProjectName string
	AppPort     int
	Runtime     string
	MainPath    string
	BuildImage  string

	ServerIP  string
	SetupUser string
	SSHPort   string

	AppSiteAddress string
	Webserver      string
	WebserverImage string

	Services []Service

	GithubActions []GithubActions

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
		BuildImage:  "latest",

		ServerIP:  "127.0.0.1",
		SetupUser: "user1",

		AppSiteAddress: "placeholder.com",
		Webserver:      "caddy",
		WebserverImage: "latest",

		Services: nil,

		GithubActions: nil,

		LocalPath: utils.GetAbsolutePath(),
		OutputDir: "deploy",

		Exit: false,
	}

	cfg.RemotePath = fmt.Sprintf("/home/tsukatsuki/%s", cfg.ProjectName)

	return cfg
}

func (app *AppConfig) SaveConfigToFile() error {
	var cfg config.AppConfigYaml

	cfg.Project.Name = app.ProjectName
	cfg.Project.Port = app.AppPort
	cfg.Project.Runtime = app.Runtime
	cfg.Project.BuildImage = app.BuildImage

	cfg.Server.IP = app.ServerIP
	cfg.Server.SetupUser = app.SetupUser

	cfg.Webserver.Domain = app.AppSiteAddress
	cfg.Webserver.Type = app.Webserver
	cfg.Webserver.DockerImage = app.WebserverImage

	for _, service := range app.Services {
		cfg.Services = append(cfg.Services, struct {
			Name        string `yaml:"name"`
			DockerImage string `yaml:"docker_image"`
		}{
			Name:        service.Name,
			DockerImage: service.DockerImage,
		})
	}

	for _, actions := range app.GithubActions {
		cfg.GithubActions = append(cfg.GithubActions, struct {
			Type string `yaml:"type"`
		}{
			Type: actions.Type,
		})
	}

	cfg.Path.LocalPath = app.LocalPath
	cfg.Path.RemotePath = fmt.Sprintf("/home/tsukatsuki/%s", app.ProjectName)
	cfg.Path.OutputDir = app.OutputDir

	return config.UpdateConfigFile(cfg)
}

func NewAppConfigFromYaml(yamlConfig config.AppConfigYaml) *AppConfig {
	var services []Service
	for _, yamlService := range yamlConfig.Services {
		services = append(services, Service{
			Name:        yamlService.Name,
			DockerImage: yamlService.DockerImage,
		})
	}

	var githubActions []GithubActions
	for _, actions := range yamlConfig.GithubActions {
		githubActions = append(githubActions, GithubActions{
			Type: actions.Type,
		})
	}

	cfg := &AppConfig{
		ProjectName: yamlConfig.Project.Name,
		AppPort:     yamlConfig.Project.Port,
		Runtime:     yamlConfig.Project.Runtime,
		BuildImage:  yamlConfig.Project.BuildImage,

		MainPath: utils.GetMainFileLocation(),

		ServerIP:  yamlConfig.Server.IP,
		SetupUser: yamlConfig.Server.SetupUser,

		Webserver:      yamlConfig.Webserver.Type,
		AppSiteAddress: yamlConfig.Webserver.Domain,
		WebserverImage: yamlConfig.Webserver.DockerImage,

		Services: services,

		GithubActions: githubActions,

		LocalPath:  yamlConfig.Path.LocalPath,
		RemotePath: yamlConfig.Path.RemotePath,
		OutputDir:  yamlConfig.Path.OutputDir,
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
