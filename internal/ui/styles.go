package ui

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	// Colors
	Primary   = lipgloss.Color("#7C3AED")
	Secondary = lipgloss.Color("#10B981")
	Warning   = lipgloss.Color("#F59E0B")
	Error     = lipgloss.Color("#EF4444")
	Muted     = lipgloss.Color("#6B7280")

	// Styles
	TitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(Primary).
			MarginBottom(1)

	SubtitleStyle = lipgloss.NewStyle().
			Foreground(Muted).
			Italic(true)

	SuccessStyle = lipgloss.NewStyle().
			Foreground(Secondary).
			Bold(true)

	WarningStyle = lipgloss.NewStyle().
			Foreground(Warning)

	ErrorStyle = lipgloss.NewStyle().
			Foreground(Error).
			Bold(true)

	FileStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#60A5FA"))

	RuleStyle = lipgloss.NewStyle().
			Foreground(Secondary)

	BoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(Primary).
			Padding(1, 2)

	DiffAddStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#22C55E")).
			Background(lipgloss.Color("#052E16"))

	DiffRemoveStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#EF4444")).
			Background(lipgloss.Color("#450A0A"))
)

func Banner() string {
	banner := `
  ____  _                       _  __
 / ___|| |__   __ _ _ __ _ __ (_)/ _|_   _
 \___ \| '_ \ / _' | '__| '_ \| | |_| | | |
  ___) | | | | (_| | |  | |_) | |  _| |_| |
 |____/|_| |_|\__,_|_|  | .__/|_|_|  \__, |
                        |_|          |___/
`
	return TitleStyle.Render(banner)
}
