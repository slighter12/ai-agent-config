//go:build aix || android || darwin || dragonfly || freebsd || ios || linux || netbsd || openbsd || solaris

package skillcatalog

import (
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"testing"
)

func TestLoadRetiredSkillsRejectsManifestFIFOWithoutOpeningIt(t *testing.T) {
	path := filepath.Join(t.TempDir(), retiredSkillsFile)
	if err := syscall.Mkfifo(path, 0o600); err != nil {
		t.Fatal(err)
	}

	if _, err := LoadRetiredSkills(path); err == nil || !strings.Contains(err.Error(), "not a regular file") {
		t.Fatalf("expected manifest FIFO rejection, got %v", err)
	}
}

func TestActiveRoutingPathsRejectsRoutingFIFO(t *testing.T) {
	repo := t.TempDir()
	skills := filepath.Join(repo, "skills")
	if err := os.Mkdir(skills, 0o755); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(repo, "README.md")
	if err := syscall.Mkfifo(path, 0o600); err != nil {
		t.Fatal(err)
	}

	if _, err := activeRoutingPaths(repo, skills, []string{"retired-skill"}); err == nil || !strings.Contains(err.Error(), "not a regular file") {
		t.Fatalf("expected routing FIFO rejection, got %v", err)
	}
}
