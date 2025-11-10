package ui

import "github.com/charmbracelet/lipgloss"

// color palette
var (
	colorPrimary   = lipgloss.Color("39")  // cyan
	colorSecondary = lipgloss.Color("135") // purple
	colorAccent    = lipgloss.Color("220") // yellow
	colorSuccess   = lipgloss.Color("34")  // green
	colorWarning   = lipgloss.Color("220") // yellow
	colorError     = lipgloss.Color("196") // red
	colorDebug     = lipgloss.Color("240") // gray
	colorBg        = lipgloss.Color("235") // dark background
	colorBorder    = lipgloss.Color("59")  // dark gray
	colorText      = lipgloss.Color("15")  // white
	colorMuted     = lipgloss.Color("244") // gray text
)

// log styles
var (
	DebugStyle = lipgloss.NewStyle().Foreground(colorDebug)
	InfoStyle  = lipgloss.NewStyle().Foreground(colorPrimary).Bold(true)
	WarnStyle  = lipgloss.NewStyle().Foreground(colorWarning).Bold(true)
	ErrorStyle = lipgloss.NewStyle().Foreground(colorError).Bold(true)
)

// text styles
var (
	TitleStyle = lipgloss.NewStyle().
			Foreground(colorText).
			Bold(true).
			MarginBottom(1)

	HeaderStyle = lipgloss.NewStyle().
			Foreground(colorPrimary).
			Bold(true).
			MarginTop(1).
			MarginBottom(1)

	DescriptionStyle = lipgloss.NewStyle().
				Foreground(colorMuted).
				Italic(true).
				MarginBottom(1)

	SubtextStyle = lipgloss.NewStyle().
			Foreground(colorMuted).
			Italic(true)

	HighlightStyle = lipgloss.NewStyle().
			Foreground(colorAccent).
			Bold(true)
)

// input/selection styles
var (
	CursorStyle = lipgloss.NewStyle().
			Foreground(colorAccent).
			Bold(true)

	SelectionStyle = lipgloss.NewStyle().
			Foreground(colorSuccess).
			Bold(true)

	UnselectedStyle = lipgloss.NewStyle().
			Foreground(colorMuted)

	InputStyle = lipgloss.NewStyle().
			Foreground(colorText).
			Background(colorBg).
			Padding(0, 1).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(colorBorder)

	InputErrorStyle = lipgloss.NewStyle().
			Foreground(colorError).
			Background(colorBg).
			Padding(0, 1).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(colorError)

	FocusedInputStyle = lipgloss.NewStyle().
				Foreground(colorText).
				Background(colorBg).
				Padding(0, 1).
				Border(lipgloss.RoundedBorder()).
				BorderForeground(colorPrimary).
				Bold(true)
)

// box styles
var (
	BoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(colorBorder).
			Padding(1).
			Margin(1)

	PrimaryBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(colorPrimary).
			Padding(1).
			Margin(1).
			Bold(true)

	ErrorBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(colorError).
			Padding(1).
			Margin(1)

	SuccessBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(colorSuccess).
			Padding(1).
			Margin(1)
)

// button/action styles
var (
	ButtonStyle = lipgloss.NewStyle().
			Foreground(colorText).
			Background(colorPrimary).
			Padding(0, 2).
			Bold(true).
			Align(lipgloss.Center)

	ButtonInactiveStyle = lipgloss.NewStyle().
				Foreground(colorMuted).
				Background(colorBorder).
				Padding(0, 2).
				Align(lipgloss.Center)
)

// helper text styles
var (
	HintStyle = lipgloss.NewStyle().
			Foreground(colorMuted).
			Italic(true).
			MarginTop(1)

	KeybindStyle = lipgloss.NewStyle().
			Foreground(colorAccent).
			Bold(true)
)

// status/indicator styles
var (
	SelectedIndicator   = "●" // filled circle
	UnselectedIndicator = "○" // empty circle

	SelectedItemStyle = lipgloss.NewStyle().
				Foreground(colorSuccess).
				Bold(true)

	HoveredItemStyle = lipgloss.NewStyle().
				Foreground(colorAccent).
				Bold(true)

	NormalItemStyle = lipgloss.NewStyle().
			Foreground(colorText)
)
