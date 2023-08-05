package docker

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/drewharris/dockercraft/types"
)

func ListenInitialBuild(sub chan types.ResponseMsg) tea.Cmd {
	return func() tea.Msg {
		return types.ResponseMsg(<-sub)
	}
}
