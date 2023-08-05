package model

import (
	"fmt"
	"math/rand"
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	dTypes "github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/drewharris/dockercraft/docker"
	"github.com/drewharris/dockercraft/types"
)

type LoadingModel struct {
	spinner       spinner.Model
	loadingOutput []string
	outputChan    chan types.ResponseMsg
}

type MainModel struct {
	isFullScreen bool
	isLoading    bool
	image        *dTypes.ImageSummary
	loadingModel LoadingModel
	width        int
	height       int
	d            *client.Client
}

func (m MainModel) Init() tea.Cmd {
	var cmds []tea.Cmd
	cmds = append(cmds, m.loadingModel.spinner.Tick)
	cmds = append(cmds, docker.ListenInitialBuild(m.loadingModel.outputChan))
	cmds = append(cmds, docker.TryInitialBuild(m.loadingModel.outputChan))

	return tea.Batch(cmds...)
}

func (m MainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			return m, tea.Quit
		case "t":
			randomNum := rand.Intn(100)
			m.loadingModel.outputChan <- types.ResponseMsg{
				Target:  types.StartupResponse,
				Message: "TEst" + fmt.Sprintf(" %d", randomNum),
			}
			return m, nil
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	// Channel response messages
	case types.ResponseMsg:
		switch msg.Target {
		case types.StartupResponse:
			m.loadingModel.loadingOutput = append(m.loadingModel.loadingOutput, string(msg.Message))
			if len(m.loadingModel.loadingOutput) > m.height/3 {
				m.loadingModel.loadingOutput = m.loadingModel.loadingOutput[1:]
			}
			return m, docker.ListenInitialBuild(m.loadingModel.outputChan)
		}

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

func InitialModel(client *client.Client, image *dTypes.ImageSummary) MainModel {
	s := spinner.New()
	s.Spinner = spinner.Line
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("#fff"))

	model := MainModel{
		isFullScreen: false,
		isLoading:    true,
		d:            client,
		image:        image,
		loadingModel: LoadingModel{
			spinner:       s,
			loadingOutput: []string{},
			outputChan:    make(chan types.ResponseMsg),
		},
	}

	return model
}
