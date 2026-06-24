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
	ProfileName        = "workspace-git"
	BasePermissionName = ":workspace"
)

var (
	topLevelDefaultPermissionsRE = regexp.MustCompile(`^\s*default_permissions\s*=`)
	topLevelApprovalPolicyRE     = regexp.MustCompile(`^\s*approval_policy\s*=`)
	topLevelApprovalsReviewerRE  = regexp.MustCompile(`^\s*approvals_reviewer\s*=`)
	topLevelSandboxModeRE        = regexp.MustCompile(`^\s*sandbox_mode\s*=`)
	execPermissionApprovalsRE    = regexp.MustCompile(`^\s*exec_permission_approvals\s*=`)
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
	text = upsertBaseDefaultPermissions(text)
	text = upsertApprovalPolicy(text)
	text = upsertApprovalsReviewer(text)
	text = removePermissionFeatureFlags(text)
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

func upsertBaseDefaultPermissions(text string) string {
	lines := strings.Split(strings.TrimRight(text, "\n"), "\n")
	if len(lines) == 1 && lines[0] == "" {
		return `default_permissions = "` + BasePermissionName + `"` + "\n"
	}
	currentTable := ""
	insertIndex := -1
	for index, line := range lines {
		trimmed := strings.TrimSpace(line)
		if table, ok := tableName(trimmed); ok {
			currentTable = table
			if insertIndex == -1 {
				insertIndex = index
			}
		}
		if currentTable == "" && topLevelDefaultPermissionsRE.MatchString(trimmed) {
			if repoManagedDefaultPermissions(trimmed) {
				lines[index] = `default_permissions = "` + BasePermissionName + `"`
			}
			return strings.TrimRight(strings.Join(lines, "\n"), "\n") + "\n"
		}
	}
	if insertIndex == -1 {
		lines = append(lines, `default_permissions = "`+BasePermissionName+`"`)
	} else {
		insertIndex = beforeTopLevelSeparator(lines, insertIndex)
		before := append([]string{}, lines[:insertIndex]...)
		before = append(before, `default_permissions = "`+BasePermissionName+`"`)
		lines = append(before, lines[insertIndex:]...)
	}
	return strings.TrimRight(strings.Join(lines, "\n"), "\n") + "\n"
}

func repoManagedDefaultPermissions(line string) bool {
	if !topLevelDefaultPermissionsRE.MatchString(line) {
		return false
	}
	parts := strings.SplitN(line, "=", 2)
	if len(parts) != 2 {
		return false
	}
	value := strings.TrimSpace(parts[1])
	if index := strings.Index(value, "#"); index >= 0 {
		value = strings.TrimSpace(value[:index])
	}
	return value == `"`+ProfileName+`"` || value == `":workspace"`
}

func upsertApprovalPolicy(text string) string {
	return upsertTopLevelString(text, topLevelApprovalPolicyRE, "approval_policy", "on-request", repoManagedApprovalPolicy)
}

func repoManagedApprovalPolicy(line string) bool {
	value, ok := topLevelStringValue(line)
	return ok && (value == "on-request" || value == "")
}

func upsertApprovalsReviewer(text string) string {
	return upsertTopLevelString(text, topLevelApprovalsReviewerRE, "approvals_reviewer", "auto_review", repoManagedApprovalsReviewer)
}

func repoManagedApprovalsReviewer(line string) bool {
	value, ok := topLevelStringValue(line)
	return ok && (value == "user" || value == "auto_review" || value == "")
}

func upsertTopLevelString(text string, lineRE *regexp.Regexp, key, value string, shouldReplace func(string) bool) string {
	lines := strings.Split(strings.TrimRight(text, "\n"), "\n")
	if len(lines) == 1 && lines[0] == "" {
		return key + ` = "` + value + `"` + "\n"
	}
	currentTable := ""
	insertIndex := -1
	for index, line := range lines {
		trimmed := strings.TrimSpace(line)
		if table, ok := tableName(trimmed); ok {
			currentTable = table
			if insertIndex == -1 {
				insertIndex = index
			}
		}
		if currentTable == "" && lineRE.MatchString(trimmed) {
			if shouldReplace(trimmed) {
				lines[index] = key + ` = "` + value + `"`
			}
			return strings.TrimRight(strings.Join(lines, "\n"), "\n") + "\n"
		}
	}
	if insertIndex == -1 {
		lines = append(lines, key+` = "`+value+`"`)
	} else {
		insertIndex = beforeTopLevelSeparator(lines, insertIndex)
		before := append([]string{}, lines[:insertIndex]...)
		before = append(before, key+` = "`+value+`"`)
		lines = append(before, lines[insertIndex:]...)
	}
	return strings.TrimRight(strings.Join(lines, "\n"), "\n") + "\n"
}

func beforeTopLevelSeparator(lines []string, index int) int {
	for index > 0 && strings.TrimSpace(lines[index-1]) == "" {
		index--
	}
	return index
}

func topLevelStringValue(line string) (string, bool) {
	parts := strings.SplitN(line, "=", 2)
	if len(parts) != 2 {
		return "", false
	}
	value := strings.TrimSpace(parts[1])
	if index := strings.Index(value, "#"); index >= 0 {
		value = strings.TrimSpace(value[:index])
	}
	return strings.Trim(value, `"`), true
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

func removePermissionFeatureFlags(text string) string {
	lines := strings.Split(strings.TrimRight(text, "\n"), "\n")
	if len(lines) == 1 && lines[0] == "" {
		return ""
	}

	featuresStart := -1
	featuresEnd := len(lines)
	for index, line := range lines {
		if table, ok := tableName(strings.TrimSpace(line)); ok {
			if featuresStart != -1 {
				featuresEnd = index
				break
			}
			if table == "features" {
				featuresStart = index
			}
		}
	}

	if featuresStart == -1 {
		return strings.TrimRight(strings.Join(lines, "\n"), "\n") + "\n"
	}

	filtered := append([]string{}, lines[:featuresStart+1]...)
	for index := featuresStart + 1; index < featuresEnd; index++ {
		trimmed := strings.TrimSpace(lines[index])
		if execPermissionApprovalsRE.MatchString(trimmed) {
			continue
		}
		if strings.HasPrefix(trimmed, "request_permissions_tool") {
			continue
		}
		filtered = append(filtered, lines[index])
	}

	updated := append(filtered, lines[featuresEnd:]...)
	return strings.TrimRight(strings.Join(updated, "\n"), "\n") + "\n"
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
