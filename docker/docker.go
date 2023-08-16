package docker

import (
	"bufio"
	"context"
	"os"
	"strings"

	dtypes "github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/drewharris/shulker/commands"
	"github.com/drewharris/shulker/types"
)

func CreateContainer(d *client.Client) (container.CreateResponse, error) {
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
				Target: "/plugins/plugins",
			},
			{
				Type:   mount.TypeBind,
				Source: cwd + "/libraries",
				Target: "/plugins/libraries",
			},
			{
				Type:   mount.TypeBind,
				Source: cwd + "/static",
				Target: "/static",
			},
		},
	}
	c, err := d.ContainerCreate(context.Background(), &container.Config{
		Image:     "shulker:latest",
		Tty:       true,
		OpenStdin: true,
	}, hostConfig, nil, nil, "shulker_c")
	return c, err
}

func RunContainerCommand(d *client.Client, cid string, log types.Logger, cmd commands.Command) (string, error) {
	var fullCmd []string
	fullCmd = append(fullCmd, cmd.Name)
	fullCmd = append(fullCmd, cmd.Args...)

	execId, err := d.ContainerExecCreate(context.TODO(), cid, dtypes.ExecConfig{
		Cmd:          fullCmd,
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
		return "", err
	}

	defer rd.Close()
	scanner := bufio.NewScanner(rd.Reader) // Scanner doesn't return newline byte
	for scanner.Scan() {
		log(strings.ReplaceAll(scanner.Text(), "\n", ""))
	}
	// After command is done, run callback

	return execId.ID, nil
}
