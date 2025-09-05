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
					{
						Title:       "Nginx",
						Description: "A battle-tested web server",
						Value:       "nginx",
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
						Title:       "CI Only",
						Description: "No need for any kind of configuration",
						Value:       "ci",
					},

					{
						Title:       "CI/CD",
						Description: "Need to setup environtment secret to support CD",
						Value:       "ci_cd",
					},
				},
			},
		},
	}
	return schema
}
