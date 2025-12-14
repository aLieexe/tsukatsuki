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
		"go":   "golang:1.24.4-bookworm",
		"node": "node:25.2-bookworm-slim",

		"caddy":   "caddy:2.10.2-alpine",
		"traefik": "traefik:v3.6",

		"postgresql": "postgres:18.0-alpine",
		"redis":      "redis:8.2-alpine3.22",
	}

	return imageMap
}

type AppConfig struct {
	ProjectName string
	AppPort     int
	Runtime     string
	EntryPoint  string
	BuildImage  string

	ServerIP  string
	SetupUser string
	SSHPort   int
	Security  bool
	IPv6      bool

	AppSiteAddress string
	Proxy          string
	ProxyImage     string
	ACMEEmail      string

	Services []Service

	GithubActions []GithubActions

	LocalPath  string
	RemotePath string
	OutputDir  string
	EnvFile    string

	Exit bool
}

func NewAppConfig() *AppConfig {
	cfg := &AppConfig{
		ProjectName: "tsukatsuki",
		AppPort:     5050,
		Runtime:     "go",
		EntryPoint:  utils.GetMainFileLocation("go"),
		BuildImage:  "latest",

		ServerIP:  "127.0.0.1",
		SetupUser: "user1",
		SSHPort:   22,
		Security:  false,
		ACMEEmail: "admin@placeholder.com",

		AppSiteAddress: "placeholder.com",
		Proxy:          "caddy",
		ProxyImage:     "latest",

		Services: nil,

		GithubActions: nil,

		LocalPath: utils.GetAbsolutePath(),
		OutputDir: "deploy",

		Exit: false,
	}

	cfg.EnvFile = ".env"
	cfg.RemotePath = fmt.Sprintf("/home/tsukatsuki/%s", cfg.ProjectName)

	return cfg
}

func (app *AppConfig) SaveConfigToFile() error {
	var cfg config.AppConfigYaml

	cfg.Project.Name = app.ProjectName
	cfg.Project.Port = app.AppPort
	cfg.Project.Runtime = app.Runtime
	cfg.Project.BuildImage = app.BuildImage
	cfg.Project.EntryPoint = app.EntryPoint

	cfg.Server.IP = app.ServerIP
	cfg.Server.SetupUser = app.SetupUser
	cfg.Server.SSHPort = app.SSHPort
	cfg.Server.Security = app.Security

	cfg.Proxy.Domain = app.AppSiteAddress
	cfg.Proxy.Type = app.Proxy
	cfg.Proxy.DockerImage = app.ProxyImage
	cfg.Proxy.ACMEEmail = app.ACMEEmail

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
	cfg.Path.RemotePath = app.RemotePath
	cfg.Path.OutputDir = app.OutputDir
	cfg.Path.EnvFile = app.EnvFile

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

		EntryPoint: yamlConfig.Project.EntryPoint,

		ServerIP:  yamlConfig.Server.IP,
		SetupUser: yamlConfig.Server.SetupUser,
		SSHPort:   yamlConfig.Server.SSHPort,
		Security:  yamlConfig.Server.Security,

		Proxy:          yamlConfig.Proxy.Type,
		AppSiteAddress: yamlConfig.Proxy.Domain,
		ProxyImage:     yamlConfig.Proxy.DockerImage,
		ACMEEmail:      yamlConfig.Proxy.ACMEEmail,

		Services: services,

		GithubActions: githubActions,

		LocalPath:  yamlConfig.Path.LocalPath,
		RemotePath: yamlConfig.Path.RemotePath,
		OutputDir:  yamlConfig.Path.OutputDir,
		EnvFile:    yamlConfig.Path.EnvFile,
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
