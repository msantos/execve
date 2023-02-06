// Fdexe embeds and runs an executable from memory. The executable must
// be copied to this directory and renamed to "bin":
//
//	cp /usr/bin/busybox bin
//	go build
//	ln -s fdexec ls # busybox uses argv[0]
//	./ls -alh
package main

import (
	_ "embed"
	"log"
	"os"

	"codeberg.org/msantos/execve"

	"golang.org/x/sys/unix"
)

//go:embed bin
var bin []byte

func main() {
	fd, err := unix.MemfdCreate("fdexe", unix.MFD_CLOEXEC)
	if err != nil {
		log.Fatalln("MemfdCreate:", err)
	}

	if n, err := unix.Write(fd, bin); err != nil || n != len(bin) {
		log.Fatalln("Write:", err)
	}

	if err := execve.Fexecve(uintptr(fd), os.Args, os.Environ()); err != nil {
		log.Fatalln("Fexecve:", err)
	}
	os.Exit(126)
}
