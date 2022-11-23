# execve

[![Go Reference](https://pkg.go.dev/badge/codeberg.org/msantos/execve.svg)](https://pkg.go.dev/codeberg.org/msantos/execve)

A Go package for fexecve(3) and execveat(2).

# EXAMPLES

* [fdexe](examples/fdexe/main.go): embed an executable in a Go binary
  and run from memory

* [dfdexe](examples/dfdexe/main.go): embed a directory of executables
  in a Go binary and run from memory

* [ioexe](examples/iodexe/main.go): read an executable from stdin and
  run from memory

## Run an executable using a file descriptor

```go
package main

import (
	"fmt"
	"os"
	"os/exec"

	"codeberg.org/msantos/execve"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "usage: <cmd> <args>\n")
		os.Exit(2)
	}

	arg0, err := exec.LookPath(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(127)
	}

	fd, err := os.Open(arg0)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if err := execve.Fexecve(fd.Fd(), os.Args[1:], os.Environ()); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	os.Exit(126)
}
```

## Execute a script using a file descriptor

```go
package main

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"

	"codeberg.org/msantos/execve"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "usage: <cmd> <args>\n")
		os.Exit(2)
	}

	arg0, err := exec.LookPath(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(127)
	}

	fd, err := syscall.Open(arg0, syscall.O_RDONLY, 0)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	os.Args[1] = fmt.Sprintf("/dev/fd/%d", fd)

	if err := execve.Fexecve(uintptr(fd), os.Args[1:], os.Environ()); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	os.Exit(126)
}
```
