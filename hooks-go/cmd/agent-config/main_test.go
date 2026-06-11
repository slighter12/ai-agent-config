package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestRunInitSkillAcceptsPathAfterSkillName(t *testing.T) {
	dir := t.TempDir()

	if err := runInitSkill([]string{"demo-skill", "--path", dir}); err != nil {
		t.Fatal(err)
	}

	if _, err := os.Stat(filepath.Join(dir, "demo-skill", "SKILL.md")); err != nil {
		t.Fatal(err)
	}
}
