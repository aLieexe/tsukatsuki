package singleselect

import (
	"strings"

	"github.com/aLieexe/tsukatsuki/internal/prompts"
	"github.com/aLieexe/tsukatsuki/internal/services"
	tea "github.com/charmbracelet/bubbletea"
)

type Output struct {
	Value string
}

func (o *Output) update(val string) {
	o.Value = val
}

type model struct {
	cursor       int
	promptSchema prompts.SelectionSchema

	output *Output
	exit   *bool
	err    error
}

func InitializeSingleSelectModel(output *Output, promptSchema prompts.SelectionSchema, appConfig *services.AppConfig) model {
	model := model{
		promptSchema: promptSchema,
		output:       output,
		exit:         &appConfig.Exit,
		err:          nil,
		cursor:       0,
	}
	return model
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			*m.exit = true
			return m, tea.Quit

		case "enter":
			m.output.update(m.promptSchema.Options[m.cursor].Title)
			return m, tea.Quit

		case "down", "j":
			m.cursor++
			if m.cursor >= len(m.promptSchema.Options) {
				m.cursor = 0
			}
		case "up", "k":
			m.cursor--
			if m.cursor < 0 {
				m.cursor = len(m.promptSchema.Options) - 1
			}
		}
	}

	return m, nil
}

func (m model) View() string {
	s := strings.Builder{}
	s.WriteString("\n")
	s.WriteString(m.promptSchema.Headers)
	s.WriteString("\n")

	for i := 0; i < len(m.promptSchema.Options); i++ {

		if m.cursor == i {
			s.WriteString("(•) ")
		} else {
			s.WriteString("( ) ")
		}

		s.WriteString(m.promptSchema.Options[i].Title)
		s.WriteString("\n")
		s.WriteString(m.promptSchema.Options[i].Description)
		s.WriteString("\n")
		s.WriteString("\n")

	}
	s.WriteString("(press q to quit, press enter to select)\n")

	return s.String()
}
