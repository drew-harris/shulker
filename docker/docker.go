package docker

import (
	"context"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	dtypes "github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/drewharris/dockercraft/commands"
	"github.com/drewharris/dockercraft/styles"
	"github.com/drewharris/dockercraft/types"
)

func PrepareContainerCmd(sub chan types.OutputMsg, d *client.Client) tea.Cmd {
	return func() tea.Msg {
		images, err := d.ImageList(context.Background(), dtypes.ImageListOptions{All: true})
		if err != nil {
			fmt.Println(styles.Error.Render("\n   Can't connect to Docker, Is it running?"))
			os.Exit(1)
		}

		var imageId = ""
		for _, image := range images {
			tag := image.RepoTags[0]
			if tag == "dockercraft:latest" {
				imageId = image.ID
			}
		}

		// - Check if image exists -> build
		if imageId == "" { // If there is no image
			commands.RunExternalCommand(sub, commands.Command{
				Name:   "docker",
				Args:   []string{"build", "-t", "dockercraft", "-f", "Dockerfile.dev", "."},
				Target: types.StartupOutput,
			})
		}

		// - Check if container exists -> create
		containers, err := d.ContainerList(context.Background(), dtypes.ContainerListOptions{
			All: true,
		})
		if err != nil {
			sub <- types.OutputMsg{
				Target:  types.ErrorOutput,
				Message: err.Error(),
			}
		}

		var containerId = ""
		for _, cont := range containers {
			if cont.ImageID == imageId {
				containerId = cont.ID
			}
		}

		// if containerId == "" {
		// 	// Create container
		// }

		// - Start container
		// - Check if all plugins built?
		// - Build plugins if not

		return types.FinishedSetupCmd{
			ImageId:     imageId,
			ContainerId: containerId,
		}
	}
}
