package docker

import (
	"bufio"
	"os/exec"
	"strings"

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
		cmd := exec.Command("docker", "build", "-t", "dockercraft", "-f", "Dockerfile.dev", ".")
		cmd.Dir = "/Users/drew/programs/mc-docker"

		// Display output
		stdout, err := cmd.StdoutPipe()
		if err != nil {
			panic(err)
		}

		// Combine stdout and stderr so that both are captured
		cmd.Stderr = cmd.Stdout

		// Start the command
		if err := cmd.Start(); err != nil {
			panic(err)
		}

		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			sub <- types.ResponseMsg{
				Target: types.StartupResponse,
				// Sanitize
				Message: strings.ReplaceAll(scanner.Text(), "\n", ""),
			}
		}

		if err := cmd.Wait(); err != nil {
			panic(err)
		}
		return nil
	}
}
