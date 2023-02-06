// Dfdexe embeds and runs a directory of executables from memory. Executables are
// copied from the "bin" directory:
//
//	cp /bin/sh bin
//	cp /bin/ls bin
//	cp /usr/bin/vi bin
//	go build
//	./dfdexe bin/ls -al
package main

import (
	"embed"
	"log"
	"os"

	"codeberg.org/msantos/execve"

	"golang.org/x/sys/unix"
)

//go:embed bin/*
var bin embed.FS

func main() {
	bin, err := bin.ReadFile(os.Args[1])
	if err != nil {
		log.Fatalln("ReadFile:", err)
	}

	fd, err := unix.MemfdCreate("dfdexe", unix.MFD_CLOEXEC)
	if err != nil {
		log.Fatalln("MemfdCreate:", err)
	}

	if n, err := unix.Write(fd, bin); err != nil || n != len(bin) {
		log.Fatalln("Write:", err)
	}

	if err := execve.Fexecve(uintptr(fd), os.Args[1:], os.Environ()); err != nil {
		log.Fatalln("Fexecve:", err)
	}
	os.Exit(126)
}
