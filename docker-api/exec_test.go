package docker

import (
	"log"
	"strings"
	"testing"
)

func TestExecNoArgs(t *testing.T) {

	TestBuild(t)

	opts := ExecOptions{
		Name:      "exec_test_hello",
		ImageName: "fx/hello",
	}

	res, err := Exec(opts)
	if err != nil {
		t.Fatalf("Exec failed: %s", err.Error())
	}

	if len(res.Stdout.String()) == 0 {
		t.Fatal("Unexpected empty output")
	}

	log.Printf("Out: \n\n%s", res.Stdout.String())
}

func TestExecWithArgs(t *testing.T) {

	TestBuild(t)

	hello := "world"
	opts := ExecOptions{
		Name:      "exec_test_hello",
		ImageName: "fx/hello",
		Stdin:     []byte(hello),
	}

	res, err := Exec(opts)
	if err != nil {
		t.Fatalf("Exec failed: %s", err.Error())
	}

	if len(res.Stdout.String()) == 0 {
		t.Fatal("Unexpected empty output")
	}

	log.Printf("Out: \n\n%s", res.Stdout.String())

	lines := strings.Split(res.Stdout.String(), "\n")
	if !strings.Contains(lines[0], hello) {
		t.Fatalf("Expecting to find `hello %s` at first line", hello)
	}

}

// func TestExecTimeout(t *testing.T) {
//
// 	ctx1 := context.Background()
// 	ctx, cancel := context.WithCancel(ctx1)
//
// 	opts := ExecOptions{
// 		Name:       "exec_test",
// 		ImageName:  "willfarrell/ping",
// 		Env:        []string{"HOSTNAME=localhost", "TIMEOUT=1"},
// 		Autoremove: true,
// 		Context:    ctx,
// 	}
//
// 	res, err := Exec(opts)
// 	if err != nil {
// 		t.Fatalf("Exec failed: %s", err.Error())
// 	}
// 	defer res.Close()
//
// 	time.Sleep(time.Millisecond * 500)
// 	cancel()
//
// 	err = Kill(res.ID)
// 	if err != nil {
// 		t.Fatalf("Exec failed: %s", err.Error())
// 	}
//
// }
