package prompts

type Item struct {
	Title       string
	Description string
	Value       string
}

type SelectionSchema struct {
	Options []Item
	Headers string
}

type Selections struct {
	Flow map[string]SelectionSchema
}

func InitializeSelectionsSchema() *Selections {
	schema := &Selections{
		map[string]SelectionSchema{
			"webserver": {
				Headers: "Webserver Choices",
				Options: []Item{
					{
						Title:       "Caddy",
						Description: "A Modern Webserver Written in Golang",
						Value:       "caddy",
					},
				},
			},

			"services": {
				Headers: "Services choices",
				Options: []Item{
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
				Headers: "What runtime duh",
				Options: []Item{
					{
						Title:       "Go",
						Description: "Wait why is it only golang here? WDF",
						Value:       "go",
					},
				},
			},

			"actions": {
				Headers: "Do you want to setup github actions ",
				Options: []Item{
					{
						Title:       "None",
						Description: "No github workflow file will be generated",
						Value:       "none",
					},
					{
						Title:       "CI Only",
						Description: "No need for any kind of configuration",
						Value:       "actions-ci",
					},
				},
			},
		},
	}
	return schema
}
