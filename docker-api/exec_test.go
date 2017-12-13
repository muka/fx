package docker

import (
	"testing"
)

func TestExec(t *testing.T) {

	opts := ExecOptions{
		Name:             "exec_test",
		ImageName:        "willfarrell/ping",
		Env:              []string{"HOSTNAME=localhost", "TIMEOUT:100"},
		Stdin:            "test",
		RecreateInstance: true,
	}

	_, err := Exec(opts)
	if err != nil {
		t.Fatalf("Exec failed: %s", err.Error())
	}

}
