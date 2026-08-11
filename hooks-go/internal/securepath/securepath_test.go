package securepath

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestOpenDirectoryRejectsSymlinkedComponent(t *testing.T) {
	root := t.TempDir()
	outside := filepath.Join(t.TempDir(), "outside")
	if err := os.Mkdir(outside, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(outside, filepath.Join(root, "link")); err != nil {
		t.Skipf("symlinks unavailable: %v", err)
	}

	file, err := OpenDirectory(filepath.Join(root, "link"))
	if err == nil {
		file.Close()
		t.Fatal("expected symlinked component rejection")
	}
}

func TestOpenFileRejectsSymlinkedLeaf(t *testing.T) {
	root := t.TempDir()
	target := filepath.Join(t.TempDir(), "target.txt")
	if err := os.WriteFile(target, []byte("outside"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(target, filepath.Join(root, "leaf.txt")); err != nil {
		t.Skipf("symlinks unavailable: %v", err)
	}

	file, err := OpenFile(filepath.Join(root, "leaf.txt"))
	if err == nil {
		file.Close()
		t.Fatal("expected symlinked leaf rejection")
	}
}

func TestOpenFileReadsRegularFile(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, "plain.txt")
	if err := os.WriteFile(path, []byte("inside"), 0o644); err != nil {
		t.Fatal(err)
	}

	file, err := OpenFile(path)
	if err != nil {
		t.Fatalf("expected regular file to open, got %v", err)
	}
	defer file.Close()
	info, err := file.Stat()
	if err != nil {
		t.Fatal(err)
	}
	if !info.Mode().IsRegular() {
		t.Fatalf("expected regular file, got %v", info.Mode())
	}
}

func TestOpenFileRelativeRejectsEscape(t *testing.T) {
	root := t.TempDir()
	directory, err := OpenDirectory(root)
	if err != nil {
		t.Fatal(err)
	}
	defer directory.Close()

	for _, relativePath := range []string{"../escape.txt", "/etc/hosts", "..", "."} {
		file, err := OpenFileRelative(directory, relativePath)
		if err == nil {
			file.Close()
			t.Fatalf("expected %q to be rejected", relativePath)
		}
		if !strings.Contains(err.Error(), "root-relative") && !strings.Contains(err.Error(), "not a file") {
			t.Fatalf("unexpected error for %q: %v", relativePath, err)
		}
	}
}

func TestOpenFileRelativeWalksNestedPath(t *testing.T) {
	root := t.TempDir()
	nested := filepath.Join(root, "a", "b")
	if err := os.MkdirAll(nested, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(nested, "leaf.txt"), []byte("nested"), 0o644); err != nil {
		t.Fatal(err)
	}
	directory, err := OpenDirectory(root)
	if err != nil {
		t.Fatal(err)
	}
	defer directory.Close()

	file, err := OpenFileRelative(directory, filepath.Join("a", "b", "leaf.txt"))
	if err != nil {
		t.Fatalf("expected nested path to open, got %v", err)
	}
	file.Close()
}

func TestOpenFileRelativeRejectsSymlinkedIntermediate(t *testing.T) {
	root := t.TempDir()
	outside := filepath.Join(t.TempDir(), "outside")
	if err := os.Mkdir(outside, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(outside, "leaf.txt"), []byte("outside"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(outside, filepath.Join(root, "link")); err != nil {
		t.Skipf("symlinks unavailable: %v", err)
	}
	directory, err := OpenDirectory(root)
	if err != nil {
		t.Fatal(err)
	}
	defer directory.Close()

	file, err := OpenFileRelative(directory, filepath.Join("link", "leaf.txt"))
	if err == nil {
		file.Close()
		t.Fatal("expected symlinked intermediate rejection")
	}
}
