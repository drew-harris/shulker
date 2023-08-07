package model

import (
	"bufio"
	"context"

	tea "github.com/charmbracelet/bubbletea"
	dtypes "github.com/docker/docker/api/types"
	"github.com/drewharris/dockercraft/types"
)

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
