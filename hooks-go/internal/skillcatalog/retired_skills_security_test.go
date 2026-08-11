package skillcatalog

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoadRetiredSkillsRejectsManifestSymlink(t *testing.T) {
	directory := t.TempDir()
	target := filepath.Join(directory, "target.json")
	manifest := filepath.Join(directory, retiredSkillsFile)
	if err := os.WriteFile(target, []byte(`{"retired_skills": []}`), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(target, manifest); err != nil {
		t.Skipf("symlinks unavailable: %v", err)
	}

	if _, err := LoadRetiredSkills(manifest); err == nil || !strings.Contains(err.Error(), "symlink") {
		t.Fatalf("expected manifest symlink rejection, got %v", err)
	}
}

func TestActiveRoutingPathsRejectsLeafSymlink(t *testing.T) {
	repo := t.TempDir()
	skills := filepath.Join(repo, "skills")
	if err := os.Mkdir(skills, 0o755); err != nil {
		t.Fatal(err)
	}
	target := filepath.Join(repo, "target.md")
	readme := filepath.Join(repo, "README.md")
	if err := os.WriteFile(target, []byte("outside"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(target, readme); err != nil {
		t.Skipf("symlinks unavailable: %v", err)
	}

	if _, err := activeRoutingPaths(repo, skills, []string{"retired-skill"}); err == nil || !strings.Contains(err.Error(), "symlink") {
		t.Fatalf("expected routing leaf symlink rejection, got %v", err)
	}
}

func TestActiveRoutingPathsRejectsIntermediateSymlink(t *testing.T) {
	repo := t.TempDir()
	skills := filepath.Join(repo, "skills")
	outside := filepath.Join(t.TempDir(), "routing")
	if err := os.Mkdir(skills, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Mkdir(outside, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(outside, "route.md"), []byte("outside"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(outside, filepath.Join(repo, "config")); err != nil {
		t.Skipf("symlinks unavailable: %v", err)
	}

	if _, err := activeRoutingPaths(repo, skills, []string{"retired-skill"}); err == nil || !strings.Contains(err.Error(), "symlink") {
		t.Fatalf("expected routing intermediate symlink rejection, got %v", err)
	}
}

func TestReadRetiredSkillsFileDetectsReplacementAfterOpen(t *testing.T) {
	directory := t.TempDir()
	path := filepath.Join(directory, retiredSkillsFile)
	replacement := filepath.Join(directory, "replacement.json")
	raw := []byte(`{"retired_skills": []}`)
	if err := os.WriteFile(path, raw, 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(replacement, raw, 0o644); err != nil {
		t.Fatal(err)
	}
	var callbackErr error
	_, err := readRetiredSkillsFile(path, func() {
		callbackErr = os.Rename(replacement, path)
	})
	if callbackErr != nil {
		t.Fatal(callbackErr)
	}
	if err == nil || !strings.Contains(err.Error(), "replaced while reading") {
		t.Fatalf("expected replacement rejection, got %v", err)
	}
}

func TestReadRetiredSkillsFileDetectsGrowthAfterOpen(t *testing.T) {
	directory := t.TempDir()
	path := filepath.Join(directory, retiredSkillsFile)
	raw := []byte(`{"retired_skills": []}`)
	if err := os.WriteFile(path, raw, 0o644); err != nil {
		t.Fatal(err)
	}
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_APPEND, 0)
	if err != nil {
		t.Fatal(err)
	}
	file.Close()

	var callbackErr error
	_, err = readRetiredSkillsFile(path, func() {
		growth, openErr := os.OpenFile(path, os.O_WRONLY|os.O_APPEND, 0)
		if openErr != nil {
			callbackErr = openErr
			return
		}
		_, callbackErr = growth.WriteString(" ")
		if closeErr := growth.Close(); callbackErr == nil {
			callbackErr = closeErr
		}
	})
	if callbackErr != nil {
		t.Fatal(callbackErr)
	}
	if err == nil || !strings.Contains(err.Error(), "grew or changed while reading") {
		t.Fatalf("expected growth rejection, got %v", err)
	}
}

func TestScanRoutingFileDetectsReplacementAfterOpen(t *testing.T) {
	directory := t.TempDir()
	path := filepath.Join(directory, "route.md")
	replacement := filepath.Join(directory, "replacement.md")
	if err := os.WriteFile(path, []byte("safe"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(replacement, []byte("safe"), 0o644); err != nil {
		t.Fatal(err)
	}
	var callbackErr error
	failures := []string{}
	budget := routingScanBudget{remaining: maxRoutingSurfaceAggregateBytes}
	err := scanRoutingFileWithHook(path, nil, &failures, &budget, func() {
		callbackErr = os.Rename(replacement, path)
	})
	if callbackErr != nil {
		t.Fatal(callbackErr)
	}
	if err == nil || !strings.Contains(err.Error(), "replaced while reading") {
		t.Fatalf("expected routing replacement rejection, got %v", err)
	}
}

func TestScanRoutingFileDetectsGrowthAfterOpen(t *testing.T) {
	directory := t.TempDir()
	path := filepath.Join(directory, "route.md")
	if err := os.WriteFile(path, []byte("safe"), 0o644); err != nil {
		t.Fatal(err)
	}
	var callbackErr error
	failures := []string{}
	budget := routingScanBudget{remaining: maxRoutingSurfaceAggregateBytes}
	err := scanRoutingFileWithHook(path, nil, &failures, &budget, func() {
		growth, openErr := os.OpenFile(path, os.O_WRONLY|os.O_APPEND, 0)
		if openErr != nil {
			callbackErr = openErr
			return
		}
		_, callbackErr = growth.WriteString(" ")
		if closeErr := growth.Close(); callbackErr == nil {
			callbackErr = closeErr
		}
	})
	if callbackErr != nil {
		t.Fatal(callbackErr)
	}
	if err == nil || !strings.Contains(err.Error(), "grew or changed while reading") {
		t.Fatalf("expected routing growth rejection, got %v", err)
	}
}

func TestLoadRetiredSkillsFromRootRejectsSymlinkedConfigDirectory(t *testing.T) {
	repo := t.TempDir()
	outside := filepath.Join(t.TempDir(), "config")
	if err := os.Mkdir(outside, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(outside, retiredSkillsFile), []byte(`{"retired_skills": []}`), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(outside, filepath.Join(repo, "config")); err != nil {
		t.Skipf("symlinks unavailable: %v", err)
	}

	if _, err := LoadRetiredSkillsFromRoot(repo); err == nil || !strings.Contains(err.Error(), "symlink") {
		t.Fatalf("expected symlinked config directory rejection, got %v", err)
	}
}

func TestLoadRetiredSkillsFromRootAcceptsSymlinkedRepositoryRoot(t *testing.T) {
	repo := t.TempDir()
	config := filepath.Join(repo, "config")
	if err := os.Mkdir(config, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(config, retiredSkillsFile), []byte(`{"retired_skills": ["retired-skill"]}`), 0o644); err != nil {
		t.Fatal(err)
	}
	link := filepath.Join(t.TempDir(), "repo-link")
	if err := os.Symlink(repo, link); err != nil {
		t.Skipf("symlinks unavailable: %v", err)
	}

	names, err := LoadRetiredSkillsFromRoot(link)
	if err != nil {
		t.Fatalf("expected symlinked repository root to load, got %v", err)
	}
	if len(names) != 1 || names[0] != "retired-skill" {
		t.Fatalf("unexpected retired names: %v", names)
	}
}
