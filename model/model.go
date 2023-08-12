package model

import (
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	dtypes "github.com/docker/docker/api/types"
	"github.com/drewharris/shulker/config"
	"github.com/drewharris/shulker/engine"
	"github.com/drewharris/shulker/types"
)

type ServerExec struct {
	Connection dtypes.HijackedResponse
}

type Loggers struct {
	error   types.Logger
	build   types.Logger
	server  types.Logger
	startup types.Logger
}

type viewMode int

const (
	serverView viewMode = iota
	startupView
	buildView
	shutdownView
	testView
)

type MainModel struct {
	width  int
	height int

	engine   engine.Engine
	config   config.Config
	keys     KeyMap
	viewMode viewMode

	isBuilding bool

	outputChan     chan types.OutputMsg
	errorMessages  []string
	serverMessages []string
	buildMessages  []string
	loadingOutput  []string

	loggers Loggers

	reloadSpigotOnBuild bool

	spinner  spinner.Model
	help     help.Model
	cmdInput textinput.Model
	viewport viewport.Model
}

func (m MainModel) Init() tea.Cmd {
	var cmds []tea.Cmd
	cmds = append(cmds, m.spinner.Tick)
	cmds = append(cmds, ListenForOutput(m.outputChan))
	cmds = append(cmds, m.ensureSetupCmd())
	cmds = append(cmds, textinput.Blink)

	return tea.Batch(cmds...)
}

func (m MainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:

		if msg.String() == "u" {
			m.viewport.GotoBottom()
			return m, nil
		}

		if m.cmdInput.Focused() {
			if msg.String() == tea.KeyEnter.String() {
				cmd := m.cmdInput.Value()
				m.cmdInput.SetValue("")
				m.cmdInput.Blur()
				return m, tea.Sequence(m.sendCmdToSpigot(cmd))
			} else if msg.String() == tea.KeyEsc.String() {
				m.cmdInput.SetValue("")
				m.cmdInput.Blur()
				return m, nil
			} else {
				newInput, cmd := m.cmdInput.Update(msg)
				m.cmdInput = newInput
				return m, cmd
			}
		}

		switch {
		case key.Matches(msg, m.keys.Quit):
			m.viewMode = shutdownView
			return m, m.Shutdown()

		case key.Matches(msg, m.keys.ToggleBuildLogs):
			if m.viewMode == serverView {
				m.viewMode = buildView
				m.viewport.SetContent(strings.Join(m.buildMessages, "\n"))
			} else if m.viewMode == buildView {
				m.viewMode = serverView
				m.viewport.SetContent(strings.Join(m.serverMessages, "\n"))
			}

			m.viewport.GotoBottom()
			return m, nil

		case key.Matches(msg, m.keys.SendCmdToSpigot):
			// Print info in non alt screen
			m.cmdInput.Focus()
			return m, nil
		case key.Matches(msg, m.keys.RebuildAll):
			// Print info in non alt screen
			m.isBuilding = true
			return m, tea.Sequence(m.rebuildAllPlugins(false), func() tea.Msg { return types.DoneBuilding })
		case key.Matches(msg, m.keys.RebuildAllNoCache):
			// Print info in non alt screen
			m.isBuilding = true
			return m, tea.Sequence(m.rebuildAllPlugins(true), func() tea.Msg { return types.DoneBuilding })
		case key.Matches(msg, m.keys.Help):
			m.help.ShowAll = !m.help.ShowAll
			return m, nil
		case key.Matches(msg, m.keys.ToggleReloadServerEveryBuild):
			m.reloadSpigotOnBuild = !m.reloadSpigotOnBuild
			return m, nil
		}

	case tea.WindowSizeMsg: // RESIZE
		m.width = msg.Width
		m.height = msg.Height
		m.viewport.Height = 30
		m.viewport.Width = 90

	// Channel output messages
	case types.OutputMsg:
		switch msg.Target {
		case types.StartupOutput:
			m.loadingOutput = append(m.loadingOutput, msg.Message)
			if len(m.loadingOutput) > m.height/2 {
				m.loadingOutput = m.loadingOutput[1:]
			}
		case types.ErrorOutput:
			m.errorMessages = append(m.errorMessages, msg.Message)
		case types.ServerOutput:
			m.serverMessages = append(m.serverMessages, msg.Message)
		case types.BuildOutput:
			m.buildMessages = append(m.buildMessages, msg.Message)
		}

		if m.viewMode == serverView {
			m.viewport.SetContent(strings.Join(m.serverMessages, "\n"))
		} else if m.viewMode == buildView {
			m.viewport.SetContent(strings.Join(m.buildMessages, "\n"))
		}

		m.viewport.GotoBottom()
		return m, ListenForOutput(m.outputChan)

	case types.QuickMsg:
		switch msg {
		case types.DoneBuilding:
			m.isBuilding = false
			if m.reloadSpigotOnBuild {
				m.loggers.server("trying to reload")
				return m, m.sendCmdToSpigot("reload")
			} else {
				return m, nil
			}
		case types.BuildStarted:
			m.isBuilding = true
			return m, nil
		case types.ErrorBuilding:
			m.isBuilding = false
			return m, nil
		case types.FinishedSetup:
			m.viewMode = serverView
			return m, m.startServerCmd()
		case types.FinishedServerStart:
			return m, nil

		}

	default:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}

	return m, nil
}

func generateLogFn(sub chan types.OutputMsg, target types.OutputTarget) func(msg string) {
	return func(msg string) {
		sub <- types.OutputMsg{
			Target:  target,
			Message: msg,
		}
	}
}

func InitialModel(engine engine.Engine, config config.Config) MainModel {
	s := spinner.New()
	s.Spinner = spinner.Line

	outputChan := make(chan types.OutputMsg)
	ti := textinput.New()
	ti.Placeholder = "Send command to server..."

	vp := viewport.New(0, 0)
	vp.Style = lipgloss.NewStyle().Border(lipgloss.NormalBorder()).Padding(1)

	model := MainModel{
		engine:              engine,
		viewMode:            startupView,
		outputChan:          outputChan,
		keys:                DefaultKeyMap,
		help:                help.New(),
		viewport:            vp,
		config:              config,
		spinner:             s,
		cmdInput:            ti,
		reloadSpigotOnBuild: false,
		loggers: Loggers{
			error:   generateLogFn(outputChan, types.ErrorOutput),
			build:   generateLogFn(outputChan, types.BuildOutput),
			server:  generateLogFn(outputChan, types.ServerOutput),
			startup: generateLogFn(outputChan, types.StartupOutput),
		},
	}

	return model
}

func lastLines(strs string, amt int) string {
	startIndex := len(strs) - amt
	if startIndex < 0 {
		startIndex = 0
	}

	lines := strings.Split(strs, "\n")
	if len(lines) > amt {
		lines = lines[len(lines)-amt:]
	}

	strs = strings.Join(lines, "\n")
	return strs
}
