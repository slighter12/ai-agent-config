package main

import (
	"os"
	"path/filepath"
	"strings"
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

func TestRepositorySkillsDirectoryAcceptsRelativePath(t *testing.T) {
	dir := t.TempDir()
	repoRoot := filepath.Join(dir, "repo")
	skillsDir := filepath.Join(repoRoot, "skills")
	if err := os.MkdirAll(skillsDir, 0o755); err != nil {
		t.Fatal(err)
	}
	oldWorkingDirectory, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := os.Chdir(oldWorkingDirectory); err != nil {
			t.Errorf("restore working directory: %v", err)
		}
	})

	got, err := repositorySkillsDirectory("repo", filepath.Join("repo", "skills"))
	if err != nil {
		t.Fatal(err)
	}
	want, err := filepath.EvalSymlinks(skillsDir)
	if err != nil {
		t.Fatal(err)
	}
	if got != want {
		t.Fatalf("repositorySkillsDirectory() = %q, want %q", got, want)
	}
}

func TestRepositorySkillsDirectoryAcceptsSymlinkAlias(t *testing.T) {
	dir := t.TempDir()
	repoRoot := filepath.Join(dir, "repo")
	skillsDir := filepath.Join(repoRoot, "skills")
	if err := os.MkdirAll(skillsDir, 0o755); err != nil {
		t.Fatal(err)
	}
	alias := filepath.Join(dir, "skills-alias")
	if err := os.Symlink(skillsDir, alias); err != nil {
		t.Fatal(err)
	}

	got, err := repositorySkillsDirectory(repoRoot, alias)
	if err != nil {
		t.Fatal(err)
	}
	want, err := filepath.EvalSymlinks(skillsDir)
	if err != nil {
		t.Fatal(err)
	}
	if got != want {
		t.Fatalf("repositorySkillsDirectory() = %q, want %q", got, want)
	}
}

func TestRepositorySkillsDirectoryRejectsDifferentDirectories(t *testing.T) {
	dir := t.TempDir()
	repoRoot := filepath.Join(dir, "repo")
	if err := os.MkdirAll(filepath.Join(repoRoot, "skills"), 0o755); err != nil {
		t.Fatal(err)
	}
	otherRepoSkills := filepath.Join(dir, "other-repo", "skills")
	siblingSkills := filepath.Join(repoRoot, "sibling-skills")
	for _, path := range []string{otherRepoSkills, siblingSkills} {
		if err := os.MkdirAll(path, 0o755); err != nil {
			t.Fatal(err)
		}
	}

	for _, skillsDir := range []string{otherRepoSkills, siblingSkills} {
		t.Run(filepath.Base(filepath.Dir(skillsDir))+"_"+filepath.Base(skillsDir), func(t *testing.T) {
			_, err := repositorySkillsDirectory(repoRoot, skillsDir)
			if err == nil {
				t.Fatalf("expected %q to be rejected", skillsDir)
			}
			if !strings.Contains(err.Error(), "pass <repo-root>/skills") {
				t.Fatalf("error lacks corrective guidance: %v", err)
			}
		})
	}
}

func TestRepositorySkillsDirectoryReportsExpectedPathGuidance(t *testing.T) {
	repoRoot := t.TempDir()
	_, err := repositorySkillsDirectory(repoRoot, filepath.Join(repoRoot, "skills"))
	if err == nil {
		t.Fatal("expected missing repository skills directory to be rejected")
	}
	if !strings.Contains(err.Error(), "pass <repo-root>/skills") {
		t.Fatalf("error lacks corrective guidance: %v", err)
	}
}
