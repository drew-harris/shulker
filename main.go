package main

import (
	"context"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/drewharris/dockercraft/model"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

// STYLES
var (
	errorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#ff8888"))
)

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

	var dcImage types.ImageSummary
	for _, image := range images {
		tag := image.RepoTags[0]
		if tag == "dockercraft:latest" {
			dcImage = image
		}
	}

	p := tea.NewProgram(model.InitialModel(d, &dcImage))

	f, err := tea.LogToFile("debug.log", "debug")
	if err != nil {
		fmt.Println("fatal:", err)
		os.Exit(1)
	}

	defer f.Close()
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
