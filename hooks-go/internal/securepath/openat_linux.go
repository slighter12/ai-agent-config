//go:build android || linux

package securepath

import "syscall"

func openAt(dirfd int, path string, flags int, mode uint32) (int, error) {
	return syscall.Openat(dirfd, path, flags, mode)
}
