package docker

import (
	"bytes"
	"io"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"

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
	Stdout *bytes.Buffer
	Stderr *bytes.Buffer
}

// Exec spawn a container and wait for its output
func Exec(opts ExecOptions) (*ExecResult, error) {
	cli, err := getClient()
	if err != nil {
		return nil, err
	}

	ctx := context.Background()

	filter := filters.NewArgs()
	filter.Add("label", "belong-to=fx")
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
			Labels:       map[string]string{"belong-to": "fx"},
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

	connOut, err := attach(ctx, cli, 1, containerID)
	if err != nil {
		return nil, err
	}
	connErr, err := attach(ctx, cli, 2, containerID)
	if err != nil {
		return nil, err
	}

	var or bytes.Buffer
	var er bytes.Buffer

	go func() {
		for {
			_, oerr := io.Copy(&or, connOut.Reader)
			if oerr != nil {
				if oerr == io.EOF {
					return
				}
				fmt.Printf("Fail stdout copy: %s\n", oerr.Error())
			}
		}
	}()

	go func() {
		for {
			_, eerr := io.Copy(&er, connErr.Reader)
			if eerr != nil {
				if eerr == io.EOF {
					return
				}
				fmt.Printf("Fail stderr copy: %s\n", eerr.Error())
			}
		}
	}()

	_, err = cli.ContainerWait(ctx, containerID)
	if err != nil {
		return nil, err
	}

	return &ExecResult{
		ID:     containerID,
		Stdout: &or,
		Stderr: &er,
	}, err
}

func attach(ctx context.Context, cli *client.Client, std int, containerID string) (*types.HijackedResponse, error) {

	stdout := false
	stderr := false

	if std == 1 {
		stdout = true
	}
	if std == 2 {
		stderr = true
	}

	attachConfig := types.ContainerAttachOptions{
		Logs:   false,
		Stdin:  false,
		Stdout: stdout,
		Stderr: stderr,
		Stream: true,
	}
	conn, err := cli.ContainerAttach(ctx, containerID, attachConfig)
	if err != nil {
		return nil, err
	}

	return &conn, err
}
