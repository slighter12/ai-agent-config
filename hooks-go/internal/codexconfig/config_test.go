package codexconfig

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestEnsureWorkspaceGitProfileCreatesProfile(t *testing.T) {
	got, err := EnsureWorkspaceGitProfile("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	mustContain(t, got, `default_permissions = ":workspace"`)
	mustContain(t, got, `approval_policy = "on-request"`)
	mustContain(t, got, `approvals_reviewer = "auto_review"`)
	mustNotContain(t, got, `exec_permission_approvals`)
	mustNotContain(t, got, `request_permissions_tool`)
	mustContain(t, got, `[permissions.workspace-git]`)
	mustContain(t, got, `".git" = "write"`)
	mustContain(t, got, `[permissions.workspace-git.network]`)
	mustContain(t, got, `"github.com" = "allow"`)
}

func TestEnsureWorkspaceGitProfileIsIdempotent(t *testing.T) {
	first, err := EnsureWorkspaceGitProfile(`model = "gpt-5.5"

[projects."/tmp/demo"]
trust_level = "trusted"
`)
	if err != nil {
		t.Fatalf("first apply: %v", err)
	}
	second, err := EnsureWorkspaceGitProfile(first)
	if err != nil {
		t.Fatalf("second apply: %v", err)
	}
	if first != second {
		t.Fatalf("expected idempotent output\nfirst:\n%s\nsecond:\n%s", first, second)
	}
}

func TestEnsureWorkspaceGitProfileReplacesExistingBlock(t *testing.T) {
	got, err := EnsureWorkspaceGitProfile(`default_permissions = "workspace-git"

[permissions.workspace-git]
description = "old"

[permissions.workspace-git.filesystem.":workspace_roots"]
".git" = "read"

[permissions.workspace-git.network]
enabled = false

[hooks.state]
`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if strings.Contains(got, `description = "old"`) || strings.Contains(got, `".git" = "read"`) {
		t.Fatalf("old workspace-git block was not removed:\n%s", got)
	}
	if strings.Contains(got, `enabled = false`) {
		t.Fatalf("old workspace-git network block was not removed:\n%s", got)
	}
	mustContain(t, got, `default_permissions = ":workspace"`)
	mustContain(t, got, `enabled = true`)
	mustContain(t, got, `"api.github.com" = "allow"`)
	mustContain(t, got, `[hooks.state]`)
}

func TestEnsureWorkspaceGitProfileConvertsRepoManagedDefault(t *testing.T) {
	got, err := EnsureWorkspaceGitProfile(`default_permissions = "workspace-git"
approvals_reviewer = "user"

[projects."/tmp/demo"]
trust_level = "trusted"
`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	mustContain(t, got, `default_permissions = ":workspace"`)
	mustContain(t, got, `approval_policy = "on-request"`)
	mustContain(t, got, `approvals_reviewer = "auto_review"`)
	mustContain(t, got, `[permissions.workspace-git]`)
	mustContain(t, got, `[projects."/tmp/demo"]`)
}

func TestEnsureWorkspaceGitProfileRemovesPermissionFeatureFlags(t *testing.T) {
	got, err := EnsureWorkspaceGitProfile(`model = "gpt-5.5"

[features]
exec_permission_approvals = false
terminal_resize_reflow = true
request_permissions_tool = true

[projects."/tmp/demo"]
trust_level = "trusted"
`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	mustNotContain(t, got, `exec_permission_approvals`)
	mustContain(t, got, `terminal_resize_reflow = true`)
	mustContain(t, got, `[projects."/tmp/demo"]`)
	mustNotContain(t, got, `request_permissions_tool`)
}

func TestEnsureWorkspaceGitProfilePreservesUserDefaultPermissions(t *testing.T) {
	got, err := EnsureWorkspaceGitProfile(`default_permissions = ":read-only"

[projects."/tmp/demo"]
trust_level = "trusted"
`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	mustContain(t, got, `default_permissions = ":read-only"`)
	mustContain(t, got, `[permissions.workspace-git]`)
	mustContain(t, got, `[projects."/tmp/demo"]`)
}

func TestEnsureWorkspaceGitProfilePreservesUnrelatedTables(t *testing.T) {
	got, err := EnsureWorkspaceGitProfile(`model = "gpt-5.5"

[hooks.state."demo"]
enabled = true

[mcp_servers.node_repl]
command = "node"
`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	mustContain(t, got, `[hooks.state."demo"]`)
	mustContain(t, got, `[mcp_servers.node_repl]`)
	mustContain(t, got, `command = "node"`)
}

func TestEnsureWorkspaceGitProfileRejectsLegacySandbox(t *testing.T) {
	_, err := EnsureWorkspaceGitProfile(`sandbox_mode = "workspace-write"
`)
	if err == nil {
		t.Fatal("expected legacy sandbox_mode error")
	}
	_, err = EnsureWorkspaceGitProfile(`[sandbox_workspace_write]
network_access = false
`)
	if err == nil {
		t.Fatal("expected legacy sandbox_workspace_write error")
	}
}

func TestApplyFileBacksUpOnlyWhenChanged(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.toml")
	if err := os.WriteFile(path, []byte("model = \"gpt-5.5\"\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	now := time.Date(2026, 6, 4, 12, 0, 0, 0, time.UTC)
	result, err := ApplyFile(path, now)
	if err != nil {
		t.Fatalf("apply: %v", err)
	}
	if !result.Changed || result.BackupPath == "" {
		t.Fatalf("expected changed result with backup, got %+v", result)
	}
	if _, err := os.Stat(result.BackupPath); err != nil {
		t.Fatalf("backup missing: %v", err)
	}
	result, err = ApplyFile(path, now.Add(time.Second))
	if err != nil {
		t.Fatalf("second apply: %v", err)
	}
	if result.Changed || result.BackupPath != "" {
		t.Fatalf("expected second apply to be no-op, got %+v", result)
	}
}

func mustContain(t *testing.T, text, want string) {
	t.Helper()
	if !strings.Contains(text, want) {
		t.Fatalf("expected output to contain %q:\n%s", want, text)
	}
}

func mustNotContain(t *testing.T, text, want string) {
	t.Helper()
	if strings.Contains(text, want) {
		t.Fatalf("expected output not to contain %q:\n%s", want, text)
	}
}
