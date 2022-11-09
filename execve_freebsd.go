// Package execve is a wrapper around the system fexecve(3) system call.
package execve

import (
	"syscall"
	"unsafe"

	"golang.org/x/sys/unix"
)

// Execveat is not supported on this platform.
func Execveat(fd uintptr, pathname string, argv []string, envv []string, flags int) error {
	return unix.ENOSYS
}

// Fexecve executes the program referred to by a file descriptor.
// The go runtime process image is replaced with the executable referred
// to by the file descriptor.
// The file descriptor should be opened with the O_CLOEXEC flag set
// (the default when using os.Open) to prevent the fd from leaking to the
// new process image.
//
// The exception to this rule is running scripts: the file descriptor
// must be opened without O_CLOEXEC.
func Fexecve(fd uintptr, argv []string, envv []string) error {
	argvp, err := syscall.SlicePtrFromStrings(argv)
	if err != nil {
		return err
	}

	envvp, err := syscall.SlicePtrFromStrings(envv)
	if err != nil {
		return err
	}

	_, _, errno := syscall.Syscall(
		unix.SYS_FEXECVE,
		fd,
		uintptr(unsafe.Pointer(&argvp[0])),
		uintptr(unsafe.Pointer(&envvp[0])),
	)

	return errno
}
