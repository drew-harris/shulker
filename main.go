package main

import (
	"context"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

// STYLES
var (
	errorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#ff8888"))
)

type MainModel struct {
	isFullScreen bool
	width        int
	height       int
	d            *client.Client
}

func (m MainModel) Init() tea.Cmd {
	return nil
}

func (m MainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}
	return m, nil
}

func (m MainModel) View() string {
	return "Hello there"
}

func InitialModel(client *client.Client) MainModel {
	return MainModel{
		isFullScreen: false,
		d:            client,
	}
}

func main() {
	os.Setenv("DOCKER_API_VERSION", "1.41")

	d, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		fmt.Println(errorStyle.Render("Error connecting to docker " + err.Error()))
		os.Exit(1)
	}

	// Test connection by getting containers
	_, err = d.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		fmt.Println(errorStyle.Render("\n   Can't connect to Docker, Is it running?"))
		os.Exit(1)
	}

	p := tea.NewProgram(InitialModel(d))

	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
