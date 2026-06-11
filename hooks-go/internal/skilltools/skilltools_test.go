package skilltools

import (
	"archive/zip"
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestValidateSkillAcceptsPortableFrontmatter(t *testing.T) {
	dir := t.TempDir()
	skillDir := filepath.Join(dir, "demo-skill")
	if err := os.MkdirAll(skillDir, 0o755); err != nil {
		t.Fatal(err)
	}
	writeFile(t, filepath.Join(skillDir, "SKILL.md"), `---
name: demo-skill
description: Perform demo validation. Use when checking the Go validator. Avoid when another skill is more specific.
license: MIT
compatibility: [codex, claude, gemini]
metadata:
  version: "0.1.0"
---

# Demo Skill
`)

	valid, message, err := ValidateSkill(skillDir)
	if err != nil {
		t.Fatal(err)
	}
	if !valid || message != "Skill is valid" {
		t.Fatalf("expected valid skill, got valid=%v message=%q", valid, message)
	}
}

func TestValidateSkillRejectsPlaceholderDescription(t *testing.T) {
	dir := t.TempDir()
	skillDir := filepath.Join(dir, "demo-skill")
	if err := os.MkdirAll(skillDir, 0o755); err != nil {
		t.Fatal(err)
	}
	writeFile(t, filepath.Join(skillDir, "SKILL.md"), `---
name: demo-skill
description: Perform [capability]. Use when testing. Avoid when done.
---

# Demo Skill
`)

	valid, message, err := ValidateSkill(skillDir)
	if err != nil {
		t.Fatal(err)
	}
	if valid {
		t.Fatal("expected invalid skill")
	}
	if !strings.Contains(message, "placeholder") {
		t.Fatalf("expected placeholder error, got %q", message)
	}
}

func TestInitSkillCreatesTemplateFiles(t *testing.T) {
	dir := t.TempDir()
	var out bytes.Buffer
	if err := InitSkill(&out, "demo-skill", dir); err != nil {
		t.Fatal(err)
	}
	skillDir := filepath.Join(dir, "demo-skill")
	mustContain(t, readFile(t, filepath.Join(skillDir, "SKILL.md")), "name: demo-skill")
	mustContain(t, readFile(t, filepath.Join(skillDir, "SKILL.md")), "# Demo Skill")
	mustContain(t, readFile(t, filepath.Join(skillDir, "references", "INDEX.md")), "# References")
}

func TestPackageSkillCreatesSkillArchive(t *testing.T) {
	dir := t.TempDir()
	skillDir := filepath.Join(dir, "demo-skill")
	if err := os.MkdirAll(filepath.Join(skillDir, "references"), 0o755); err != nil {
		t.Fatal(err)
	}
	writeFile(t, filepath.Join(skillDir, "SKILL.md"), `---
name: demo-skill
description: Perform demo packaging. Use when packaging a test skill. Avoid when validation should fail.
metadata:
  version: "0.1.0"
---

# Demo Skill
`)
	writeFile(t, filepath.Join(skillDir, "references", "INDEX.md"), "# References\n")
	outDir := filepath.Join(dir, "dist")

	var out bytes.Buffer
	archivePath, err := PackageSkill(&out, skillDir, outDir)
	if err != nil {
		t.Fatal(err)
	}
	reader, err := zip.OpenReader(archivePath)
	if err != nil {
		t.Fatal(err)
	}
	defer reader.Close()
	names := map[string]bool{}
	for _, file := range reader.File {
		names[file.Name] = true
	}
	for _, name := range []string{"demo-skill/SKILL.md", "demo-skill/references/INDEX.md"} {
		if !names[name] {
			t.Fatalf("expected archive entry %s, got %#v", name, names)
		}
	}
}

func writeFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func readFile(t *testing.T, path string) string {
	t.Helper()
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return string(raw)
}

func mustContain(t *testing.T, text, want string) {
	t.Helper()
	if !strings.Contains(text, want) {
		t.Fatalf("expected %q in:\n%s", want, text)
	}
}
