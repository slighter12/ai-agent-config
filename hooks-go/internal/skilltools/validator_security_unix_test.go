//go:build android || darwin || dragonfly || freebsd || ios || linux || netbsd || openbsd

package skilltools

import (
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"testing"
)

func TestValidateSkillRejectsFIFOWithoutOpeningIt(t *testing.T) {
	skillDir := filepath.Join(t.TempDir(), "demo")
	if err := os.Mkdir(skillDir, 0o755); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(skillDir, "SKILL.md")
	if err := syscall.Mkfifo(path, 0o600); err != nil {
		t.Fatal(err)
	}
	if _, _, err := ValidateSkill(skillDir); err == nil || !strings.Contains(err.Error(), "not a regular file") {
		t.Fatalf("expected FIFO rejection, got %v", err)
	}
}
