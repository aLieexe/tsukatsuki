package ui

import "github.com/charmbracelet/lipgloss"

var (
	DebugStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	InfoStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("39"))
	WarnStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("220"))
	ErrorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("196"))
)
