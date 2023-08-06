package docker

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	dtypes "github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
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

		if containerId == "" {

			c, err := createContainer(d)
			if err != nil {
				sub <- types.OutputMsg{
					Target:  types.ErrorOutput,
					Message: "Could not create container: " + err.Error(),
				}
			}

			// Trim whitespace of output
			containerId = c.ID
			sub <- types.OutputMsg{
				Target:  types.StartupOutput,
				Message: "Created container: " + containerId,
			}
		} else {
			sub <- types.OutputMsg{
				Target:  types.StartupOutput,
				Message: "Found Container: " + containerId,
			}

		}

		// - Start container
		err = d.ContainerStart(context.Background(), containerId, dtypes.ContainerStartOptions{})
		if err != nil {
			sub <- types.OutputMsg{
				Target:  types.ErrorOutput,
				Message: "Could not start container: " + err.Error(),
			}
		}
		sub <- types.OutputMsg{
			Target:  types.StartupOutput,
			Message: "Container started...",
		}

		// - Check if all plugins built?
		// var pluginsBuilt bool
		// stat, err := d.ContainerStatPath(context.Background(), containerId, "/server/plugins/HuMInGameLabsPlugin.jar")
		var errors []error
		_, err = d.ContainerStatPath(context.Background(), containerId, "/server/plugins/HuMInGameLabsPlugin.jar")
		errors = append(errors, err)
		_, err = d.ContainerStatPath(context.Background(), containerId, "/server/plugins/Contraption.jar")
		errors = append(errors, err)
		_, err = d.ContainerStatPath(context.Background(), containerId, "/server/plugins/Recycler.jar")
		errors = append(errors, err)
		_, err = d.ContainerStatPath(context.Background(), containerId, "/server/plugins/SplitterNode.jar")
		errors = append(errors, err)

		var notBuilt bool = false
		for _, err := range errors {
			if err != nil {
				notBuilt = true
				break
			}
		}

		if notBuilt {
			execId, err := d.ContainerExecCreate(context.TODO(), containerId, dtypes.ExecConfig{
				Cmd:          []string{"./build_all.sh"},
				AttachStderr: true,
				AttachStdout: true,
				Tty:          true,
			})
			if err != nil {
				panic(err)
			}

			rd, err := d.ContainerExecAttach(context.Background(), execId.ID, dtypes.ExecStartCheck{})
			d.ContainerExecStart(context.Background(), execId.ID, dtypes.ExecStartCheck{
				ConsoleSize: &[2]uint{800, 4},
			})

			if err != nil {
				panic(err)
			}
			defer rd.Close()
			scanner := bufio.NewScanner(rd.Reader) // Scanner doesn't return newline byte
			for scanner.Scan() {
				sub <- types.OutputMsg{
					Target:  types.StartupOutput,
					Message: strings.ReplaceAll(scanner.Text(), "\n", ""),
				}
			}

		}

		sub <- types.OutputMsg{
			Target:  types.StartupOutput,
			Message: "All Plugins Built",
		}

		return types.FinishedSetupCmd{
			ImageId:     imageId,
			ContainerId: containerId,
		}
	}
}

func createContainer(d *client.Client) (container.CreateResponse, error) {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	// Create container
	hostConfig := &container.HostConfig{
		PortBindings: nat.PortMap{
			"25565/tcp": []nat.PortBinding{
				{
					HostIP:   "0.0.0.0",
					HostPort: "25565",
				},
			},
		},
		Mounts: []mount.Mount{
			{
				Type:   mount.TypeBind,
				Source: cwd + "/plugins",
				Target: "/plugins",
			},
			{
				Type:   mount.TypeBind,
				Source: cwd + "/static",
				Target: "/static",
			},
		},
	}
	c, err := d.ContainerCreate(context.Background(), &container.Config{
		Image:     "dockercraft:latest",
		Tty:       true,
		OpenStdin: true,
	}, hostConfig, nil, nil, "dockercraft_c")
	return c, err
}
