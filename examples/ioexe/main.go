// Ioexe reads and runs an executable from stdin.
//
//	cat /bin/sh | ./ioexe -c "echo hello world"
package main

import (
	"bufio"
	"io"
	"log"
	"os"

	"codeberg.org/msantos/execve"

	"golang.org/x/sys/unix"
)

func main() {
	flag := unix.MFD_CLOEXEC

	stdin := bufio.NewReader(os.Stdin)
	if p, err := stdin.Peek(2); err == nil && p[0] == '#' && p[1] == '!' {
		flag &= ^unix.MFD_CLOEXEC
	}

	fd, err := unix.MemfdCreate("ioexe", flag)
	if err != nil {
		log.Fatalln("MemfdCreate:", err)
	}

	f := os.NewFile(uintptr(fd), "ioexe")

	if _, err := io.Copy(f, stdin); err != nil {
		log.Fatalln("Copy:", err)
	}

	if err := execve.Fexecve(f.Fd(), os.Args, os.Environ()); err != nil {
		log.Fatalln("Fexecve:", err)
	}
	os.Exit(126)
}
