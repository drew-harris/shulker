package model

import (
	"github.com/charmbracelet/bubbles/key"
)

type KeyMap struct {
	Help key.Binding
	Quit key.Binding

	SendCmdToSpigot              key.Binding
	ToggleBuildLogs              key.Binding
	RebuildAll                   key.Binding
	RebuildAllNoCache            key.Binding
	ToggleReloadServerEveryBuild key.Binding
}

func (k KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit}
}

// FullHelp returns keybindings for the expanded help view. It's part of the
// key.Map interface.
func (k KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.SendCmdToSpigot, k.ToggleBuildLogs, k.ToggleReloadServerEveryBuild}, // second column
		{k.RebuildAll, k.RebuildAllNoCache},                                    // second column
		{k.Help, k.Quit},                                                       // second column
	}
}

var DefaultKeyMap = KeyMap{
	// Up: key.NewBinding(
	// 	key.WithKeys("k", "up"),        // actual keybindings
	// 	key.WithHelp("↑/k", "Move Up"), // corresponding help text
	// ),
	// Down: key.NewBinding(
	// 	key.WithKeys("j", "down"),
	// 	key.WithHelp("↓/j", "Move Down"),
	// ),

	Quit: key.NewBinding(
		key.WithKeys("q", "esc"),
		key.WithHelp("q", "Quit"),
	),

	Help: key.NewBinding(
		key.WithKeys("h", "?"),
		key.WithHelp("h", "Toggle Help"),
	),

	SendCmdToSpigot: key.NewBinding(
		key.WithKeys("c"),
		key.WithHelp("c", "Send command to spigot"),
	),

	ToggleBuildLogs: key.NewBinding(
		key.WithKeys("b"),
		key.WithHelp("b", "Toggle Build Logs"),
	),

	RebuildAll: key.NewBinding(
		key.WithKeys("r"),
		key.WithHelp("r", "Rebuild Plugins"),
	),

	RebuildAllNoCache: key.NewBinding(
		key.WithKeys("R"),
		key.WithHelp("R", "Rebuild Plugins (without cache)"),
	),

	ToggleReloadServerEveryBuild: key.NewBinding(
		key.WithKeys("t"),
		key.WithHelp("t", "Reload server every build (toggle)"),
	),
}
