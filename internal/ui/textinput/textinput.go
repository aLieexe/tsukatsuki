package textinput

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/aLieexe/tsukatsuki/internal/prompts"
	"github.com/aLieexe/tsukatsuki/internal/services"
)

type Output struct {
	Value string
}

func (o *Output) update(val string) {
	o.Value = val
}

type model struct {
	textInput textinput.Model
	err       error
	output    *Output
	exit      *bool
	header    string
	validator func(string) error
}

func InitializeTextinputModel(output *Output, question prompts.Question, appConfig *services.AppConfig, validator func(string) error) model {
	ti := textinput.New()
	ti.Focus() // focus so it’s ready to type
	ti.CharLimit = 100
	ti.Width = 40
	ti.Placeholder = question.Placeholder

	return model{
		textInput: ti,
		err:       nil,
		output:    output,
		exit:      &appConfig.Exit,
		header:    question.Header,
		validator: validator, // set it
	}
}

func InitializePasswordInputModel(output *Output, header string, placeholder string, appConfig *services.AppConfig) model {
	ti := textinput.New()
	ti.Focus() // focus so it’s ready to type
	ti.CharLimit = 100
	ti.Width = 20
	ti.Placeholder = placeholder

	ti.EchoMode = textinput.EchoPassword
	ti.EchoCharacter = '•'

	return model{
		textInput: ti,
		err:       nil,
		output:    output,
		exit:      &appConfig.Exit,
		header:    header,
	}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	// basically this is to check what key is pressed
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			val := m.textInput.Value()
			if m.validator != nil {
				if err := m.validator(val); err != nil {
					m.err = err
					return m, nil
				}
			}
			m.output.update(val)
			return m, tea.Quit

		case tea.KeyCtrlC, tea.KeyEsc:
			*m.exit = true
			return m, tea.Quit
		}
	}

	// if its not a key that we handle that mean we can update it again
	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m model) View() string {
	view := fmt.Sprintf("%s: \n%s\n", m.header, m.textInput.View())
	if m.err != nil {
		view += fmt.Sprintf("error: %s\n", m.err.Error())
	}
	return view
}
