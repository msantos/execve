package execve_test

import (
	"fmt"
	"os"
	"os/exec"
	"testing"

	"codeberg.org/msantos/execve"
	"golang.org/x/sys/unix"
)

func TestExecveat(t *testing.T) {
	if os.Getenv("TESTING_EXECVE_TESTEXECVEAT") == "1" {
		ExampleExecveat()
		return
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestExecveat")
	cmd.Env = append(os.Environ(), "TESTING_EXECVE_TESTEXECVEAT=1")

	if err := run(cmd, "test\n"); err != nil {
		t.Errorf("%v", err)
		return
	}
}

func ExampleExecveat() {
	if err := os.Chdir("/bin"); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	n := int32(unix.AT_FDCWD)
	AT_FDCWD := uintptr(uint32(n))

	if err := execve.Execveat(
		AT_FDCWD,
		"sh", []string{"sh", "-c", "echo test"},
		[]string{},
		0,
	); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
}
