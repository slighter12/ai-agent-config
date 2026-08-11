//go:build aix || android || darwin || dragonfly || freebsd || ios || linux || netbsd || openbsd || solaris

package securepath

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"syscall"

	"ai-agent-config/hooks/internal/pathidentity"
)

// OpenDirectory opens every component from a held starting directory. It never
// asks the kernel to resolve the complete caller-supplied path, so an ancestor
// symlink or a concurrent directory rename cannot redirect traversal.
func OpenDirectory(path string) (*os.File, error) {
	clean := filepath.Clean(path)
	if clean == "" {
		return nil, errors.New("directory path is empty")
	}
	clean = pathidentity.NormalizeSystemAlias(clean)

	startPath := "."
	components := strings.Split(clean, string(filepath.Separator))
	if filepath.IsAbs(clean) {
		startPath = string(filepath.Separator)
		components = strings.Split(strings.TrimPrefix(clean, string(filepath.Separator)), string(filepath.Separator))
	}
	current, err := openDirectoryBase(startPath)
	if err != nil {
		return nil, err
	}
	for _, component := range components {
		if component == "" || component == "." {
			continue
		}
		next, err := OpenDirectoryAt(current, component)
		if err != nil {
			current.Close()
			return nil, fmt.Errorf("open directory component %s: %w", component, err)
		}
		current.Close()
		current = next
	}
	return current, nil
}

func openDirectoryBase(path string) (*os.File, error) {
	file, err := os.OpenFile(path, os.O_RDONLY|syscall.O_DIRECTORY|syscall.O_NOFOLLOW, 0)
	if err != nil {
		return nil, err
	}
	return file, nil
}

// OpenDirectoryAt opens name inside parent without following a symlink.
func OpenDirectoryAt(parent *os.File, name string) (*os.File, error) {
	if parent == nil {
		return nil, errors.New("parent directory is nil")
	}
	fd, err := openAt(int(parent.Fd()), name, os.O_RDONLY|syscall.O_DIRECTORY|syscall.O_NOFOLLOW, 0)
	if err != nil {
		return nil, err
	}
	file := os.NewFile(uintptr(fd), filepath.Join(parent.Name(), name))
	if file == nil {
		_ = syscall.Close(fd)
		return nil, errors.New("could not create directory handle")
	}
	info, err := file.Stat()
	if err != nil {
		file.Close()
		return nil, err
	}
	if !info.IsDir() {
		file.Close()
		return nil, fmt.Errorf("%s is not a directory", name)
	}
	return file, nil
}

// OpenFileAt opens the regular file name inside parent without following a
// symlink.
func OpenFileAt(parent *os.File, name string) (*os.File, error) {
	if parent == nil {
		return nil, errors.New("parent directory is nil")
	}
	// O_NONBLOCK prevents a FIFO or device from blocking before fstat can
	// reject it as a non-regular input.
	fd, err := openAt(int(parent.Fd()), name, os.O_RDONLY|syscall.O_NOFOLLOW|syscall.O_NONBLOCK, 0)
	if err != nil {
		return nil, err
	}
	file := os.NewFile(uintptr(fd), filepath.Join(parent.Name(), name))
	if file == nil {
		_ = syscall.Close(fd)
		return nil, errors.New("could not create file handle")
	}
	info, err := file.Stat()
	if err != nil {
		file.Close()
		return nil, err
	}
	if !info.Mode().IsRegular() {
		file.Close()
		return nil, fmt.Errorf("%s is not a regular file", name)
	}
	return file, nil
}

// OpenFile opens path's parent by component-wise traversal, then opens the
// leaf without following a symlink.
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

// OpenFileRelative walks relativePath from the held root descriptor. The path
// must stay inside root: absolute paths and traversal are rejected.
func OpenFileRelative(root *os.File, relativePath string) (*os.File, error) {
	clean := filepath.Clean(relativePath)
	if filepath.IsAbs(relativePath) || clean == "." || clean == ".." || strings.HasPrefix(clean, ".."+string(filepath.Separator)) {
		return nil, fmt.Errorf("path %q is not root-relative", relativePath)
	}
	components := strings.Split(clean, string(filepath.Separator))
	if len(components) == 0 || components[len(components)-1] == "" {
		return nil, fmt.Errorf("path %q is not a file", relativePath)
	}
	current := root
	owned := false
	defer func() {
		if owned {
			current.Close()
		}
	}()
	for _, component := range components[:len(components)-1] {
		if component == "" || component == "." {
			continue
		}
		next, err := OpenDirectoryAt(current, component)
		if err != nil {
			return nil, fmt.Errorf("open parent %s: %w", component, err)
		}
		if owned {
			current.Close()
		}
		current = next
		owned = true
	}
	return OpenFileAt(current, components[len(components)-1])
}

// OpenNoFollow opens a single file without following a symlink at the leaf. It
// performs no component traversal; callers that need ancestor guarantees must
// use OpenFile or OpenFileRelative.
func OpenNoFollow(path string) (*os.File, error) {
	return os.OpenFile(path, os.O_RDONLY|syscall.O_NOFOLLOW|syscall.O_NONBLOCK, 0)
}
