// Dfdexe embeds and runs a directory of executables from memory. Executables are
// copied from the `exe` directory:
//
//	cp /bin/sh exe
//	cp /bin/ls exe
//	cp /usr/bin/vi exe
//	go build
//	./dfdexe ls -al
package main

import (
	"embed"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"codeberg.org/msantos/execve"

	"golang.org/x/sys/unix"
)

//go:embed exe/*
var exe embed.FS

func main() {
	dir, err := exe.ReadDir("exe")
	if err != nil {
		log.Fatalln("ReadDir:", err)
	}

	fd, err := unix.MemfdCreate("dfdexe", unix.MFD_CLOEXEC)
	if err != nil {
		log.Fatalln("MemfdCreate:", err)
	}

	var bin []byte

	for _, e := range dir {
		if len(os.Args) == 1 {
			fmt.Println(e.Name())
			continue
		}
		if !e.IsDir() && e.Name() == os.Args[1] {
			b, err := exe.ReadFile(filepath.Join("exe", e.Name()))
			if err != nil {
				log.Fatalln("ReadFile:", e.Name(), err)
			}
			bin = b
		}
	}

	if len(os.Args) == 1 || len(bin) == 0 {
		os.Exit(127)
	}

	if n, err := unix.Write(fd, bin); err != nil || n != len(bin) {
		log.Fatalln("Write:", err)
	}

	if err := execve.Fexecve(uintptr(fd), os.Args[1:], os.Environ()); err != nil {
		log.Fatalln("Fexecve:", err)
	}
	os.Exit(126)
}
