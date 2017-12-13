package docker

import (
	"io"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"

	"context"
	"fmt"
)

//ExecOptions control how a container is executed
type ExecOptions struct {
	Name             string
	Cmd              []string
	Env              []string
	Stdin            string
	ImageName        string
	RecreateInstance bool
	Context          context.Context
}

//ExecResult return the execution results
type ExecResult struct {
	Stdout string
}

// Exec spawn a container and wait for its output
func Exec(opts ExecOptions) (*ExecResult, error) {
	cli, err := getClient()
	if err != nil {
		return nil, err
	}

	ctx := context.Background()

	containerConfig := &container.Config{
		Cmd:          opts.Cmd,
		Env:          opts.Env,
		Image:        opts.ImageName,
		AttachStdin:  false,
		AttachStderr: false,
		AttachStdout: false,
		Tty:          false,
		StdinOnce:    false,
		Labels:       map[string]string{"belongs-to": "fx"},
	}

	hostConfig := &container.HostConfig{}
	netConfig := &network.NetworkingConfig{}
	resp, err := cli.ContainerCreate(ctx, containerConfig, hostConfig, netConfig, opts.Name)
	if err != nil {
		return nil, err
	}

	if opts.Context != nil {
		ctx = opts.Context
	}

	startConfig := types.ContainerStartOptions{}

	if err = cli.ContainerStart(ctx, resp.ID, startConfig); err != nil {
		return nil, err
	}

	fmt.Println("Started " + resp.ID)

	attachConfig := types.ContainerAttachOptions{
		Stdin:  true,
		Stdout: true,
		Stderr: true,
		Stream: true,
	}
	conn, err := cli.ContainerAttach(ctx, resp.ID, attachConfig)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	fmt.Println("Attached")

	go func() {
		io.Copy(os.Stdout, conn.Reader)
	}()

	//blocks here
	_, err = cli.ContainerWait(ctx, resp.ID)
	if err != nil {
		return nil, err
	}

	out := ""

	fmt.Println(resp.ID)
	return &ExecResult{Stdout: out}, err
}
