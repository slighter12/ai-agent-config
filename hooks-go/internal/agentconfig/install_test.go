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
	if err := config.createSkillSymlinks(sourceDir, targetDir, []string{"retired-alpha"}); err != nil {
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
	if strings.Contains(out.String(), "README.md") {
		t.Fatalf("ordinary files should be skipped silently:\n%s", out.String())
	}
}

func TestCreateSkillSymlinksRemovesOnlyRetiredRepoSymlinks(t *testing.T) {
	dir := t.TempDir()
	sourceDir := filepath.Join(dir, "repo", "skills")
	targetDir := filepath.Join(dir, "home", "skills")
	validSkill := filepath.Join(sourceDir, "valid-skill")
	if err := os.MkdirAll(validSkill, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(validSkill, "SKILL.md"), []byte("name: valid-skill\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(targetDir, 0o755); err != nil {
		t.Fatal(err)
	}
	retiredNames := []string{"retired-alpha", "retired-beta", "retired-gamma"}
	ownedRetired := filepath.Join(targetDir, retiredNames[0])
	if err := os.Symlink(filepath.Join(sourceDir, retiredNames[0]), ownedRetired); err != nil {
		t.Fatal(err)
	}
	foreignRetired := filepath.Join(targetDir, retiredNames[1])
	foreignSource := filepath.Join(dir, "foreign", retiredNames[1])
	if err := os.Symlink(foreignSource, foreignRetired); err != nil {
		t.Fatal(err)
	}
	userRetired := filepath.Join(targetDir, retiredNames[2])
	if err := os.WriteFile(userRetired, []byte("user-owned"), 0o644); err != nil {
		t.Fatal(err)
	}

	var out bytes.Buffer
	config := Config{Out: &out}
	if err := config.createSkillSymlinks(sourceDir, targetDir, retiredNames); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Lstat(ownedRetired); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("expected repo-owned retired symlink removed, err=%v", err)
	}
	if link, err := os.Readlink(foreignRetired); err != nil || link != foreignSource {
		t.Fatalf("foreign retired symlink changed: link=%q err=%v", link, err)
	}
	if got := readFile(t, userRetired); got != "user-owned" {
		t.Fatalf("user-owned retired target changed: %q", got)
	}
	if _, err := os.Readlink(filepath.Join(targetDir, "valid-skill")); err != nil {
		t.Fatalf("current skill was not linked: %v", err)
	}
	mustContain(t, out.String(), "Removed retired repo skill symlink")
	mustContain(t, out.String(), "leaving unchanged")
}

