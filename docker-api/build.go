package docker

import (
	"bufio"
	"encoding/json"

	"github.com/docker/docker/api/types"
	"github.com/jhoonb/archivex"

	"context"
	"fmt"
	"log"
	"os"
)

// Build builds a docker image from the image directory
func Build(name string, dir string) error {
	cli, err := getClient()
	if err != nil {
		return err
	}

	tar := new(archivex.TarFile)
	err = tar.Create(dir)
	if err != nil {
		return err
	}
	err = tar.AddAll(dir, false)
	if err != nil {
		return err
	}
	err = tar.Close()
	if err != nil {
		return err
	}

	dockerBuildContext, buildContextErr := os.Open(dir + ".tar")
	if buildContextErr != nil {
		return buildContextErr
	}
	defer dockerBuildContext.Close()

	buildOptions := types.ImageBuildOptions{
		Dockerfile: "Dockerfile", // optional, is the default
		Tags:       []string{name},
		Labels:     map[string]string{"belong-to": "fx"},
	}
	buildResponse, buildErr := cli.ImageBuild(context.Background(), dockerBuildContext, buildOptions)
	if buildErr != nil {
		return buildErr
	}
	log.Println("build", buildResponse.OSType)
	defer buildResponse.Body.Close()

	scanner := bufio.NewScanner(buildResponse.Body)
	for scanner.Scan() {
		var info dockerInfo
		err := json.Unmarshal(scanner.Bytes(), &info)
		if err != nil {
			return err
		}
		fmt.Printf(info.Stream)
	}

	return nil
}
