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

func TestExecWithTimeout(t *testing.T) {

	imageName := "fx/test-timeout"
	doBuild(t, "../test/timeout", imageName)

	opts := ExecOptions{
		Name:      "exec_test_timeout",
		ImageName: imageName,
		Stdin:     []byte("test"),
		Timeout:   2,
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
