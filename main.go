package main

import (
	"context"
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/charmbracelet/lipgloss"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

// STYLES
var (
	errorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#ff8888"))
)

type LoadingModel struct {
	spinner spinner.Model
}

type MainModel struct {
	isFullScreen bool
	isLoading    bool
	loadingModel LoadingModel
	width        int
	height       int
	d            *client.Client
}

func (m MainModel) Init() tea.Cmd {
	var cmds []tea.Cmd
	cmds = append(cmds, m.loadingModel.spinner.Tick)
	return tea.Batch(cmds...)

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

	default:
		var cmd tea.Cmd
		m.loadingModel.spinner, cmd = m.loadingModel.spinner.Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m MainModel) View() string {
	if m.isLoading {
		loadingStyle := lipgloss.NewStyle().Padding(2).Bold(true)
		return loadingStyle.Render(m.loadingModel.spinner.View() + " Starting server")
	}
	return "Hello there"
}

func InitialModel(client *client.Client) MainModel {
	s := spinner.New()
	s.Spinner = spinner.Line
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("#fff"))

	return MainModel{
		isFullScreen: false,
		isLoading:    true,
		d:            client,
		loadingModel: LoadingModel{
			spinner: s,
		},
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
