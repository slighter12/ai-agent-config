package integration_test

import (
	"path/filepath"
	"strings"
	"testing"
)

// Codex role files inline their instructions because the installer copies them
// verbatim; agents/<role>.md stays the single editable source of that text.
func TestCodexRoleInstructionsMatchAgentDefinitions(t *testing.T) {
	repoRoot := filepath.Clean(filepath.Join("..", ".."))
	for _, role := range []string{"builder", "explorer", "librarian", "oracle", "orchestrator", "reviewer", "test-runner"} {
		t.Run(role, func(t *testing.T) {
			roleFile := readFile(t, filepath.Join(repoRoot, "config", "codex-agents", role+".toml"))
			for _, key := range []string{`name = "` + role + `"`, "\nmodel = ", "\nmodel_reasoning_effort = "} {
				mustContain(t, roleFile, key)
			}
			mustNotContain(t, roleFile, "instructions_path")
			// Custom agents inherit these from the parent session; setting them here is silently ignored.
			for _, ignored := range []string{"\ndefault_permissions =", "\nsandbox_mode =", "\napproval_policy =", "[mcp_servers."} {
				mustNotContain(t, roleFile, ignored)
			}

			_, rest, found := strings.Cut(roleFile, "developer_instructions = '''\n")
			if !found {
				t.Fatalf("%s.toml has no literal developer_instructions block", role)
			}
			instructions, _, closed := strings.Cut(rest, "\n'''")
			if !closed {
				t.Fatalf("%s.toml developer_instructions block is not closed", role)
			}

			definition := readFile(t, filepath.Join(repoRoot, "agents", role+".md"))
			if strings.HasPrefix(definition, "---\n") {
				if close := strings.Index(definition[4:], "\n---\n"); close >= 0 {
					definition = definition[4+close+len("\n---\n"):]
				}
			}
			if got, want := strings.TrimSpace(instructions), strings.TrimSpace(definition); got != want {
				t.Fatalf("%s.toml developer_instructions drifted from agents/%s.md; regenerate it from the agent body", role, role)
			}
		})
	}
}
