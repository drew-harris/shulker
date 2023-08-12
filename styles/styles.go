package styles

import "github.com/charmbracelet/lipgloss"

var (
	Error  = lipgloss.NewStyle().Foreground(lipgloss.Color("#ff8888"))
	Dimmed = lipgloss.NewStyle().Foreground(lipgloss.Color("#707070"))

	LeftMenuContainer = lipgloss.NewStyle().Padding(1, 3).BorderStyle(lipgloss.NormalBorder()).BorderForeground(lipgloss.Color("#fff"))
	InlineLog         = lipgloss.NewStyle().Inline(true).Italic(true).MaxWidth(30)
	LogContainer      = lipgloss.NewStyle().Padding(0, 1).AlignVertical(lipgloss.Top).AlignHorizontal(lipgloss.Left).Border(lipgloss.NormalBorder())
)
