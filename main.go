package main

import (
	"fmt"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
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

func InitialModel() MainModel {
	return MainModel{
		isFullScreen: false,
	}
}

func main() {
	p := tea.NewProgram(InitialModel())

	_, err := client.NewClientWithOpts(client.FromEnv)
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(errorStyle.Copy().Foreground(lipgloss.Color("#f07f60")).Render("Error connecting to docker"))
			log.Println("panic occurred:", err)
			log.Println(errorStyle.Render("Are you running docker???"))
		}
	}()
	if err == nil {
		fmt.Println(errorStyle.Render(err.Error()))
		fmt.Println(errorStyle.Copy().Foreground(lipgloss.Color("#f07f60")).Render("Error connecting to docker"))
	}

	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
