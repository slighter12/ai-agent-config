package agentconfig

import (
	"bytes"
	"errors"
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

func TestCreateSkillSymlinksSkipsEntriesWithoutSkillManifest(t *testing.T) {
	dir := t.TempDir()
	sourceDir := filepath.Join(dir, "skills")
	targetDir := filepath.Join(dir, "home", "skills")
	validSkill := filepath.Join(sourceDir, "valid-skill")
	emptyDir := filepath.Join(sourceDir, "empty-dir")
	if err := os.MkdirAll(validSkill, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(emptyDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(validSkill, "SKILL.md"), []byte("name: valid-skill\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(sourceDir, "README.md"), []byte("not a skill"), 0o644); err != nil {
		t.Fatal(err)
	}
	var out bytes.Buffer
	config := Config{Out: &out}
	if err := config.createSkillSymlinks(sourceDir, targetDir); err != nil {
		t.Fatal(err)
	}
	if link, err := os.Readlink(filepath.Join(targetDir, "valid-skill")); err != nil || link != validSkill {
		t.Fatalf("unexpected valid skill symlink %q, err %v", link, err)
	}
	if _, err := os.Lstat(filepath.Join(targetDir, "empty-dir")); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("expected empty dir to be skipped, err %v", err)
	}
	if _, err := os.Lstat(filepath.Join(targetDir, "README.md")); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("expected file entry to be skipped, err %v", err)
	}
	mustContain(t, out.String(), "Skipping non-skill entry")
}

func TestCodexMarketplaceConfiguredDetectsRepoMarketplace(t *testing.T) {
	dir := t.TempDir()
	codexHome := filepath.Join(dir, ".codex")
	if err := os.MkdirAll(codexHome, 0o755); err != nil {
		t.Fatal(err)
	}
	config := Config{CodexHome: codexHome, RepoRoot: filepath.Join(dir, "repo")}
	if config.codexMarketplaceConfigured() {
		t.Fatal("expected missing config to be unconfigured")
	}
	if err := os.WriteFile(filepath.Join(codexHome, "config.toml"), []byte(`
[marketplaces.other]
source = "/tmp/other"

[marketplaces.ai-agent-config]
source_type = "local"
source = "`+filepath.ToSlash(filepath.Join(dir, "other-repo"))+`"
`), 0o644); err != nil {
		t.Fatal(err)
	}
	if config.codexMarketplaceConfigured() {
		t.Fatal("expected marketplace with different source to be unconfigured")
	}
	if err := os.WriteFile(filepath.Join(codexHome, "config.toml"), []byte(`
[marketplaces.ai-agent-config]
source_type = "local"
source = "`+filepath.ToSlash(config.RepoRoot)+`"
`), 0o644); err != nil {
		t.Fatal(err)
	}
	if !config.codexMarketplaceConfigured() {
		t.Fatal("expected ai-agent-config marketplace to be configured")
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
`
	if err := os.WriteFile(filepath.Join(sourceDir, "demo.toml"), []byte(template), 0o644); err != nil {
		t.Fatal(err)
	}
	writeRoleManifest(t, sourceDir, `{
  "skill_groups": {
    "demo": [
      "{{REPO_ROOT}}/skills/demo",
      "{{AGENTS_HOME}}/skills/demo"
    ]
  },
  "roles": {
    "demo": {
      "disabled_skill_groups": ["demo"]
    }
  }
}
`)
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
	mustContain(t, got, "enabled = false")
	if got := readFile(t, userTarget); got != "user-owned" {
		t.Fatalf("user-owned role file was overwritten: %q", got)
	}
}

func TestGitCommitTemplateDeclaresWorkspaceGitProfile(t *testing.T) {
	repoRoot := filepath.Clean(filepath.Join("..", "..", ".."))
	template := readFile(t, filepath.Join(repoRoot, "config", "codex-agents", "git-commit.toml"))
	mustContain(t, template, `default_permissions = "workspace-git"`)
	mustContain(t, template, `[permissions.workspace-git]`)
	mustContain(t, template, `[permissions.workspace-git.filesystem.":workspace_roots"]`)
	mustContain(t, template, `".git" = "write"`)
	mustContain(t, template, `[permissions.workspace-git.network.domains]`)
}

func TestGenerateCodexRoleFilesSkipsAlreadyCurrentManagedFiles(t *testing.T) {
	repo := t.TempDir()
	sourceDir := filepath.Join(repo, "config", "codex-agents")
	targetDir := filepath.Join(repo, "home", ".codex", "agents")
	if err := os.MkdirAll(sourceDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(targetDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(sourceDir, "demo.toml"), []byte(`name = "demo"`), 0o644); err != nil {
		t.Fatal(err)
	}
	writeRoleManifest(t, sourceDir, `{"skill_groups":{"demo":["{{REPO_ROOT}}/skills/demo"]},"roles":{"demo":{"disabled_skill_groups":["demo"]}}}`)
	var out bytes.Buffer
	config := Config{
		RepoRoot:   repo,
		Home:       filepath.Join(repo, "home"),
		CodexHome:  filepath.Join(repo, "home", ".codex"),
		AgentsHome: filepath.Join(repo, "home", ".agents"),
		Out:        &out,
	}
	if err := config.generateCodexRoleFiles(targetDir); err != nil {
		t.Fatal(err)
	}
	if err := config.generateCodexRoleFiles(targetDir); err != nil {
		t.Fatal(err)
	}
	mustContain(t, out.String(), "Role file already current")
}

func TestGenerateCodexRoleFilesRejectsUnknownSkillToggleGroup(t *testing.T) {
	repo := t.TempDir()
	sourceDir := filepath.Join(repo, "config", "codex-agents")
	targetDir := filepath.Join(repo, "home", ".codex", "agents")
	if err := os.MkdirAll(sourceDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(targetDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(sourceDir, "demo.toml"), []byte(`name = "demo"`), 0o644); err != nil {
		t.Fatal(err)
	}
	writeRoleManifest(t, sourceDir, `{"skill_groups":{"demo":["{{REPO_ROOT}}/skills/demo"]},"roles":{"demo":{"disabled_skill_groups":["missing"]}}}`)
	config := Config{RepoRoot: repo, Out: &bytes.Buffer{}}
	err := config.generateCodexRoleFiles(targetDir)
	if err == nil || !strings.Contains(err.Error(), `unknown skill group "missing"`) {
		t.Fatalf("expected unknown group error, got %v", err)
	}
}

func TestGenerateCodexRoleFilesRejectsUnknownMCPServer(t *testing.T) {
	repo := t.TempDir()
	sourceDir := filepath.Join(repo, "config", "codex-agents")
	targetDir := filepath.Join(repo, "home", ".codex", "agents")
	if err := os.MkdirAll(sourceDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(targetDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(sourceDir, "demo.toml"), []byte(`name = "demo"`), 0o644); err != nil {
		t.Fatal(err)
	}
	writeRoleManifest(t, sourceDir, `{
  "mcp_servers": [{"name": "context7", "command": "bunx"}],
  "skill_groups": {"demo": ["{{REPO_ROOT}}/skills/demo"]},
  "roles": {"demo": {"allowed_mcp_servers": ["missing"]}}
}`)
	config := Config{RepoRoot: repo, Out: &bytes.Buffer{}}
	err := config.generateCodexRoleFiles(targetDir)
	if err == nil || !strings.Contains(err.Error(), `unknown mcp server "missing"`) {
		t.Fatalf("expected unknown mcp server error, got %v", err)
	}
}

func TestGenerateCodexRoleFilesRendersMCPServersWithAvailability(t *testing.T) {
	repo := t.TempDir()
	sourceDir := filepath.Join(repo, "config", "codex-agents")
	targetDir := filepath.Join(repo, "home", ".codex", "agents")
	if err := os.MkdirAll(sourceDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(targetDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(repo, "home", ".codex", "config.toml"), []byte("[mcp_servers.sequential-thinking]\n\n[mcp_servers.clickup]\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(sourceDir, "demo.toml"), []byte(`name = "demo"`), 0o644); err != nil {
		t.Fatal(err)
	}
	writeRoleManifest(t, sourceDir, `{
  "mcp_servers": [
    {
      "name": "context7",
      "command": "bunx",
      "args": ["-y", "@upstash/context7-mcp"],
      "env": {"DEFAULT_MINIMUM_TOKENS": "40000"}
    },
    {
      "name": "sequential-thinking",
      "command": "bunx",
      "args": ["-y", "@modelcontextprotocol/server-sequential-thinking"]
    },
    {
      "name": "clickup",
      "command": "bunx",
      "args": ["-y", "mcp-remote", "https://mcp.clickup.com/mcp"]
    }
  ],
  "skill_groups": {"demo": ["{{REPO_ROOT}}/skills/demo"]},
  "roles": {
    "demo": {
      "allowed_mcp_servers": ["sequential-thinking", "clickup"],
      "disabled_skill_groups": ["demo"]
    }
  }
}`)
	config := Config{
		RepoRoot:   repo,
		Home:       filepath.Join(repo, "home"),
		CodexHome:  filepath.Join(repo, "home", ".codex"),
		AgentsHome: filepath.Join(repo, "home", ".agents"),
		Out:        &bytes.Buffer{},
		LookPath: func(command string) (string, error) {
			if command == "bunx" {
				return filepath.Join(repo, "bin", "bunx"), nil
			}
			return "", os.ErrNotExist
		},
	}
	if err := config.generateCodexRoleFiles(targetDir); err != nil {
		t.Fatal(err)
	}
	got := readFile(t, filepath.Join(targetDir, "demo.toml"))
	mustContain(t, got, "[mcp_servers.context7]")
	mustContain(t, got, `args = ["-y", "@upstash/context7-mcp"]`)
	mustContain(t, got, "enabled = false\n\n[mcp_servers.context7.env]")
	mustContain(t, got, `DEFAULT_MINIMUM_TOKENS = "40000"`)
	mustContain(t, got, "[mcp_servers.sequential-thinking]")
	mustContain(t, got, "enabled = true\n\n[mcp_servers.clickup]")
	mustContain(t, got, `path = "`+filepath.ToSlash(filepath.Join(repo, "skills", "demo"))+`"`)
}

func TestGenerateCodexRoleFilesDisablesAllowedMCPWhenCommandMissing(t *testing.T) {
	repo := t.TempDir()
	sourceDir := filepath.Join(repo, "config", "codex-agents")
	targetDir := filepath.Join(repo, "home", ".codex", "agents")
	if err := os.MkdirAll(sourceDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(targetDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(repo, "home", ".codex", "config.toml"), []byte("[mcp_servers.sequential-thinking]\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(sourceDir, "demo.toml"), []byte(`name = "demo"`), 0o644); err != nil {
		t.Fatal(err)
	}
	writeRoleManifest(t, sourceDir, `{
  "mcp_servers": [{"name": "sequential-thinking", "command": "bunx"}],
  "skill_groups": {"demo": ["{{REPO_ROOT}}/skills/demo"]},
  "roles": {"demo": {"allowed_mcp_servers": ["sequential-thinking"]}}
}`)
	config := Config{
		RepoRoot: repo,
		Out:      &bytes.Buffer{},
		LookPath: func(string) (string, error) {
			return "", os.ErrNotExist
		},
	}
	if err := config.generateCodexRoleFiles(targetDir); err != nil {
		t.Fatal(err)
	}
	got := readFile(t, filepath.Join(targetDir, "demo.toml"))
	mustContain(t, got, "[mcp_servers.sequential-thinking]")
	mustContain(t, got, "enabled = false")
	mustNotContain(t, got, "enabled = true")
}

func TestGenerateCodexRoleFilesDisablesAllowedMCPWhenNotConfigured(t *testing.T) {
	repo := t.TempDir()
	sourceDir := filepath.Join(repo, "config", "codex-agents")
	targetDir := filepath.Join(repo, "home", ".codex", "agents")
	if err := os.MkdirAll(sourceDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(targetDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(sourceDir, "demo.toml"), []byte(`name = "demo"`), 0o644); err != nil {
		t.Fatal(err)
	}
	writeRoleManifest(t, sourceDir, `{
  "mcp_servers": [{"name": "sequential-thinking", "command": "bunx"}],
  "skill_groups": {"demo": ["{{REPO_ROOT}}/skills/demo"]},
  "roles": {"demo": {"allowed_mcp_servers": ["sequential-thinking"]}}
}`)
	config := Config{
		RepoRoot:  repo,
		CodexHome: filepath.Join(repo, "home", ".codex"),
		Out:       &bytes.Buffer{},
		LookPath: func(command string) (string, error) {
			if command == "bunx" {
				return filepath.Join(repo, "bin", "bunx"), nil
			}
			return "", os.ErrNotExist
		},
	}
	if err := config.generateCodexRoleFiles(targetDir); err != nil {
		t.Fatal(err)
	}
	got := readFile(t, filepath.Join(targetDir, "demo.toml"))
	mustContain(t, got, "[mcp_servers.sequential-thinking]")
	mustContain(t, got, "enabled = false")
	mustNotContain(t, got, "enabled = true")
}

func TestRenderRoleTemplateRejectsNameMismatch(t *testing.T) {
	dir := t.TempDir()
	source := filepath.Join(dir, "demo.toml")
	if err := os.WriteFile(source, []byte(`name = "other"`), 0o644); err != nil {
		t.Fatal(err)
	}
	config := Config{RepoRoot: dir, Out: &bytes.Buffer{}}
	_, err := config.renderRoleTemplate(source, "demo", codexRoleManifest{})
	if err == nil || !strings.Contains(err.Error(), `declares role name "other", expected "demo"`) {
		t.Fatalf("expected role name mismatch error, got %v", err)
	}
}

func TestRenderRoleTemplateRejectsMCPBlocks(t *testing.T) {
	dir := t.TempDir()
	source := filepath.Join(dir, "demo.toml")
	if err := os.WriteFile(source, []byte("name = \"demo\"\n[mcp_servers.context7]\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	config := Config{RepoRoot: dir, Out: &bytes.Buffer{}}
	_, err := config.renderRoleTemplate(source, "demo", codexRoleManifest{})
	if err == nil || !strings.Contains(err.Error(), "role-manifest.json") {
		t.Fatalf("expected mcp block rejection, got %v", err)
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
	writeRoleManifest(t, sourceDir, `{"skill_groups":{"demo":["{{REPO_ROOT}}/skills/demo"]},"roles":{"demo":{"disabled_skill_groups":["demo"]}}}`)
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

func mustNotContain(t *testing.T, text, unwanted string) {
	t.Helper()
	if strings.Contains(text, unwanted) {
		t.Fatalf("did not expect %q in:\n%s", unwanted, text)
	}
}

func writeRoleManifest(t *testing.T, sourceDir, text string) {
	t.Helper()
	if err := os.WriteFile(filepath.Join(sourceDir, "role-manifest.json"), []byte(text), 0o644); err != nil {
		t.Fatal(err)
	}
}
