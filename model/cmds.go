package model

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/drewharris/shulker/types"
)

func generateLogFn(sub chan types.OutputMsg, target types.OutputTarget) func(msg string) {
	return func(msg string) {
		sub <- types.OutputMsg{
			Target:  target,
			Message: msg,
		}
	}
}

func ListenForOutput(sub chan types.OutputMsg) tea.Cmd {
	return func() tea.Msg {
		return types.OutputMsg(<-sub)
	}
}

func (m *MainModel) ensureSetupCmd() tea.Cmd {
	return func() tea.Msg {
		err := m.engine.EnsureSetup(m.outputChan)
		if err != nil {
			panic(err)
		}

		return types.FinishedSetup
	}
}

func (m *MainModel) startServerCmd() tea.Cmd {
	return func() tea.Msg {
		err := m.engine.StartServer(m.outputChan)
		if err != nil {
			panic(err)
		}

		return types.FinishedServerStart
	}
}

func (m *MainModel) Shutdown() tea.Cmd {
	return func() tea.Msg {
		err := m.engine.Shutdown()
		if err != nil {
			panic(err)
		}
		return tea.Quit()
	}
}

func (m *MainModel) rebuildAllPlugins() tea.Cmd {
	return func() tea.Msg {
		err := m.engine.RebuildAllPlugins(m.outputChan)
		if err != nil {
			return types.ErrorBuilding
		}

		return types.BuildStarted
	}
}
