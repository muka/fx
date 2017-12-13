package docker

import (
	"testing"
	"time"
)

func TestExec(t *testing.T) {

	opts := ExecOptions{
		Name:       "exec_test",
		ImageName:  "willfarrell/ping",
		Env:        []string{"HOSTNAME=localhost", "TIMEOUT=1"},
		Stdin:      "test",
		Autoremove: true,
	}

	res, err := Exec(opts)
	if err != nil {
		t.Fatalf("Exec failed: %s", err.Error())
	}

	time.Sleep(time.Second * 5)

	err = Kill(res.ID)
	if err != nil {
		t.Fatalf("Exec failed: %s", err.Error())
	}

}
