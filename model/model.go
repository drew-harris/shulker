package model

import (
	"context"
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	dtypes "github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/drewharris/dockercraft/docker"
	"github.com/drewharris/dockercraft/styles"
	"github.com/drewharris/dockercraft/types"
)

type LoadingModel struct {
	spinner       spinner.Model
	loadingOutput []string
}

type ServerExec struct {
	ExecId     string
	Connection dtypes.HijackedResponse
}

type ExecId string

type MainModel struct {
	isLoading    bool
	loadingModel LoadingModel

	width  int
	height int
	d      *client.Client

	ImageId      string
	ConatainerId string

	outputChan    chan types.OutputMsg
	errorMessages []string

	ServerExec ServerExec
	OtherExecs []ExecId
}

func ListenForOutput(sub chan types.OutputMsg) tea.Cmd {
	return func() tea.Msg {
		return types.OutputMsg(<-sub)
	}
}

func (m MainModel) Init() tea.Cmd {
	var cmds []tea.Cmd
	cmds = append(cmds, m.loadingModel.spinner.Tick)
	cmds = append(cmds, ListenForOutput(m.outputChan))

	cmds = append(cmds, docker.PrepareContainerCmd(m.outputChan, m.d))

	return tea.Batch(cmds...)
}

func (m *MainModel) Shutdown() {
	m.d.ContainerStop(context.Background(), m.ConatainerId, container.StopOptions{})
}

func (m MainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c": // QUIT
			m.Shutdown()
			return m, tea.Quit
		}

	case tea.WindowSizeMsg: // RESIZE
		m.width = msg.Width
		m.height = msg.Height

	// Channel output messages
	case types.OutputMsg:
		switch msg.Target {
		case types.StartupOutput:
			m.loadingModel.loadingOutput = append(m.loadingModel.loadingOutput, msg.Message)
			if len(m.loadingModel.loadingOutput) > m.height/3 {
				m.loadingModel.loadingOutput = m.loadingModel.loadingOutput[1:]
			}
		case types.ErrorOutput:
			m.errorMessages = append(m.errorMessages, msg.Message)
		}
		return m, ListenForOutput(m.outputChan)

	case types.FinishedSetupCmd:
		m.ConatainerId = msg.ContainerId
		m.ImageId = msg.ImageId
		m.isLoading = false
		return m, tea.Batch(tea.EnterAltScreen, m.startServerExecCmd())

	case ServerExec:
		m.ServerExec = msg
		return m, nil

	default:
		var cmd tea.Cmd
		m.loadingModel.spinner, cmd = m.loadingModel.spinner.Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m MainModel) View() string {
	if m.isLoading {
		errors := styles.Error.Render(strings.Join(m.errorMessages, "\n"))
		loadingStyle := lipgloss.NewStyle().Padding(2).Bold(true).Italic(true)
		loading := loadingStyle.Height(m.height / 3).Width(m.width / 3).Align(lipgloss.Center).AlignVertical(lipgloss.Center).Render(m.loadingModel.spinner.View() + "   Starting Server...")

		// var nonAlphanumericRegex = regexp.MustCompile(`[^\x20-\x7e]`)
		// for i, str := range m.loadingModel.loadingOutput {
		// 	m.loadingModel.loadingOutput[i] = nonAlphanumericRegex.ReplaceAllString(str, "")
		// }
		output := styles.Dimmed.Render(strings.Join(m.loadingModel.loadingOutput, "\n"))
		screen := lipgloss.NewStyle().Margin(1).MaxWidth(m.width - 8)
		return screen.Render(
			lipgloss.JoinVertical(
				lipgloss.Left,
				errors,
				lipgloss.JoinHorizontal(lipgloss.Top, loading, output)),
		)
	}

	errors := styles.Error.Render(strings.Join(m.errorMessages, "\n"))

	return lipgloss.JoinVertical(
		lipgloss.Center,
		errors,
		"Hellow",
	)
}

func InitialModel(client *client.Client) MainModel {
	s := spinner.New()
	s.Spinner = spinner.Line
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("#fff"))

	model := MainModel{
		isLoading:  true,
		d:          client,
		outputChan: make(chan types.OutputMsg),
		loadingModel: LoadingModel{
			spinner:       s,
			loadingOutput: []string{},
		},
	}

	return model
}
