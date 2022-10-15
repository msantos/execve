package execve

import (
	"syscall"
	"unsafe"

	"golang.org/x/sys/unix"
)

func execveat(fd uintptr, pathname string, argv []string, envv []string, flags int) error {
	return unix.ENOSYS
}

func fexecve(fd uintptr, argv []string, envv []string) error {
	argvp, err := syscall.SlicePtrFromStrings(argv)
	if err != nil {
		return err
	}

	envvp, err := syscall.SlicePtrFromStrings(envv)
	if err != nil {
		return err
	}

	_, _, err = syscall.Syscall(
		unix.SYS_FEXECVE,
		fd,
		uintptr(unsafe.Pointer(&argvp[0])),
		uintptr(unsafe.Pointer(&envvp[0])),
	)

	return err
}
