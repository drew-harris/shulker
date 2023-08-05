package model

import (
	"fmt"
	"math/rand"
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/docker/docker/client"
)

type LoadingModel struct {
	spinner       spinner.Model
	loadingOutput []string
	outputChan    chan responseMsg
}

type MainModel struct {
	isFullScreen bool
	isLoading    bool
	loadingModel LoadingModel
	width        int
	height       int
	d            *client.Client
}

type responseMsg struct {
	target  responseTarget
	message string
}

type responseTarget string

var (
	startupResponse responseTarget = "startup"
)

func (m MainModel) Init() tea.Cmd {
	var cmds []tea.Cmd
	cmds = append(cmds, m.loadingModel.spinner.Tick)
	cmds = append(cmds, listenInitialBuild(m.loadingModel.outputChan))
	return tea.Batch(cmds...)
}

func listenInitialBuild(sub chan responseMsg) tea.Cmd {
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
			m.loadingModel.outputChan <- responseMsg{
				target:  startupResponse,
				message: "TEst" + fmt.Sprintf(" %d", randomNum),
			}
			return m, nil
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	// Channel response messages
	case responseMsg:
		switch msg.target {
		case startupResponse:
			m.loadingModel.loadingOutput = append(m.loadingModel.loadingOutput, string(msg.message))
			if len(m.loadingModel.loadingOutput) > m.height/3 {
				m.loadingModel.loadingOutput = m.loadingModel.loadingOutput[1:]
			}
			return m, listenInitialBuild(m.loadingModel.outputChan)
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
			outputChan:    make(chan responseMsg),
		},
	}

	return model
}
