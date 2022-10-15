// Package execve is a wrapper around the system execve(2) and fexecve(3)
// system calls.
package execve

// Execveat executes at program at a path using a file descriptor.
// The go runtime process image is replaced by the executable described
// by the directory file descriptor and pathname.
func Execveat(fd uintptr, pathname string, argv []string, envv []string, flags int) error {
	return execveat(fd, pathname, argv, envv, flags)
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
// described in NOTES (fexecve(3))
func Fexecve(fd uintptr, argv []string, envv []string) error {
	return fexecve(fd, argv, envv)
}
