package styles

import "github.com/charmbracelet/lipgloss"

var (
	Error     = lipgloss.NewStyle().Foreground(lipgloss.Color("#ff8888"))
	Highlight = lipgloss.NewStyle().Foreground(lipgloss.Color("#23c223"))
	Dimmed    = lipgloss.NewStyle().Foreground(lipgloss.Color("#707070"))
	Purple    = lipgloss.NewStyle().Foreground(lipgloss.Color("#bf34ed"))
)
