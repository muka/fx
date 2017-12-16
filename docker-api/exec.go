package docker

import (
	"bytes"
	"io"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/network"

	"context"
	"fmt"
)

var defaultTimeout int64 = 10

//ExecOptions control how a container is executed
type ExecOptions struct {
	Name      string
	Cmd       []string
	Env       []string
	Stdin     []byte
	ImageName string
	// Timeout in second to stop the container
	Timeout int64
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

	// set default TTL
	if opts.Timeout == 0 {
		opts.Timeout = defaultTimeout
	}

	var cmd []string
	if len(opts.Cmd) == 0 {
		imgfilter := filters.NewArgs()
		imgfilter.Add("reference", opts.ImageName)
		listOptions := types.ImageListOptions{
			All:     true,
			Filters: imgfilter,
		}
		var imageID string
		imageList, ilerr := cli.ImageList(ctx, listOptions)
		if ilerr != nil {
			return nil, ilerr
		}
		if len(imageList) == 1 {
			imageID = imageList[0].ID
		} else {
			return nil, fmt.Errorf("Image not found %s", opts.ImageName)
		}

		imageInfo, _, iierr := cli.ImageInspectWithRaw(ctx, imageID)
		if iierr != nil {
			return nil, iierr
		}

		cmd = append(imageInfo.Config.Cmd, string(opts.Stdin))
	} else {
		cmd = append(opts.Cmd, string(opts.Stdin))
	}

	fmt.Printf("Creating container %s (from %s)\n", opts.Name, opts.ImageName)
	containerConfig := &container.Config{
		Cmd:          cmd,
		Env:          opts.Env,
		Image:        opts.ImageName,
		AttachStdin:  false,
		AttachStderr: true,
		AttachStdout: true,
		Tty:          true,
		StdinOnce:    true,
		Labels: map[string]string{
			"belong-to": "fx",
			"fx-type":   "container",
		},
	}

	hostConfig := &container.HostConfig{
		AutoRemove: true,
	}
	netConfig := &network.NetworkingConfig{}
	resp, cerr := cli.ContainerCreate(ctx, containerConfig, hostConfig, netConfig, opts.Name)
	if cerr != nil {
		return nil, cerr
	}

	containerID := resp.ID

	attachConfig := types.ContainerAttachOptions{
		Logs:   false,
		Stdin:  false,
		Stdout: true,
		Stderr: true,
		Stream: true,
	}
	conn, err := cli.ContainerAttach(ctx, containerID, attachConfig)
	if err != nil {
		return nil, err
	}
	// defer conn.Close()

	startConfig := types.ContainerStartOptions{}
	if err = cli.ContainerStart(ctx, containerID, startConfig); err != nil {
		return nil, err
	}

	fmt.Printf("Started %s\n", containerID)

	var outBuffer bytes.Buffer
	go func() {
		for {
			_, outBuffer := io.Copy(&outBuffer, conn.Reader)
			if outBuffer != nil {
				if outBuffer == io.EOF {
					return
				}
				fmt.Printf("Fail stdout copy: %s\n", outBuffer.Error())
			}
		}
	}()

	wait := make(chan bool, 1)

	go func() {
		d := time.Second * time.Duration(opts.Timeout)
		time.Sleep(d)
		kerr := Kill(containerID)
		if kerr != nil {
			fmt.Printf("Error on kill: %s", kerr.Error())
		}
		wait <- true
	}()

	go func() {
		ch := getEventsChannel()
		for {
			select {
			case ev := <-ch:
				// fmt.Printf("%++v", ev)
				if ev.ID == containerID {
					if ev.Action == "die" {
						wait <- true
						return
					}
				}
			}
		}
	}()

	// sleep until something happens
	<-wait

	return &ExecResult{
		ID:     containerID,
		Stdout: &outBuffer,
		Stderr: nil,
	}, err
}
