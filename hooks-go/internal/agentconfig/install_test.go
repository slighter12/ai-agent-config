package agentconfig

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestCreateSymlinkIsIdempotentAndSkipsUserFiles(t *testing.T) {
	dir := t.TempDir()
	source := filepath.Join(dir, "source")
	target := filepath.Join(dir, "target")
	userFile := filepath.Join(dir, "user-file")
	var out bytes.Buffer
	config := Config{Out: &out}
	if err := os.WriteFile(source, []byte("source"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := config.createSymlink(source, target); err != nil {
		t.Fatal(err)
	}
	if err := config.createSymlink(source, target); err != nil {
		t.Fatal(err)
	}
	if link, err := os.Readlink(target); err != nil || link != source {
		t.Fatalf("unexpected symlink %q, err %v", link, err)
	}
	if err := os.WriteFile(userFile, []byte("mine"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := config.createSymlink(source, userFile); err != nil {
		t.Fatal(err)
	}
	if got := readFile(t, userFile); got != "mine" {
		t.Fatalf("user-owned file was overwritten: %q", got)
	}
}

func TestGenerateCodexRoleFilesRendersTemplatesAndProtectsUserFiles(t *testing.T) {
	repo := t.TempDir()
	sourceDir := filepath.Join(repo, "config", "codex-agents")
	targetDir := filepath.Join(repo, "home", ".codex", "agents")
	if err := os.MkdirAll(sourceDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(targetDir, 0o755); err != nil {
		t.Fatal(err)
	}
	template := `name = "demo"
[[skills.config]]
path = "{{REPO_ROOT}}/skills/demo"
[[skills.config]]
path = "{{AGENTS_HOME}}/skills/demo"
`
	if err := os.WriteFile(filepath.Join(sourceDir, "demo.toml"), []byte(template), 0o644); err != nil {
		t.Fatal(err)
	}
	userTarget := filepath.Join(targetDir, "user.toml")
	if err := os.WriteFile(filepath.Join(sourceDir, "user.toml"), []byte(`name = "user"`), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(userTarget, []byte("user-owned"), 0o644); err != nil {
		t.Fatal(err)
	}
	config := Config{
		RepoRoot:   repo,
		Home:       filepath.Join(repo, "home"),
		CodexHome:  filepath.Join(repo, "home", ".codex"),
		AgentsHome: filepath.Join(repo, "home", ".agents"),
		Out:        &bytes.Buffer{},
		Now:        func() time.Time { return time.Unix(0, 0) },
	}
	if err := config.generateCodexRoleFiles(targetDir); err != nil {
		t.Fatal(err)
	}
	got := readFile(t, filepath.Join(targetDir, "demo.toml"))
	mustContain(t, got, managedHeader)
	mustContain(t, got, filepath.ToSlash(filepath.Join(repo, "skills", "demo")))
	mustContain(t, got, filepath.ToSlash(filepath.Join(repo, "home", ".agents", "skills", "demo")))
	if got := readFile(t, userTarget); got != "user-owned" {
		t.Fatalf("user-owned role file was overwritten: %q", got)
	}
}

func TestGenerateCodexRoleFilesReplacesLegacyTemplateSymlink(t *testing.T) {
	repo := t.TempDir()
	sourceDir := filepath.Join(repo, "config", "codex-agents")
	targetDir := filepath.Join(repo, "home", ".codex", "agents")
	if err := os.MkdirAll(sourceDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(targetDir, 0o755); err != nil {
		t.Fatal(err)
	}
	source := filepath.Join(sourceDir, "demo.toml")
	target := filepath.Join(targetDir, "demo.toml")
	if err := os.WriteFile(source, []byte(`path = "{{HOME}}/demo"`), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(source, target); err != nil {
		t.Fatal(err)
	}
	config := Config{
		RepoRoot:   repo,
		Home:       filepath.Join(repo, "home"),
		CodexHome:  filepath.Join(repo, "home", ".codex"),
		AgentsHome: filepath.Join(repo, "home", ".agents"),
		Out:        &bytes.Buffer{},
	}
	if err := config.generateCodexRoleFiles(targetDir); err != nil {
		t.Fatal(err)
	}
	info, err := os.Lstat(target)
	if err != nil {
		t.Fatal(err)
	}
	if info.Mode()&os.ModeSymlink != 0 {
		t.Fatal("expected generated regular file, got symlink")
	}
	mustContain(t, readFile(t, target), filepath.ToSlash(filepath.Join(repo, "home", "demo")))
}

func TestAppendZshrcSourceIsIdempotent(t *testing.T) {
	dir := t.TempDir()
	zshrc := filepath.Join(dir, ".zshrc")
	var out bytes.Buffer
	config := Config{Home: dir, Out: &out}
	if err := config.appendZshrcSource(zshrc); err != nil {
		t.Fatal(err)
	}
	if err := config.appendZshrcSource(zshrc); err != nil {
		t.Fatal(err)
	}
	got := readFile(t, zshrc)
	if count := strings.Count(got, "# Codex repo-managed shell helpers"); count != 1 {
		t.Fatalf("expected one managed block, got %d:\n%s", count, got)
	}
}

func TestPrepareCodexAgentDirSkipsLegacyDirectorySymlink(t *testing.T) {
	repo := t.TempDir()
	target := filepath.Join(repo, "home", ".codex", "agents")
	if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(repo, "agents"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(filepath.Join(repo, "agents"), target); err != nil {
		t.Fatal(err)
	}
	config := Config{RepoRoot: repo, Out: &bytes.Buffer{}}
	ok, err := config.prepareCodexAgentDir(target)
	if err != nil {
		t.Fatal(err)
	}
	if ok {
		t.Fatal("expected legacy symlink to be skipped")
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
