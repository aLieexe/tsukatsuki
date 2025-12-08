package prompts

import "github.com/aLieexe/tsukatsuki/internal/utils"

type Choice struct {
	Title       string
	Description string
	Value       string
}

type ChoiceQuestion struct {
	Choices     []Choice
	Headers     string
	Description string
}

type ChoiceQuestionSchema struct {
	Questions map[string]ChoiceQuestion
}

type Question struct {
	Header      string
	Placeholder string
}

type QuestionSchema struct {
	Questions map[string]Question
}

func NewQuestionSchema() *QuestionSchema {
	schema := &QuestionSchema{
		map[string]Question{
			"app-name": {
				Header:      "What is your application name",
				Placeholder: utils.GetProjectDirectory(),
			},

			"app-port": {
				Header:      "In what port is your application running",
				Placeholder: "6969",
			},

			"app-entrypoint": {
				Header:      "Where is your application entrypoint",
				Placeholder: "",
			},

			"server-ip": {
				Header:      "What is your server IP",
				Placeholder: "127.0.0.1",
			},

			"server-user": {
				Header:      "Please provide a sudo user that is not root",
				Placeholder: "user1",
			},

			"server-port": {
				Header:      "What is the custom SSH Port you want to be exposed",
				Placeholder: "222",
			},

			"proxy-endpoint": {
				Header:      "What is the endpoint that will be used for this App (enter to use ip)",
				Placeholder: "subdomain.placeholder.com",
			},
		},
	}

	return schema
}

func NewSelectionsSchema() *ChoiceQuestionSchema {
	schema := &ChoiceQuestionSchema{
		map[string]ChoiceQuestion{
			"proxy": {
				Headers:     "Proxy Choices",
				Description: "Select the proxy to host your application.",
				Choices: []Choice{
					{
						Title:       "Caddy",
						Description: "A modern, automated web server written in Go, with built-in HTTPS.",
						Value:       "caddy",
					},
					{
						Title:       "Traefik",
						Description: "A leading modern open source reverse proxy and ingress controller that makes deploying services and APIs easy.",
						Value:       "traefik",
					},
				},
			},

			"services": {
				Headers:     "Services Choices",
				Description: "Pick one or more services to include in your setup.",
				Choices: []Choice{
					{
						Title:       "Postgresql",
						Description: "A reliable, open-source relational database system.",
						Value:       "postgresql",
					},

					{
						Title:       "Redis",
						Description: "An in-memory key-value store for caching and fast data access.",
						Value:       "redis",
					},
				},
			},

			"runtime": {
				Headers:     "Runtime Environment",
				Description: "Choose the programming runtime for your application.",
				Choices: []Choice{
					{
						Title:       "Go",
						Description: "Use Go (Golang) as your application runtime.",
						Value:       "go",
					},
					{
						Title:       "Node",
						Description: "Use Node JS as your application runtime.",
						Value:       "node",
					},
				},
			},

			"actions": {
				Headers:     "GitHub Actions",
				Description: "Select one or more workflows for CI/CD automation.",
				Choices: []Choice{
					{
						Title:       "CI",
						Description: "Set up continuous integration workflows to automatically test your code.",
						Value:       "actions-ci",
					},
					{
						Title:       "CD",
						Description: "Set up continuous deployment workflows to automatically deploy updates to your server.",
						Value:       "actions-cd",
					},
				},
			},

			"security": {
				Headers:     "Server Hardening",
				Description: "Enable security enhancements like SSH hardening and SELinux for a safer server environment.",
				Choices: []Choice{
					{
						Title:       "Yes",
						Description: "Enable server hardening features.",
						Value:       "true",
					},
					{
						Title:       "No",
						Description: "Do not enable additional security measures.",
						Value:       "false",
					},
					// Maybe add Backup, Notifications?
				},
			},
		},
	}
	return schema
}
