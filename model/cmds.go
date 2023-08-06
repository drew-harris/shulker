package model

import (
	"bufio"
	"context"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	dtypes "github.com/docker/docker/api/types"
	"github.com/drewharris/dockercraft/types"
)

func (m MainModel) startServerExecCmd() tea.Cmd {
	return func() tea.Msg {
		execId, err := m.d.ContainerExecCreate(context.TODO(), m.ConatainerId, dtypes.ExecConfig{
			Cmd:          []string{"./start.sh"},
			AttachStderr: true,
			AttachStdout: true,
			Tty:          true,
		})
		if err != nil {
			panic(err)
		}

		rd, err := m.d.ContainerExecAttach(context.Background(), execId.ID, dtypes.ExecStartCheck{})
		m.d.ContainerExecStart(context.Background(), execId.ID, dtypes.ExecStartCheck{
			ConsoleSize: &[2]uint{800, 4},
		})

		if err != nil {
			panic(err)
		}

		go func() {
			scanner := bufio.NewScanner(rd.Reader) // Scanner doesn't return newline byte
			for scanner.Scan() {
				m.outputChan <- types.OutputMsg{
					Target:  types.ServerOutput,
					Message: strings.ReplaceAll(scanner.Text(), "\n", ""),
				}
			}
		}()
		return ServerExec{
			ExecId:     execId.ID,
			Connection: rd,
		}
	}
}
