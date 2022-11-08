// Package execve is a wrapper around the system execve(2) and fexecve(3)
// system calls.
package execve

import (
	"syscall"
	"unsafe"

	"golang.org/x/sys/unix"
)

// Execveat executes at program at a path using a file descriptor.
// The go runtime process image is replaced by the executable described
// by the directory file descriptor and pathname.
func Execveat(fd uintptr, pathname string, argv []string, envv []string, flags int) error {
	pathnamep, err := syscall.BytePtrFromString(pathname)
	if err != nil {
		return err
	}

	argvp, err := syscall.SlicePtrFromStrings(argv)
	if err != nil {
		return err
	}

	envvp, err := syscall.SlicePtrFromStrings(envv)
	if err != nil {
		return err
	}

	_, _, err = syscall.Syscall6(
		unix.SYS_EXECVEAT,
		fd,
		uintptr(unsafe.Pointer(pathnamep)),
		uintptr(unsafe.Pointer(&argvp[0])),
		uintptr(unsafe.Pointer(&envvp[0])),
		uintptr(flags),
		0,
	)

	return err
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
//
// BUG(fexecve): If fd refers to a script (i.e., it is an executable
// text file that names a script interpreter with a first line that begins
// with the characters #!) and the close-on-exec flag has been set for fd,
// then fexecve() fails with the error ENOENT.  This error occurs because,
// by  the time  the script interpreter is executed, fd has already been
// closed because of the close-on-exec flag.  Thus, the close-on-exec flag
// can't be set on fd if it refers to a script, leading to the problems
// described in NOTES (fexecve(3)).
func Fexecve(fd uintptr, argv []string, envv []string) error {
	return Execveat(fd, "", argv, envv, unix.AT_EMPTY_PATH)
}
