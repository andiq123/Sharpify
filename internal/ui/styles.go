package ui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

var (
	
	Primary   = lipgloss.Color("#6366F1") 
	Secondary = lipgloss.Color("#10B981") 
	Warning   = lipgloss.Color("#F59E0B") 
	Error     = lipgloss.Color("#EF4444") 
	Muted     = lipgloss.Color("#9CA3AF") 
	Accent    = lipgloss.Color("#8B5CF6") 
	Info      = lipgloss.Color("#3B82F6") 

	
	TitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(Primary)

	SubtitleStyle = lipgloss.NewStyle().
			Foreground(Muted)

	SuccessStyle = lipgloss.NewStyle().
			Foreground(Secondary).
			Bold(true)

	WarningStyle = lipgloss.NewStyle().
			Foreground(Warning)

	ErrorStyle = lipgloss.NewStyle().
			Foreground(Error).
			Bold(true)

	InfoStyle = lipgloss.NewStyle().
			Foreground(Info)

	FileStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#60A5FA"))

	RuleStyle = lipgloss.NewStyle().
			Foreground(Secondary)

	AccentStyle = lipgloss.NewStyle().
			Foreground(Accent).
			Bold(true)

	
	BoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(Primary).
			Padding(1, 2)

	HeaderStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(Primary).
			Padding(0, 2)

	DiffAddStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#22C55E")).
			Background(lipgloss.Color("#052E16"))

	DiffRemoveStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#EF4444")).
			Background(lipgloss.Color("#450A0A"))

	
	DotSuccess = lipgloss.NewStyle().Foreground(Secondary).Render("●")
	DotWarning = lipgloss.NewStyle().Foreground(Warning).Render("●")
	DotError   = lipgloss.NewStyle().Foreground(Error).Render("●")
	DotInfo    = lipgloss.NewStyle().Foreground(Info).Render("●")
)

func Banner() string {
	banner := `
  ╔═══════════════════════════════════════════╗
  ║                                           ║
  ║   ⚡  S H A R P I F Y                     ║
  ║                                           ║
  ║   Modernize your C# code instantly        ║
  ║                                           ║
  ╚═══════════════════════════════════════════╝`
	return AccentStyle.Render(banner)
}

func SmallBanner() string {
	return AccentStyle.Render("⚡ Sharpify")
}

func Divider() string {
	return SubtitleStyle.Render("─────────────────────────────────────────────")
}

func StatusBadge(label string, value string, color lipgloss.Color) string {
	badge := lipgloss.NewStyle().
		Foreground(color).
		Bold(true).
		Render(value)
	return SubtitleStyle.Render(label+": ") + badge
}

func ProgressBar(current, total int, width int) string {
	if total == 0 {
		return ""
	}

	filled := int(float64(current) / float64(total) * float64(width))
	if filled > width {
		filled = width
	}

	bar := ""
	for i := 0; i < width; i++ {
		if i < filled {
			bar += "█"
		} else {
			bar += "░"
		}
	}

	return InfoStyle.Render(bar) + SubtitleStyle.Render(" "+fmt.Sprintf("%d/%d", current, total))
}


func Success(msg string) string {
	return SuccessStyle.Render("✓ ") + msg
}

func Warn(msg string) string {
	return WarningStyle.Render("⚠ ") + msg
}

func Fail(msg string) string {
	return ErrorStyle.Render("✗ ") + msg
}

func Tip(msg string) string {
	return InfoStyle.Render("💡 ") + SubtitleStyle.Render(msg)
}
