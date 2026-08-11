//go:build js || plan9 || wasip1 || windows

package securepath

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// These targets do not expose a portable no-follow *at API through the Go
// standard library. Every pathname component is validated with Lstat and any
// symlink is rejected before the open; the Unix implementation provides the
// race-free descriptor-relative variant where the kernel API is available.
// The residual same-user TOCTOU window here is deliberate and documented.
func OpenDirectory(path string) (*os.File, error) {
	if err := validatePathComponents(path); err != nil {
		return nil, err
	}
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	info, err := file.Stat()
	if err != nil {
		file.Close()
		return nil, err
	}
	if !info.IsDir() {
		file.Close()
		return nil, fmt.Errorf("%s is not a directory", path)
	}
	return file, nil
}

// OpenDirectoryAt opens name inside parent after validating every component.
func OpenDirectoryAt(parent *os.File, name string) (*os.File, error) {
	if parent == nil {
		return nil, errors.New("parent directory is nil")
	}
	return OpenDirectory(filepath.Join(parent.Name(), name))
}

// OpenFileAt opens the regular file name inside parent after validating every
// component.
func OpenFileAt(parent *os.File, name string) (*os.File, error) {
	if parent == nil {
		return nil, errors.New("parent directory is nil")
	}
	path := filepath.Join(parent.Name(), name)
	if err := validatePathComponents(path); err != nil {
		return nil, err
	}
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	info, err := file.Stat()
	if err != nil {
		file.Close()
		return nil, err
	}
	if !info.Mode().IsRegular() {
		file.Close()
		return nil, fmt.Errorf("%s is not a regular file", path)
	}
	return file, nil
}

// OpenFile validates every component of path and opens the leaf.
func OpenFile(path string) (*os.File, error) {
	parent, name, err := OpenParent(path)
	if err != nil {
		return nil, err
	}
	defer parent.Close()
	return OpenFileAt(parent, name)
}

// OpenParent opens path's parent directory and returns it with the leaf name.
func OpenParent(path string) (*os.File, string, error) {
	clean := filepath.Clean(path)
	if clean == "." || clean == ".." || strings.HasSuffix(clean, string(filepath.Separator)) {
		return nil, "", fmt.Errorf("path %q is not a file", path)
	}
	parent, err := OpenDirectory(filepath.Dir(clean))
	if err != nil {
		return nil, "", err
	}
	return parent, filepath.Base(clean), nil
}

// OpenFileRelative resolves relativePath under root after rejecting absolute
// paths and traversal.
func OpenFileRelative(root *os.File, relativePath string) (*os.File, error) {
	clean := filepath.Clean(relativePath)
	if filepath.IsAbs(relativePath) || clean == "." || clean == ".." || strings.HasPrefix(clean, ".."+string(filepath.Separator)) {
		return nil, fmt.Errorf("path %q is not root-relative", relativePath)
	}
	return OpenFile(filepath.Join(root.Name(), clean))
}

// OpenNoFollow validates every component of path and opens the leaf. Without a
// kernel no-follow flag the symlink rejection is a check, not an open-time
// guarantee.
func OpenNoFollow(path string) (*os.File, error) {
	if err := validatePathComponents(path); err != nil {
		return nil, err
	}
	return os.Open(path)
}

func validatePathComponents(path string) error {
	clean := filepath.Clean(path)
	for current := clean; ; current = filepath.Dir(current) {
		info, err := os.Lstat(current)
		if err != nil {
			return err
		}
		if info.Mode()&os.ModeSymlink != 0 {
			return fmt.Errorf("path %s is a symlink", current)
		}
		parent := filepath.Dir(current)
		if parent == current {
			break
		}
	}
	return nil
}
