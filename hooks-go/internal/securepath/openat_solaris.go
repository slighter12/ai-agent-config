//go:build solaris

package securepath

import "errors"

// Solaris's syscall package does not expose an openat bridge. Failing closed
// avoids falling back to pathname traversal on this target.
func openAt(dirfd int, path string, flags int, mode uint32) (int, error) {
	return 0, errors.New("descriptor-relative traversal is unavailable on Solaris")
}
