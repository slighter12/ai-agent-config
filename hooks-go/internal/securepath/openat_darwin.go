//go:build darwin || ios

package securepath

import _ "unsafe"

// syscall keeps openat unexported on Darwin even though its libc entry point
// is available. This narrow bridge matches the standard library's Unix helper.
//
//go:linkname syscallOpenat syscall.openat
func syscallOpenat(dirfd int, path string, flags int, mode uint32) (int, error)

func openAt(dirfd int, path string, flags int, mode uint32) (int, error) {
	return syscallOpenat(dirfd, path, flags, mode)
}
