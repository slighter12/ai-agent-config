package integration_test

import (
	"os"
	"path/filepath"
	"sort"
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

// The installer links claude-agents/ while agents/ stays the source the Codex
// roles are checked against, so the two directories must not drift.
func TestAgentDefinitionDirectoriesStayIdentical(t *testing.T) {
	repoRoot := filepath.Clean(filepath.Join("..", ".."))
	roles := func(dir, ext string) []string {
		entries, err := os.ReadDir(filepath.Join(repoRoot, dir))
		if err != nil {
			t.Fatal(err)
		}
		var names []string
		for _, entry := range entries {
			if !entry.IsDir() && filepath.Ext(entry.Name()) == ext {
				names = append(names, strings.TrimSuffix(entry.Name(), ext))
			}
		}
		sort.Strings(names)
		return names
	}
	source := roles("agents", ".md")
	if len(source) == 0 {
		t.Fatal("no agent definitions found")
	}
	for dir, ext := range map[string]string{"claude-agents": ".md", filepath.Join("config", "codex-agents"): ".toml"} {
		if got, want := strings.Join(roles(dir, ext), ","), strings.Join(source, ","); got != want {
			t.Errorf("%s roles = %s, want %s", dir, got, want)
		}
	}
	for _, role := range source {
		agent := readFile(t, filepath.Join(repoRoot, "agents", role+".md"))
		claude := readFile(t, filepath.Join(repoRoot, "claude-agents", role+".md"))
		if agent != claude {
			t.Errorf("claude-agents/%s.md differs from agents/%s.md; copy it byte for byte", role, role)
		}
	}
}

// Role definitions call skills by name and stay inside their authority:
// read-only roles get read-only skills, and no role loads verify-app or the
// user-invoked step-back.
func TestAgentDefinitionsCallSkillsByName(t *testing.T) {
	repoRoot := filepath.Clean(filepath.Join("..", ".."))
	read := func(role string) string {
		return readFile(t, filepath.Join(repoRoot, "agents", role+".md"))
	}
	mustContain(t, read("explorer"), "use the `how` skill")
	oracle := read("oracle")
	mustContain(t, oracle, "Use the `codebase-design` skill")
	mustContain(t, oracle, "use the `why` skill")
	mustContain(t, oracle, "use the `blast-radius` skill")
	mustNotContain(t, oracle, "(`codebase-design`)")
	reviewer := read("reviewer")
	mustContain(t, reviewer, "use the `why` skill")
	mustContain(t, reviewer, "use the `blast-radius` skill")
	builder := read("builder")
	mustContain(t, builder, "Use the `tidy` skill only when the ticket or the user explicitly asks for a tidy step")
	mustContain(t, builder, "Use the `tdd` skill")
	mustNotContain(t, builder, "(`tdd`)")
	orchestrator := read("orchestrator")
	mustContain(t, orchestrator, "Use the `code-review` skill")
	mustNotContain(t, orchestrator, "Run `code-review`")
	for _, role := range []string{"builder", "explorer", "librarian", "oracle", "orchestrator", "reviewer", "test-runner"} {
		definition := read(role)
		for _, banned := range []string{"verify-app", "step-back", "../"} {
			mustNotContain(t, definition, banned)
		}
		if role != "builder" {
			mustNotContain(t, definition, "`tidy`")
			mustNotContain(t, definition, "`unslop`")
		}
	}
}
