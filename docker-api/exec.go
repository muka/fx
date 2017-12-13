package docker

import (
	"io"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/network"

	"context"
	"fmt"
)

//ExecOptions control how a container is executed
type ExecOptions struct {
	Name       string
	Cmd        []string
	Env        []string
	Stdin      string
	ImageName  string
	Autoremove bool
	Context    context.Context
}

//ExecResult return the execution results
type ExecResult struct {
	ID     string
	Stdout string
}

// Exec spawn a container and wait for its output
func Exec(opts ExecOptions) (*ExecResult, error) {
	cli, err := getClient()
	if err != nil {
		return nil, err
	}

	ctx := context.Background()

	filter := filters.NewArgs()
	filter.Add("label", "belongs-to=fx")
	filter.Add("name", opts.Name)
	list, err := cli.ContainerList(ctx, types.ContainerListOptions{
		Filters: filter,
	})
	if err != nil {
		return nil, err
	}

	var containerID string
	if len(list) == 0 {

		fmt.Printf("Creating container %s (from %s)\n", opts.Name, opts.ImageName)

		containerConfig := &container.Config{
			Cmd:          opts.Cmd,
			Env:          opts.Env,
			Image:        opts.ImageName,
			AttachStdin:  true,
			AttachStderr: true,
			AttachStdout: true,
			Tty:          true,
			StdinOnce:    true,
			Labels:       map[string]string{"belongs-to": "fx"},
		}

		hostConfig := &container.HostConfig{
			AutoRemove: opts.Autoremove,
		}
		netConfig := &network.NetworkingConfig{}
		resp, cerr := cli.ContainerCreate(ctx, containerConfig, hostConfig, netConfig, opts.Name)
		if cerr != nil {
			return nil, cerr
		}

		containerID = resp.ID

	} else {
		containerID = list[0].ID
	}

	if opts.Context != nil {
		ctx = opts.Context
	}

	startConfig := types.ContainerStartOptions{}
	if err = cli.ContainerStart(ctx, containerID, startConfig); err != nil {
		return nil, err
	}
	fmt.Println("Started " + containerID)

	attachConfig := types.ContainerAttachOptions{
		Logs:   true,
		Stdin:  true,
		Stdout: true,
		Stderr: true,
		Stream: true,
	}
	conn, err := cli.ContainerAttach(ctx, containerID, attachConfig)
	if err != nil {
		return nil, err
	}
	// defer conn.Close()
	fmt.Println("Attached")

	go func() {
		for {
			_, serr := io.Copy(os.Stdout, conn.Reader)
			if serr != nil {
				fmt.Printf("Fail copy: %s\n", serr.Error())
				return
			}
		}

	}()

	// logsConfig := types.ContainerLogsOptions{
	// 	ShowStdout: true,
	// }
	// out, err := cli.ContainerLogs(ctx, containerID, logsConfig)
	// if err != nil {
	// 	return nil, err
	// }

	// _, err = cli.ContainerWait(ctx, containerID)
	// if err != nil {
	// 	return nil, err
	// }

	return &ExecResult{
		ID:     containerID,
		Stdout: "",
	}, err
}
