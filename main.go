package main

import (
	"context"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/drewharris/dockercraft/model"
	"github.com/drewharris/dockercraft/styles"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

// STYLES
func main() {
	os.Setenv("DOCKER_API_VERSION", "1.41")

	d, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		fmt.Println(styles.Error.Render("Error connecting to docker " + err.Error()))
		os.Exit(1)
	}

	_, err = d.ImageList(context.Background(), types.ImageListOptions{All: true})
	if err != nil {
		fmt.Println(styles.Error.Render("\n   Can't connect to Docker, Is it running? \n" + err.Error()))
		os.Exit(1)
	}

	// Test connection by getting containers
	p := tea.NewProgram(model.InitialModel(d))

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
