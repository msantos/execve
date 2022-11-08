package execve_test

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"testing"

	"codeberg.org/msantos/execve"
)

var errInvalidOutput = errors.New("unexpected output")

func run(cmd *exec.Cmd, output string) error {
	var buf bytes.Buffer

	cmd.Stdout = &buf
	cmd.Stderr = &buf
	cmd.Stdin = os.Stdin

	if err := cmd.Run(); err != nil {
		return err
	}

	if !strings.HasPrefix(buf.String(), output) {
		return fmt.Errorf("Expected: %s\nOutput: %s\nError: %w",
			output,
			buf.String(),
			errInvalidOutput,
		)
	}

	return nil
}

func TestFexecve(t *testing.T) {
	if os.Getenv("TESTING_EXECVE_TESTFEXECVE") == "1" {
		ExampleFexecve()
		return
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestFexecve")
	cmd.Env = append(os.Environ(), "TESTING_EXECVE_TESTFEXECVE=1")

	if err := run(cmd, "test"); err != nil {
		t.Errorf("%v", err)
		return
	}
}

func ExampleFexecve() {
	fd, err := os.Open("/bin/sh")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	if err := execve.Fexecve(fd.Fd(), []string{"sh", "-c", "echo test"}, []string{}); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
}

func TestFexecveScript(t *testing.T) {
	if os.Getenv("TESTING_EXECVE_TESTFEXECVESCRIPT") == "1" {
		ExampleFexecve_script()
		return
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestFexecveScript")
	cmd.Env = append(os.Environ(), "TESTING_EXECVE_TESTFEXECVESCRIPT=1")

	if err := run(cmd, "test"); err != nil {
		t.Errorf("%v", err)
		return
	}
}

func ExampleFexecve_script() {
	// Make a shell script
	file, err := os.CreateTemp("", "script")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	name := file.Name()

	if _, err := file.Write([]byte("#!/bin/sh\necho $@")); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	if err := file.Chmod(0o755); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	if err := file.Close(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	// Get a file descriptor with O_CLOEXEC unset
	fd, err := syscall.Open(name, syscall.O_RDONLY, 0)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	sh := fmt.Sprintf("/dev/fd/%d", fd)

	if err := execve.Fexecve(uintptr(fd), []string{sh, "test"}, []string{}); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
}
