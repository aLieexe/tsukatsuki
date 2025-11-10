package singleselect

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/aLieexe/tsukatsuki/internal/prompts"
	"github.com/aLieexe/tsukatsuki/internal/services"
	"github.com/aLieexe/tsukatsuki/internal/ui"
)

type Output struct {
	Value string
}

func (o *Output) update(val string) {
	o.Value = val
}

type model struct {
	cursor       int
	promptSchema prompts.ChoiceQuestion

	output *Output
	exit   *bool
	err    error
}

func InitializeSingleSelectModel(output *Output, promptSchema prompts.ChoiceQuestion, appConfig *services.AppConfig) model {
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
			m.output.update(m.promptSchema.Choices[m.cursor].Value)
			return m, tea.Quit

		case "down", "j":
			m.cursor++
			if m.cursor >= len(m.promptSchema.Choices) {
				m.cursor = 0
			}
		case "up", "k":
			m.cursor--
			if m.cursor < 0 {
				m.cursor = len(m.promptSchema.Choices) - 1
			}
		}
	}

	return m, nil
}

func (m model) View() string {
	s := strings.Builder{}

	// header
	s.WriteString(ui.HeaderStyle.Render(m.promptSchema.Headers))
	s.WriteString("\n")

	// description
	s.WriteString(ui.DescriptionStyle.Render(m.promptSchema.Description))
	s.WriteString("\n\n")

	// options
	for i := 0; i < len(m.promptSchema.Choices); i++ {
		choice := m.promptSchema.Choices[i]

		if m.cursor == i {
			s.WriteString(ui.HoveredItemStyle.Render(ui.SelectedIndicator))
			s.WriteString(" ")
			s.WriteString(ui.HighlightStyle.Render(choice.Title))
		} else {
			s.WriteString(ui.UnselectedStyle.Render(ui.UnselectedIndicator))
			s.WriteString(" ")
			s.WriteString(ui.NormalItemStyle.Render(choice.Title))
		}

		s.WriteString("\n")
		s.WriteString(ui.SubtextStyle.Render("  " + choice.Description))
		s.WriteString("\n\n")
	}

	// hint
	s.WriteString(ui.HintStyle.Render("(↑/k up  ↓/j down  enter select  q quit)"))
	s.WriteString("\n")

	return s.String()
}