func TestCreateSkillSymlinksPreservesRepromotedRetiredSkill(t *testing.T) {
	dir := t.TempDir()
	sourceDir := filepath.Join(dir, "repo", "skills")
	targetDir := filepath.Join(dir, "home", "skills")
	retiredName := "retired-alpha"
	repromoted := filepath.Join(sourceDir, retiredName)
	if err := os.MkdirAll(repromoted, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(repromoted, "SKILL.md"), []byte("name: "+retiredName+"\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(targetDir, 0o755); err != nil {
		t.Fatal(err)
	}
	target := filepath.Join(targetDir, retiredName)
	if err := os.Symlink(repromoted, target); err != nil {
		t.Fatal(err)
	}

	var out bytes.Buffer
	config := Config{Out: &out}
	if err := config.createSkillSymlinks(sourceDir, targetDir, []string{retiredName}); err != nil {
		t.Fatal(err)
	}
	if link, err := os.Readlink(target); err != nil || link != repromoted {
		t.Fatalf("repromoted skill symlink changed: link=%q err=%v", link, err)
	}
	if strings.Contains(out.String(), "Removed retired repo skill symlink") {
		t.Fatalf("repromoted skill was logged as retired:\n%s", out.String())
	}
}

func TestCleanupRetiredSymlinkUsesCanonicalSourceParent(t *testing.T) {
	dir := t.TempDir()
	repoDir := filepath.Join(dir, "repo")
	sourceDir := filepath.Join(repoDir, "skills")
	aliasRepo := filepath.Join(dir, "repo-alias")
	targetDir := filepath.Join(dir, "home", "skills")
	retiredName := "retired-alpha"
	if err := os.MkdirAll(sourceDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(repoDir, aliasRepo); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(targetDir, 0o755); err != nil {
		t.Fatal(err)
	}
	target := filepath.Join(targetDir, retiredName)
	if err := os.Symlink(filepath.Join("..", "..", "repo-alias", "skills", retiredName), target); err != nil {
		t.Fatal(err)
	}

	var out bytes.Buffer
	config := Config{Out: &out}
	if err := config.cleanupRetiredSkillSymlinks(sourceDir, targetDir, []string{retiredName}); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Lstat(target); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("canonical repo-owned alias was not removed, err=%v", err)
	}
}

func TestInstallRejectsInvalidRetiredManifestBeforeProviderMutation(t *testing.T) {
	repo := t.TempDir()
	if err := os.MkdirAll(filepath.Join(repo, "config"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(repo, "config", "retired-skills.json"), []byte(`{"retired_skills": ["Invalid"]}`), 0o644); err != nil {
		t.Fatal(err)
	}
	codexHome := filepath.Join(t.TempDir(), "codex")
	if err := os.MkdirAll(codexHome, 0o755); err != nil {
		t.Fatal(err)
	}
	config := Config{
		RepoRoot:        repo,
		CodexHome:       codexHome,
		ClaudeHome:      filepath.Join(t.TempDir(), "claude"),
		OpenCodeHome:    filepath.Join(t.TempDir(), "opencode"),
		AntigravityHome: filepath.Join(t.TempDir(), "antigravity"),
		AgentsHome:      filepath.Join(t.TempDir(), "agents"),
		Out:             &bytes.Buffer{},
	}
	if err := config.Install(); err == nil {
		t.Fatal("expected invalid retired manifest to stop install")
	}
	if _, err := os.Lstat(filepath.Join(codexHome, "AGENTS.md")); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("provider target changed before manifest validation, err=%v", err)
	}
}

func TestSetupCodexAgentsRejectsInvalidRetiredManifestBeforeMutation(t *testing.T) {
	repo := t.TempDir()
	if err := os.MkdirAll(filepath.Join(repo, "config"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(repo, "config", "retired-skills.json"), []byte(`{"retired_skills": ["Invalid"]}`), 0o644); err != nil {
		t.Fatal(err)
	}
	codexHome := filepath.Join(t.TempDir(), "codex")
	config := Config{RepoRoot: repo, CodexHome: codexHome, AgentsHome: filepath.Join(t.TempDir(), "agents"), Out: &bytes.Buffer{}}
	if err := config.SetupCodexAgents(); err == nil {
		t.Fatal("expected invalid retired manifest to stop Codex setup")
	}
	if _, err := os.Lstat(codexHome); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("Codex provider root changed before manifest validation, err=%v", err)
	}
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
	staleManagedTarget := filepath.Join(targetDir, "stale-managed.toml")
	if err := os.WriteFile(staleManagedTarget, []byte(managedHeader+"name = \"stale-managed\"\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	staleUserTarget := filepath.Join(targetDir, "stale-user.toml")
	if err := os.WriteFile(staleUserTarget, []byte("user-owned stale role"), 0o644); err != nil {
		t.Fatal(err)
	}
	staleForeignLink := filepath.Join(targetDir, "stale-link.toml")
	if err := os.Symlink(filepath.Join(repo, "foreign-role.toml"), staleForeignLink); err != nil {
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
	if _, err := os.Stat(staleManagedTarget); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("expected stale managed role removed, err=%v", err)
	}
	if got := readFile(t, staleUserTarget); got != "user-owned stale role" {
		t.Fatalf("stale user-owned role file changed: %q", got)
	}
	if _, err := os.Lstat(staleForeignLink); err != nil {
		t.Fatalf("stale foreign role link was removed: %v", err)
	}
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
      "args": ["-y", "mcp-remote", "https://mcp.clickup.com/mcp"],
      "enabled_tools": ["clickup_search", "clickup_update_task"]
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
	mustContain(t, got, `enabled_tools = ["clickup_search", "clickup_update_task"]`)
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

func TestGenerateCodexRoleFilesDoesNotFollowReplacementSymlink(t *testing.T) {
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
	writeRoleManifest(t, sourceDir, `{"skill_groups":{"demo":["{{REPO_ROOT}}/skills/demo"]},"roles":{"demo":{}}}`)
	external := filepath.Join(repo, "external.toml")
	if err := os.WriteFile(external, []byte("external-owned"), 0o644); err != nil {
		t.Fatal(err)
	}
	target := filepath.Join(targetDir, "demo.toml")
	config := Config{
		RepoRoot: repo,
		Out:      &bytes.Buffer{},
		beforeGeneratedRoleWrite: func(path string) {
			if err := os.Symlink(external, path); err != nil {
				t.Fatal(err)
			}
		},
	}
	if err := config.generateCodexRoleFiles(targetDir); err != nil {
		t.Fatal(err)
	}
	if got := readFile(t, external); got != "external-owned" {
		t.Fatalf("external replacement was modified: %q", got)
	}
	if link, err := os.Readlink(target); err != nil || link != external {
		t.Fatalf("replacement symlink was not preserved: link=%q err=%v", link, err)
	}
}

func TestGenerateCodexRoleFilesDoesNotRemoveReplacementRegularStaleFile(t *testing.T) {
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
	writeRoleManifest(t, sourceDir, `{"skill_groups":{"demo":["{{REPO_ROOT}}/skills/demo"]},"roles":{"demo":{}}}`)
	stale := filepath.Join(targetDir, "stale.toml")
	if err := os.WriteFile(stale, []byte(managedHeader+"name = \"stale\"\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	config := Config{
		RepoRoot: repo,
		Out:      &bytes.Buffer{},
		beforeStaleRoleRemove: func(path string) {
			if err := os.Remove(path); err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(path, []byte("replacement-owned"), 0o644); err != nil {
				t.Fatal(err)
			}
		},
	}
	if err := config.generateCodexRoleFiles(targetDir); err != nil {
		t.Fatal(err)
	}
	if got := readFile(t, stale); got != "replacement-owned" {
		t.Fatalf("replacement regular file was removed or changed: %q", got)
	}
}

func TestCleanupRetiredSkillSymlinksDoesNotRemoveReplacementRegularFile(t *testing.T) {
	dir := t.TempDir()
	sourceDir := filepath.Join(dir, "repo", "skills")
	targetDir := filepath.Join(dir, "home", "skills")
	retiredName := "retired-alpha"
	if err := os.MkdirAll(sourceDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(targetDir, 0o755); err != nil {
		t.Fatal(err)
	}
	target := filepath.Join(targetDir, retiredName)
	if err := os.Symlink(filepath.Join(sourceDir, retiredName), target); err != nil {
		t.Fatal(err)
	}
	config := Config{
		Out: &bytes.Buffer{},
		beforeRetiredSkillRemove: func(path string) {
			if err := os.Remove(path); err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(path, []byte("replacement-owned"), 0o644); err != nil {
				t.Fatal(err)
			}
		},
	}
	if err := config.cleanupRetiredSkillSymlinks(sourceDir, targetDir, []string{retiredName}); err != nil {
		t.Fatal(err)
	}
	if got := readFile(t, target); got != "replacement-owned" {
		t.Fatalf("replacement regular file was removed or changed: %q", got)
	}
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
