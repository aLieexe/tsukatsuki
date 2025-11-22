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

			"webserver-endpoint": {
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
			"webserver": {
				Headers:     "Webserver Choices",
				Description: "Webserver is a thing",
				Choices: []Choice{
					{
						Title:       "Caddy",
						Description: "A Modern Webserver Written in Golang",
						Value:       "caddy",
					},
				},
			},

			"services": {
				Headers:     "Services choices",
				Description: "You can pick more than 1 btw",
				Choices: []Choice{
					{
						Title:       "Postgresql",
						Description: "A free and open-source relational database management system",
						Value:       "postgresql",
					},

					{
						Title:       "Redis",
						Description: "An in-memory key-value database",
						Value:       "redis",
					},
				},
			},

			"runtime": {
				Headers:     "Runtime to use",
				Description: "What runtime dawg, pick one",
				Choices: []Choice{
					{
						Title:       "Go",
						Description: "Wait why is it only golang here? WDF",
						Value:       "go",
					},
				},
			},

			"actions": {
				Headers:     "Github actions to generate",
				Description: "U can pick more than one",
				Choices: []Choice{
					{
						Title:       "CI",
						Description: "Continous Integrations Github Actions Workflows",
						Value:       "actions-ci",
					},
					{
						Title:       "CD",
						Description: "Continous Deployments automatically pull from github then built it in server",
						Value:       "actions-cd",
					},
				},
			},

			"security": {
				Headers:     "Do you want Server Hardening",
				Description: "Server hardening include SSH Hardening, SELinux, and other security practices",
				Choices: []Choice{
					{
						Title:       "Yes",
						Description: "",
						Value:       "true",
					},
					{
						Title:       "No",
						Description: "",
						Value:       "false",
					},
					// Maybe add Backup, Notifications?
				},
			},
		},
	}
	return schema
}
