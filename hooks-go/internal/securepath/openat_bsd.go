//go:build dragonfly || freebsd || netbsd

package securepath

import (
	"syscall"
	"unsafe"
)

func openAt(dirfd int, path string, flags int, mode uint32) (int, error) {
	p, err := syscall.BytePtrFromString(path)
	if err != nil {
		return 0, err
	}
	r0, _, errno := syscall.Syscall6(syscall.SYS_OPENAT, uintptr(dirfd), uintptr(unsafe.Pointer(p)), uintptr(flags), uintptr(mode), 0, 0)
	if errno != 0 {
		return 0, errno
	}
	return int(r0), nil
}
