package model

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

type KeyMap struct {
	Up     key.Binding
	Down   key.Binding
	Select key.Binding
	Help   key.Binding
	Quit   key.Binding

	Attach          key.Binding
	ToggleBuildLogs key.Binding
	RebuildAll      key.Binding
}

var DefaultKeyMap = KeyMap{
	Up: key.NewBinding(
		key.WithKeys("k", "up"),        // actual keybindings
		key.WithHelp("↑/k", "Move Up"), // corresponding help text
	),
	Down: key.NewBinding(
		key.WithKeys("j", "down"),
		key.WithHelp("↓/j", "Move Down"),
	),

	Quit: key.NewBinding(
		key.WithKeys("q", "esc", "ctrl+c"),
		key.WithHelp("q / esc / ctrl+c", "Quit"),
	),

	Select: key.NewBinding(
		key.WithKeys(tea.KeyEnter.String()),
	),

	Help: key.NewBinding(
		key.WithKeys("h", "?"),
		key.WithHelp("h", "Show Help"),
	),

	Attach: key.NewBinding(
		key.WithKeys("a"),
		key.WithHelp("a", "Attach To Server"),
	),

	ToggleBuildLogs: key.NewBinding(
		key.WithKeys("b"),
	),

	RebuildAll: key.NewBinding(
		key.WithKeys("r"),
	),
}
