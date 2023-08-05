package main

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"strings"

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
	spinner       spinner.Model
	loadingOutput []string
	outputChan    chan string
}

type MainModel struct {
	isFullScreen bool
	isLoading    bool
	loadingModel LoadingModel
	width        int
	height       int
	d            *client.Client
}

type responseMsg string

func (m MainModel) Init() tea.Cmd {
	var cmds []tea.Cmd
	cmds = append(cmds, m.loadingModel.spinner.Tick)
	cmds = append(cmds, waitForActivity(m.loadingModel.outputChan))
	return tea.Batch(cmds...)
}

func waitForActivity(sub chan string) tea.Cmd {
	return func() tea.Msg {
		return responseMsg(<-sub)
	}
}

func (m MainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			return m, tea.Quit
		case "t":
			randomNum := rand.Intn(100)
			m.loadingModel.outputChan <- "test" + fmt.Sprint(randomNum)
			return m, nil
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case responseMsg:
		m.loadingModel.loadingOutput = append(m.loadingModel.loadingOutput, string(msg))
		// limit output to half screen
		if len(m.loadingModel.loadingOutput) > m.height/3 {
			m.loadingModel.loadingOutput = m.loadingModel.loadingOutput[1:]
		}
		return m, waitForActivity(m.loadingModel.outputChan)

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
		loading := loadingStyle.Height(m.height / 3).Width(m.width / 4).Align(lipgloss.Center).AlignVertical(lipgloss.Center).Render(m.loadingModel.spinner.View() + "   Starting server")
		output := lipgloss.NewStyle().Width(m.width / 2).Render(strings.Join(m.loadingModel.loadingOutput, "\n"))

		return lipgloss.JoinHorizontal(lipgloss.Top, loading, output)
	}
	return "Hello there"
}

func InitialModel(client *client.Client) MainModel {
	s := spinner.New()
	s.Spinner = spinner.Line
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("#fff"))

	model := MainModel{
		isFullScreen: false,
		isLoading:    true,
		d:            client,
		loadingModel: LoadingModel{
			spinner:       s,
			loadingOutput: []string{},
			outputChan:    make(chan string),
		},
	}

	return model
}

func main() {
	os.Setenv("DOCKER_API_VERSION", "1.41")

	d, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		fmt.Println(errorStyle.Render("Error connecting to docker " + err.Error()))
		os.Exit(1)
	}

	// Test connection by getting containers
	images, err := d.ImageList(context.Background(), types.ImageListOptions{All: true})
	if err != nil {
		fmt.Println(errorStyle.Render("\n   Can't connect to Docker, Is it running?"))
		os.Exit(1)
	}

	for _, image := range images {
		tag := image.RepoTags[0]
		if tag == "dockercraft:latest" {
			fmt.Println("Found image")
		}
	}

	p := tea.NewProgram(InitialModel(d))

	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
