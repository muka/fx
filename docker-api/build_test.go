package docker

import (
	"context"
	"os"
	"testing"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
)

func TestBuild(t *testing.T) {
	doBuild(t, "../test/hello", "fx/hello")
}

// Build builds a docker image from the image directory
func doBuild(t *testing.T, srcPath, imageName string) {

	cli, err := getClient()
	if err != nil {
		t.Fatal(err)
	}

	//HACK: for some reason docker complain about the symlink in node_modules/.bin
	dotbin := srcPath + "/node_modules/.bin"
	if _, serr := os.Stat(dotbin); !os.IsNotExist(serr) {
		err = os.RemoveAll(dotbin)
		if err != nil {
			t.Fatal(err)
		}
	}

	err = Build(imageName, srcPath)
	if err != nil {
		t.Fatal(err)
	}

	f := filters.NewArgs()
	f.Add("reference", imageName)
	list, err := cli.ImageList(context.Background(), types.ImageListOptions{
		Filters: f,
	})

	if err != nil {
		t.Fatal(err)
	}

	if len(list) == 0 {
		t.Fatal("Found 0 images")
	}

}
