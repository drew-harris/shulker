package engine

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	dtypes "github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/drewharris/shulker/commands"
	"github.com/drewharris/shulker/docker"
	"github.com/drewharris/shulker/styles"
	"github.com/drewharris/shulker/types"
)

type DockerEngine struct {
	client           *client.Client
	spigotConnection dtypes.HijackedResponse
	containerId      string
	imageId          string

	execs []string
}

func NewDockerEngine(client *client.Client) *DockerEngine {
	return &DockerEngine{
		client: client,
	}
}

func (e *DockerEngine) CanAttach() bool {
	return true
}

func (e *DockerEngine) EnsureSetup(sub chan types.OutputMsg) error {
	images, err := e.client.ImageList(context.Background(), dtypes.ImageListOptions{All: true})
	if err != nil {
		fmt.Println(styles.Error.Render("\n   Can't connect to Docker, Is it running?"))
		os.Exit(1)
	}

	var imageId = ""
	for _, image := range images {
		tag := image.RepoTags[0]
		if tag == "shulker:latest" {
			imageId = image.ID
		}
	}

	// - Check if image exists -> build
	if imageId == "" { // If there is no image
		err := commands.RunExternalCommand(sub, commands.Command{
			Name:   "docker",
			Args:   []string{"build", "-t", "shulker", "-f", "Dockerfile.dev", "."},
			Target: types.StartupOutput,
		})
		if err != nil {
			panic(err)
		}
	}

	// - Check if container exists -> create
	containers, err := e.client.ContainerList(context.Background(), dtypes.ContainerListOptions{
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
		c, err := docker.CreateContainer(e.client)
		if err != nil {
			panic(err)
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

	e.containerId = containerId

	// - Start container
	err = e.client.ContainerStart(context.Background(), containerId, dtypes.ContainerStartOptions{})
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
	_, err = e.client.ContainerStatPath(context.Background(), containerId, "/server/plugins/HuMInGameLabsPlugin.jar")
	errors = append(errors, err)
	_, err = e.client.ContainerStatPath(context.Background(), containerId, "/server/plugins/Contraption.jar")
	errors = append(errors, err)
	_, err = e.client.ContainerStatPath(context.Background(), containerId, "/server/plugins/Recycler.jar")
	errors = append(errors, err)
	_, err = e.client.ContainerStatPath(context.Background(), containerId, "/server/plugins/SplitterNode.jar")
	errors = append(errors, err)

	var notBuilt bool = false
	for _, err := range errors {
		if err != nil {
			notBuilt = true
			break
		}
	}

	if notBuilt {
		execId, err := e.client.ContainerExecCreate(context.TODO(), containerId, dtypes.ExecConfig{
			Cmd:          []string{"./build_all.sh"},
			AttachStderr: true,
			AttachStdout: true,
			Tty:          true,
		})
		if err != nil {
			panic(err)
		}

		rd, err := e.client.ContainerExecAttach(context.Background(), execId.ID, dtypes.ExecStartCheck{})
		e.client.ContainerExecStart(context.Background(), execId.ID, dtypes.ExecStartCheck{
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
		Message: "All Plugins Built!",
	}

	return nil
}

func (e *DockerEngine) StartServer(sub chan types.OutputMsg) error {
	waiter, err := e.client.ContainerAttach(context.Background(), e.containerId, dtypes.ContainerAttachOptions{
		Stderr: true,
		Stdout: true,
		Stdin:  true,
		Stream: true,
	})
	if err != nil {
		panic("Couldn't attach to container " + err.Error() + "Container id: " + e.containerId)
	}

	go func() {
		scanner := bufio.NewScanner(waiter.Reader)
		for scanner.Scan() {
			sub <- types.OutputMsg{
				Target:  types.ServerOutput,
				Message: scanner.Text(),
			}
		}
	}()

	waiter.Conn.Write([]byte("./start.sh\n"))

	e.spigotConnection = waiter

	return nil
}

func (e *DockerEngine) Shutdown() error {
	e.spigotConnection.Close()
	err := e.client.ContainerStop(context.Background(), e.containerId, container.StopOptions{})
	return err
}

func (e *DockerEngine) RebuildAllPlugins(sub chan types.OutputMsg) error {
	result, err := docker.RunContainerCommand(e.client, e.containerId, sub, commands.Command{
		Target: types.BuildOutput,
		Name:   "./build_all.sh",
	})

	if err != nil {
		return err
	} else {
		e.execs = append(e.execs, result)
		return nil
	}
}
