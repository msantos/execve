package execve_test

import (
	"fmt"
	"io"
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

func TestFexecveMemfd(t *testing.T) {
	if os.Getenv("TESTING_EXECVE_TESTFEXECVEMEMFD") == "1" {
		ExampleFexecve_memfd()
		return
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestFexecveMemfd")
	cmd.Env = append(os.Environ(), "TESTING_EXECVE_TESTFEXECVEMEMFD=1")

	if err := run(cmd, "test"); err != nil {
		t.Errorf("%v", err)
		return
	}
}

// ExampleFexecve_memfd is an example of running an executable from memory.
func ExampleFexecve_memfd() {
	path, err := exec.LookPath("echo")
	if err != nil {
		fmt.Fprintf(os.Stderr, "LookPath: %v", err)
		return
	}

	exe, err := os.Open(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Open: %v", err)
		return
	}

	fd, err := unix.MemfdCreate("ExampleFexecve_memfd", unix.MFD_CLOEXEC)
	if err != nil {
		fmt.Fprintf(os.Stderr, "MemfdCreate: %v", err)
		return
	}

	f := os.NewFile(uintptr(fd), "ExampleFexecve_memfd")

	if _, err := io.Copy(f, exe); err != nil {
		fmt.Fprintf(os.Stderr, "Copy: %v", err)
		return
	}

	if err := execve.Fexecve(f.Fd(), []string{"-n", "test"}, os.Environ()); err != nil {
		fmt.Fprintf(os.Stderr, "Fexecve: %v", err)
		return
	}
}
