package model

import (
	"bufio"
	"context"

	tea "github.com/charmbracelet/bubbletea"
	dtypes "github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/drewharris/dockercraft/commands"
	"github.com/drewharris/dockercraft/docker"
	"github.com/drewharris/dockercraft/types"
)

func ListenForOutput(sub chan types.OutputMsg) tea.Cmd {
	return func() tea.Msg {
		return types.OutputMsg(<-sub)
	}
}

func (m MainModel) startServerExecCmd() tea.Cmd {
	return func() tea.Msg {
		waiter, err := m.d.ContainerAttach(context.Background(), m.ConatainerId, dtypes.ContainerAttachOptions{
			Stderr: true,
			Stdout: true,
			Stdin:  true,
			Stream: true,
		})
		if err != nil {
			panic("Couldn't attach to container " + err.Error())
		}

		go func() {
			scanner := bufio.NewScanner(waiter.Reader)
			for scanner.Scan() {
				m.outputChan <- types.OutputMsg{
					Target:  types.ServerOutput,
					Message: scanner.Text(),
				}
			}
		}()

		waiter.Conn.Write([]byte("./start.sh\n"))
		return ServerExec{
			Connection: waiter,
		}
	}
}

func (m *MainModel) reloadServer() {
	m.ServerExec.Connection.Conn.Write([]byte("help\n"))
}

func (m *MainModel) shutdown() tea.Cmd {
	return func() tea.Msg {
		m.ServerExec.Connection.Conn.Close() // Close running server
		m.d.ContainerStop(context.Background(), m.ConatainerId, container.StopOptions{})
		return tea.Quit()
	}
}

func (m *MainModel) rebuildAllPlugins() tea.Cmd {
	return func() tea.Msg {
		docker.RunContainerCommandAsync(m.d, m.ConatainerId, m.outputChan, commands.Command{
			Target: types.BuildOutput,
			Name:   "./build_all.sh",
		}, nil)

		return nil // TODO: CHANGE
	}
}
