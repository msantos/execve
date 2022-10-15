package execve

import (
	"syscall"
	"unsafe"

	"golang.org/x/sys/unix"
)

func execveat(fd uintptr, pathname string, argv []string, envv []string, flags int) error {
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

func fexecve(fd uintptr, argv []string, envv []string) error {
	return execveat(fd, "", argv, envv, unix.AT_EMPTY_PATH)
}
