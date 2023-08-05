package docker

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/drewharris/dockercraft/types"
)

func ListenInitialBuild(sub chan types.ResponseMsg) tea.Cmd {
	return func() tea.Msg {
		return types.ResponseMsg(<-sub)
	}
}

func TryInitialBuild(sub chan types.ResponseMsg) tea.Cmd {
	return func() tea.Msg {
		for {
			time.Sleep(time.Second * 4)
			sub <- types.ResponseMsg{
				Target:  types.StartupResponse,
				Message: "No Image Found",
			}
		}
	}
}
