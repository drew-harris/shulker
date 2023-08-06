package model

import (
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/docker/docker/client"
	"github.com/drewharris/dockercraft/commands"
	"github.com/drewharris/dockercraft/types"
)

type LoadingModel struct {
	spinner       spinner.Model
	loadingOutput []string
}

type MainModel struct {
	isFullScreen  bool
	isLoading     bool
	imageId       string
	loadingModel  LoadingModel
	width         int
	height        int
	d             *client.Client
	outputChan    chan types.ResponseMsg
	errorMessages []string
}

func ListenForOutput(sub chan types.ResponseMsg) tea.Cmd {
	return func() tea.Msg {
		return types.ResponseMsg(<-sub)
	}
}

func (m MainModel) Init() tea.Cmd {
	var cmds []tea.Cmd
	cmds = append(cmds, m.loadingModel.spinner.Tick)
	cmds = append(cmds, ListenForOutput(m.outputChan))

	cmds = append(cmds, commands.TeaRunCommandWithOutput(m.outputChan, commands.Command{
		Name:   "docker",
		Args:   []string{"build", "-t", "dockercraft", "-f", "Dockerfile.dev", "."},
		Target: types.StartupResponse,
	}, nil))

	return tea.Batch(cmds...)
}

func (m MainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c": // QUIT
			return m, tea.Quit
		}

	case tea.WindowSizeMsg: // RESIZE
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
		case types.ErrorResponse:
			m.errorMessages = append(m.errorMessages, msg.Message)
		}
		return m, ListenForOutput(m.outputChan)

	default:
		var cmd tea.Cmd
		m.loadingModel.spinner, cmd = m.loadingModel.spinner.Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m MainModel) View() string {
	if m.isLoading {
		errors := lipgloss.NewStyle().Foreground(lipgloss.Color("#ff0000")).Render(strings.Join(m.errorMessages, "\n"))
		loadingStyle := lipgloss.NewStyle().Padding(2).Bold(true)
		loading := loadingStyle.Height(m.height / 3).Width(m.width / 3).Align(lipgloss.Center).AlignVertical(lipgloss.Center).Render(m.loadingModel.spinner.View() + "   Starting server")
		output := lipgloss.NewStyle().Foreground(lipgloss.Color("#707070")).Render(strings.Join(m.loadingModel.loadingOutput, "\n"))
		screen := lipgloss.NewStyle().Margin(1).MaxWidth(m.width - 8)
		return screen.Render(
			lipgloss.JoinVertical(
				lipgloss.Left,
				errors,
				lipgloss.JoinHorizontal(lipgloss.Top, loading, output)),
		)
	}
	return "Hello there"
}

func InitialModel(client *client.Client, imageId string) MainModel {
	s := spinner.New()
	s.Spinner = spinner.Line
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("#fff"))

	model := MainModel{
		isFullScreen: false,
		isLoading:    true,
		d:            client,
		imageId:      imageId,
		outputChan:   make(chan types.ResponseMsg),
		loadingModel: LoadingModel{
			spinner:       s,
			loadingOutput: []string{},
		},
	}

	return model
}
