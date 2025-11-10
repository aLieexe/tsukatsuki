package multiselect

import (
	"slices"
	"strings"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/aLieexe/tsukatsuki/internal/prompts"
	"github.com/aLieexe/tsukatsuki/internal/services"
	"github.com/aLieexe/tsukatsuki/internal/ui"
)

type Output struct {
	Value []string
}

type model struct {
	promptSchema prompts.ChoiceQuestion

	cursor int
	choice []string
	output *Output

	exit *bool
	err  error
}

func InitializeMultiSelectModel(output *Output, promptSchema prompts.ChoiceQuestion, appConfig *services.AppConfig) model {
	model := model{
		promptSchema: promptSchema,

		cursor: 0,
		output: output,

		exit: &appConfig.Exit,
		err:  nil,
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
		case "ctrl+c", "q", "esc":
			return m, tea.Quit

		case "enter":
			m.output.Value = m.choice
			return m, tea.Quit

		case " ":
			currentOptions := m.promptSchema.Choices[m.cursor]
			if slices.Contains(m.choice, currentOptions.Value) {
				for i, choice := range m.choice {
					if choice == currentOptions.Value {
						// removes the element at index i
						m.choice = append(m.choice[:i], m.choice[i+1:]...)
						break
					}
				}
			} else {
				// add to selection
				m.choice = append(m.choice, currentOptions.Value)
			}
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
		currentOptions := m.promptSchema.Choices[i]
		isSelected := slices.Contains(m.choice, currentOptions.Value)
		isHovered := m.cursor == i

		indicator := ui.UnselectedIndicator
		style := ui.NormalItemStyle

		if isSelected && isHovered {
			indicator = ui.SelectedIndicator
			style = ui.HoveredItemStyle
		} else if isSelected {
			indicator = ui.SelectedIndicator
			style = ui.SelectedItemStyle
		} else if isHovered {
			indicator = ui.UnselectedIndicator
			style = ui.HoveredItemStyle
		}

		s.WriteString(style.Render(indicator))
		s.WriteString(" ")
		s.WriteString(style.Render(currentOptions.Title))
		s.WriteString("\n")
		s.WriteString(ui.SubtextStyle.Render("  " + currentOptions.Description))
		s.WriteString("\n\n")
	}

	// hint
	s.WriteString(ui.HintStyle.Render("(↑/k up  ↓/j down  space select  enter confirm  q quit)"))
	s.WriteString("\n")

	return s.String()
}
