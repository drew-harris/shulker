package main

import (
	"context"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/drewharris/shulker/engine"
	"github.com/drewharris/shulker/model"
	"github.com/drewharris/shulker/styles"
	"github.com/integrii/flaggy"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

var version string = "unversioned"

func main() {
	flaggy.SetName("Shulker")
	flaggy.SetVersion(version)
	flaggy.SetDescription("A tui (terminal user interface) for efficient plugin development at the HuMIn Game Labs")

	var noDocker bool

	flaggy.Bool(&noDocker, "n", "no-docker", "Use the host system instead of docker")

	flaggy.Parse()

	var program *tea.Program

	if noDocker {
		hostEngine, err := engine.NewHostEngine()
		if err != nil {
			panic(err)
		}
		program = tea.NewProgram(model.InitialModel(hostEngine), tea.WithAltScreen())
	} else {
		os.Setenv("DOCKER_API_VERSION", "1.41")

		d, err := client.NewClientWithOpts(client.FromEnv)
		if err != nil {
			fmt.Println(styles.Error.Render("Error connecting to docker " + err.Error()))
			os.Exit(1)
		}

		// Test connection by getting containers
		_, err = d.ImageList(context.Background(), types.ImageListOptions{All: true})
		if err != nil {
			fmt.Println(styles.Error.Render("\n   Can't connect to Docker, Is it running?"))
			os.Exit(1)
		}

		dEngine := engine.NewDockerEngine(d)

		program = tea.NewProgram(model.InitialModel(dEngine), tea.WithAltScreen())
	}

	// Docker

	f, err := tea.LogToFile("debug.log", "debug")
	if err != nil {
		fmt.Println("fatal:", err)
		os.Exit(1)
	}

	defer f.Close()
	if _, err := program.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
