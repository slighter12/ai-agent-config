//go:build android || darwin || dragonfly || freebsd || ios || linux || netbsd || openbsd

package securepath

import (
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"testing"
)

func TestOpenFileRejectsFIFOWithoutBlocking(t *testing.T) {
	path := filepath.Join(t.TempDir(), "pipe")
	if err := syscall.Mkfifo(path, 0o600); err != nil {
		t.Fatal(err)
	}

	file, err := OpenFile(path)
	if err == nil {
		file.Close()
		t.Fatal("expected FIFO rejection")
	}
	if !strings.Contains(err.Error(), "not a regular file") {
		t.Fatalf("expected non-regular rejection, got %v", err)
	}
}

func TestOpenNoFollowRejectsSymlink(t *testing.T) {
	root := t.TempDir()
	target := filepath.Join(t.TempDir(), "target.txt")
	if err := os.WriteFile(target, []byte("outside"), 0o644); err != nil {
		t.Fatal(err)
	}
	link := filepath.Join(root, "link.txt")
	if err := os.Symlink(target, link); err != nil {
		t.Skipf("symlinks unavailable: %v", err)
	}

	file, err := OpenNoFollow(link)
	if err == nil {
		file.Close()
		t.Fatal("expected symlink rejection")
	}
}
