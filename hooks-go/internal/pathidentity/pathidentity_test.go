package pathidentity

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func TestResolveDirectoryResolvesSymlink(t *testing.T) {
	root := t.TempDir()
	target := filepath.Join(root, "target")
	alias := filepath.Join(root, "alias")
	if err := os.Mkdir(target, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(target, alias); err != nil {
		t.Fatal(err)
	}

	canonical, info, err := ResolveDirectory(alias)
	if err != nil {
		t.Fatal(err)
	}
	want, err := filepath.EvalSymlinks(target)
	if err != nil {
		t.Fatal(err)
	}
	if canonical != want {
		t.Fatalf("ResolveDirectory() path = %q, want %q", canonical, want)
	}
	if !os.SameFile(info, statResolvedTarget(t, target)) {
		t.Fatal("ResolveDirectory() returned information for the wrong directory")
	}
}

func TestResolveDirectoryRejectsNonDirectory(t *testing.T) {
	path := filepath.Join(t.TempDir(), "file")
	if err := os.WriteFile(path, []byte("not a directory"), 0o644); err != nil {
		t.Fatal(err)
	}

	if _, _, err := ResolveDirectory(path); err == nil {
		t.Fatal("ResolveDirectory() accepted a regular file")
	}
}

func TestResolveDirectoryRejectsMissingPath(t *testing.T) {
	path := filepath.Join(t.TempDir(), "missing")
	if _, _, err := ResolveDirectory(path); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("ResolveDirectory() error = %v, want not-exist error", err)
	}
}

func statResolvedTarget(t *testing.T, path string) os.FileInfo {
	t.Helper()
	info, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	return info
}
