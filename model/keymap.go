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
}

var DefaultKeyMap = KeyMap{
	Up: key.NewBinding(
		key.WithKeys("k", "up"),        // actual keybindings
		key.WithHelp("↑/k", "move up"), // corresponding help text
	),
	Down: key.NewBinding(
		key.WithKeys("j", "down"),
		key.WithHelp("↓/j", "move down"),
	),

	Quit: key.NewBinding(
		key.WithKeys("q", "esc", "ctrl+c"),
	),

	Select: key.NewBinding(
		key.WithKeys(tea.KeyEnter.String()),
	),

	Help: key.NewBinding(
		key.WithKeys("h", "?"),
	),

	Attach: key.NewBinding(
		key.WithKeys("a"),
	),

	ToggleBuildLogs: key.NewBinding(
		key.WithKeys("b"),
	),
}
