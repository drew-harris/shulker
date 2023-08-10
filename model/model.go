package model

import (
	"regexp"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	dtypes "github.com/docker/docker/api/types"
	"github.com/drewharris/shulker/engine"
	"github.com/drewharris/shulker/styles"
	"github.com/drewharris/shulker/types"
)

type LoadingModel struct {
	spinner       spinner.Model
	loadingOutput []string
}

type ServerExec struct {
	Connection dtypes.HijackedResponse
}

type ViewSelection int

const (
	loadingView ViewSelection = iota
	normalView
	shutdownView
	buildLogsView
	helpView
)

type MainModel struct {
	// TODO: CHANGE VIEW SELECTION TO ENUM
	isLoading          bool
	isShuttingDown     bool
	isBuilding         bool
	isViewingBuildLogs bool
	loadingModel       LoadingModel

	width  int
	height int

	engine engine.Engine

	outputChan     chan types.OutputMsg
	errorMessages  []string
	serverMessages []string
	buildMessages  []string
}

func (m MainModel) Init() tea.Cmd {
	var cmds []tea.Cmd
	cmds = append(cmds, m.loadingModel.spinner.Tick)
	cmds = append(cmds, ListenForOutput(m.outputChan))

	cmds = append(cmds, m.ensureSetupCmd())

	return tea.Batch(cmds...)
}

func (m MainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, DefaultKeyMap.Quit):
			m.isShuttingDown = true
			return m, m.Shutdown()

		case key.Matches(msg, DefaultKeyMap.ToggleBuildLogs):
			m.isViewingBuildLogs = !m.isViewingBuildLogs
			return m, nil

		case key.Matches(msg, DefaultKeyMap.Attach):
			// Print info in non alt screen
			// return m, tea.ExecProcess(exec.Command("docker", "attach", m.ConatainerId), func(err error) tea.Msg { return nil })
		case key.Matches(msg, DefaultKeyMap.RebuildAll):
			// Print info in non alt screen
			m.isBuilding = true
			return m, tea.Sequence(m.rebuildAllPlugins(), func() tea.Msg { return types.DoneBuilding })
		}

	case tea.WindowSizeMsg: // RESIZE
		m.width = msg.Width
		m.height = msg.Height

	// Channel output messages
	case types.OutputMsg:
		switch msg.Target {
		case types.StartupOutput:
			m.loadingModel.loadingOutput = append(m.loadingModel.loadingOutput, msg.Message)
			if len(m.loadingModel.loadingOutput) > m.height/2 {
				m.loadingModel.loadingOutput = m.loadingModel.loadingOutput[1:]
			}
		case types.ErrorOutput:
			m.errorMessages = append(m.errorMessages, msg.Message)
		case types.ServerOutput:
			m.serverMessages = append(m.serverMessages, msg.Message)
		case types.BuildOutput:
			m.buildMessages = append(m.buildMessages, msg.Message)

		}
		return m, ListenForOutput(m.outputChan)

	case types.QuickMsg:
		switch msg {
		case types.DoneBuilding:
			m.isBuilding = false
			return m, nil
		case types.BuildStarted:
			m.isBuilding = true
			return m, nil
		case types.ErrorBuilding:
			m.isBuilding = false
			return m, nil
		case types.FinishedSetup:
			m.isLoading = false
			return m, m.startServerCmd()
		case types.FinishedServerStart:
			return m, nil

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
		errors := styles.Error.Render(strings.Join(m.errorMessages, "\n"))
		loadingStyle := lipgloss.NewStyle().Padding(2).Bold(true).Italic(true)
		loading := loadingStyle.Height(m.height / 2).Width(m.width / 3).Align(lipgloss.Center).AlignVertical(lipgloss.Center).Render(m.loadingModel.spinner.View() + "   Starting Server...")

		var nonAlphanumericRegex = regexp.MustCompile(`[^\x20-\x7e]`)
		for i, str := range m.loadingModel.loadingOutput {
			m.loadingModel.loadingOutput[i] = nonAlphanumericRegex.ReplaceAllString(str, "")
		}
		output := styles.Dimmed.Render(strings.Join(m.loadingModel.loadingOutput, "\n"))
		screen := lipgloss.NewStyle().Margin(1).MaxWidth(m.width - 8).Height(m.height).AlignVertical(lipgloss.Center)
		return screen.Render(
			lipgloss.JoinVertical(
				lipgloss.Left,
				errors,
				lipgloss.JoinHorizontal(lipgloss.Top, loading, output)),
		)
	} else if m.isShuttingDown {
		return lipgloss.NewStyle().Width(m.width).Height(m.height).Align(lipgloss.Center).AlignVertical(lipgloss.Center).Render(m.loadingModel.spinner.View() + " Shutting down... Please be patient")
	} else if m.isViewingBuildLogs {
		return strings.Join(m.buildMessages, "\n")
	}

	// Main interface
	doc := strings.Builder{}
	var remainingHeight int = m.height

	// Status bar
	statusStyle := lipgloss.NewStyle().Width(m.width).Background(lipgloss.Color("#555"))
	var statusBar string
	if m.isBuilding {
		statusBar = statusStyle.Copy().Background(lipgloss.Color("#ab009d")).Render(m.loadingModel.spinner.View() + "  BUILDING... ")
	} else {
		statusBar = statusStyle.Render("IDLE")
	}
	remainingHeight = remainingHeight - lipgloss.Height(statusBar)

	leftMenuContainer := lipgloss.NewStyle().Padding(1).BorderStyle(lipgloss.NormalBorder()).BorderForeground(lipgloss.Color("#fff")).Height(9)
	menu := leftMenuContainer.Render("Shulker Menu\n* Rebuild\n* Restart")

	logContainer := lipgloss.NewStyle().Padding(1).BorderStyle(lipgloss.NormalBorder()).BorderForeground(lipgloss.Color("#fff")).Width(m.width - lipgloss.Width(menu) - 2).AlignVertical(lipgloss.Top)
	serverLogs := logContainer.Render(lastLines(m.serverMessages, m.height-1))
	// serverLogs := logContainer.Render("Helllo there")
	middle := lipgloss.JoinHorizontal(lipgloss.Bottom, menu, serverLogs)

	doc.WriteString(middle + "\n")
	doc.WriteString(statusBar)

	return doc.String()
}

func InitialModel(engine engine.Engine) MainModel {
	s := spinner.New()
	s.Spinner = spinner.Line

	model := MainModel{
		isLoading:  true,
		engine:     engine,
		outputChan: make(chan types.OutputMsg),
		loadingModel: LoadingModel{
			spinner:       s,
			loadingOutput: []string{},
		},
	}

	return model
}

func lastLines(strs []string, amt int) string {
	startIndex := len(strs) - amt
	if startIndex < 0 {
		startIndex = 0
	}

	lastElements := strs[startIndex:]
	return strings.Join(lastElements, "\n")
}
