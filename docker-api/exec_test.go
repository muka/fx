package docker

import (
	"log"
	"testing"
)

func TestExec(t *testing.T) {

	opts := ExecOptions{
		Name:       "exec_test_hello",
		ImageName:  "fx/hello",
		Autoremove: true,
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
