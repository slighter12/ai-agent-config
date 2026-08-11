//go:build openbsd

package securepath

import (
	_ "unsafe"
)

// syscall keeps openat unexported on OpenBSD; use its generated libc bridge.
//
//go:linkname syscallOpenat syscall.openat
func syscallOpenat(dirfd int, path string, flags int, mode uint32) (int, error)

func openAt(dirfd int, path string, flags int, mode uint32) (int, error) {
	return syscallOpenat(dirfd, path, flags, mode)
}
