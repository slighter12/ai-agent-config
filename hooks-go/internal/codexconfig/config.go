package codexconfig

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

const (
	ProfileName = "workspace-git"
)

var (
	topLevelDefaultPermissionsRE = regexp.MustCompile(`^\s*default_permissions\s*=`)
	topLevelSandboxModeRE        = regexp.MustCompile(`^\s*sandbox_mode\s*=`)
)

type Result struct {
	Path       string
	BackupPath string
	Changed    bool
	Created    bool
}

func DefaultCodexHome() (string, error) {
	if value := strings.TrimSpace(os.Getenv("CODEX_HOME")); value != "" {
		return value, nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".codex"), nil
}

func Apply(codexHome string, now time.Time) (Result, error) {
	if strings.TrimSpace(codexHome) == "" {
		var err error
		codexHome, err = DefaultCodexHome()
		if err != nil {
			return Result{}, err
		}
	}
	path := filepath.Join(codexHome, "config.toml")
	return ApplyFile(path, now)
}

func ApplyFile(path string, now time.Time) (Result, error) {
	result := Result{Path: path}
	originalBytes, err := os.ReadFile(path)
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return result, err
		}
		result.Created = true
	}
	original := string(originalBytes)
	updated, err := EnsureWorkspaceGitProfile(original)
	if err != nil {
		return result, err
	}
	if original == updated {
		return result, nil
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return result, err
	}
	if !result.Created {
		backupPath := fmt.Sprintf("%s.bak-%s", path, now.Format("20060102150405"))
		if err := os.WriteFile(backupPath, originalBytes, 0o600); err != nil {
			return result, err
		}
		result.BackupPath = backupPath
	}
	if err := os.WriteFile(path, []byte(updated), 0o600); err != nil {
		return result, err
	}
	result.Changed = true
	return result, nil
}

func EnsureWorkspaceGitProfile(input string) (string, error) {
	text := normalizeNewline(input)
	if err := rejectLegacySandbox(text); err != nil {
		return "", err
	}
	text = removeWorkspaceGitBlocks(text)
	text = upsertDefaultPermissions(text)
	text = insertWorkspaceGitBlock(text)
	return normalizeNewline(text), nil
}

func rejectLegacySandbox(text string) error {
	currentTable := ""
	for _, line := range strings.Split(text, "\n") {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			continue
		}
		if table, ok := tableName(trimmed); ok {
			currentTable = table
			if table == "sandbox_workspace_write" || strings.HasPrefix(table, "sandbox_workspace_write.") {
				return fmt.Errorf("legacy sandbox table [%s] is present; remove it before using permission profiles", table)
			}
			continue
		}
		if currentTable == "" && topLevelSandboxModeRE.MatchString(trimmed) {
			return errors.New("legacy top-level sandbox_mode is present; remove it before using permission profiles")
		}
	}
	return nil
}

func removeWorkspaceGitBlocks(text string) string {
	lines := strings.Split(strings.TrimRight(text, "\n"), "\n")
	if len(lines) == 1 && lines[0] == "" {
		return ""
	}
	var out []string
	skipping := false
	for _, line := range lines {
		if table, ok := tableName(strings.TrimSpace(line)); ok {
			if table == "permissions."+ProfileName || strings.HasPrefix(table, "permissions."+ProfileName+".") {
				skipping = true
				continue
			}
			skipping = false
		}
		if skipping {
			continue
		}
		out = append(out, line)
	}
	return strings.TrimRight(strings.Join(out, "\n"), "\n") + "\n"
}

func upsertDefaultPermissions(text string) string {
	lines := strings.Split(strings.TrimRight(text, "\n"), "\n")
	if len(lines) == 1 && lines[0] == "" {
		return `default_permissions = "` + ProfileName + `"` + "\n"
	}
	insertIndex := -1
	for index, line := range lines {
		trimmed := strings.TrimSpace(line)
		if table, ok := tableName(trimmed); ok {
			_ = table
			insertIndex = index
			break
		}
		if topLevelDefaultPermissionsRE.MatchString(trimmed) {
			lines[index] = `default_permissions = "` + ProfileName + `"`
			return strings.TrimRight(strings.Join(lines, "\n"), "\n") + "\n"
		}
	}
	if insertIndex == -1 {
		lines = append(lines, `default_permissions = "`+ProfileName+`"`)
	} else {
		before := append([]string{}, lines[:insertIndex]...)
		before = append(before, `default_permissions = "`+ProfileName+`"`)
		lines = append(before, lines[insertIndex:]...)
	}
	return strings.TrimRight(strings.Join(lines, "\n"), "\n") + "\n"
}

func insertWorkspaceGitBlock(text string) string {
	block := workspaceGitBlock()
	lines := strings.Split(strings.TrimRight(text, "\n"), "\n")
	if len(lines) == 1 && lines[0] == "" {
		return block
	}
	insertIndex := len(lines)
	for index, line := range lines {
		table, ok := tableName(strings.TrimSpace(line))
		if ok && strings.HasPrefix(table, "projects.") {
			insertIndex = index
			break
		}
	}
	before := strings.TrimRight(strings.Join(lines[:insertIndex], "\n"), "\n")
	after := strings.TrimLeft(strings.Join(lines[insertIndex:], "\n"), "\n")
	switch {
	case before == "" && after == "":
		return block
	case before == "":
		return block + "\n" + after + "\n"
	case after == "":
		return before + "\n\n" + block
	default:
		return before + "\n\n" + block + "\n" + after + "\n"
	}
}

func workspaceGitBlock() string {
	return `[permissions.workspace-git]
description = "Workspace editing with git metadata writes and GitHub network access."
extends = ":workspace"

[permissions.workspace-git.filesystem.":workspace_roots"]
"." = "write"
".git" = "write"
".codex" = "read"
".agents" = "read"
"**/.env" = "deny"
"**/.env.*" = "deny"
"**/*.env" = "deny"

[permissions.workspace-git.network]
enabled = true

[permissions.workspace-git.network.domains]
"github.com" = "allow"
"api.github.com" = "allow"
"uploads.github.com" = "allow"
`
}

func tableName(line string) (string, bool) {
	if strings.HasPrefix(line, "[[") && strings.HasSuffix(line, "]]") {
		return strings.TrimSpace(strings.TrimSuffix(strings.TrimPrefix(line, "[["), "]]")), true
	}
	if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
		return strings.TrimSpace(strings.TrimSuffix(strings.TrimPrefix(line, "["), "]")), true
	}
	return "", false
}

func normalizeNewline(text string) string {
	text = strings.ReplaceAll(text, "\r\n", "\n")
	text = strings.ReplaceAll(text, "\r", "\n")
	return strings.TrimRight(text, "\n") + "\n"
}
