//go:build !linux && !freebsd
// +build !linux,!freebsd

// Package execve is a wrapper around the system execve(2) and fexecve(3)
// system calls.
package execve

import (
	"golang.org/x/sys/unix"
)

// Execveat is not supported on this platform.
func Execveat(fd uintptr, pathname string, argv []string, envv []string, flags int) error {
	return unix.ENOSYS
}

// Fexecve is not supported on this platform.
func Fexecve(fd uintptr, argv []string, envv []string) error {
	return unix.ENOSYS
}
