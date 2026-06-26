package hookjson

import (
	"strings"
	"testing"
)

func TestClassifyFileToolDotenvWriteDenies(t *testing.T) {
	decision := ClassifyFileTool(Event{
		"tool_input": map[string]any{"file_path": ".env.local"},
	})
	if decision.Behavior != Deny {
		t.Fatalf("expected deny, got %q", decision.Behavior)
	}
}

func TestClassifyFileToolGlobalConfigWriteDenies(t *testing.T) {
	decision := ClassifyFileTool(Event{
		"tool_input": map[string]any{"file_path": "~/.codex/config.toml"},
	})
	if decision.Behavior != Deny {
		t.Fatalf("expected deny, got %q", decision.Behavior)
	}
}

func TestClassifyPreToolUseInvalidJSONDenies(t *testing.T) {
	decision := ClassifyPreToolUse(Event{"_hook_parse_error": true})
	if decision.Behavior != Deny {
		t.Fatalf("expected deny, got %q", decision.Behavior)
	}
}

func TestSimpleGitActionPromptMatchesIndependentActions(t *testing.T) {
	prompts := []string{
		"幫我創建分支",
		"先幫我 commit",
		"幫我 push",
		"幫我發 PR",
		"幫我創建分支 commit 並且發 PR",
	}
	for _, prompt := range prompts {
		if !SimpleGitActionPrompt(prompt) {
			t.Fatalf("expected prompt to match: %q", prompt)
		}
	}
}

func TestSimpleGitActionPromptIgnoresTextOnlyPrompts(t *testing.T) {
	prompts := []string{
		"幫我看一下這段設定",
		"幫我想 commit message",
		"幫我取 branch name",
		"幫我看一下 PR",
		"幫我寫 PR title/body",
	}
	for _, prompt := range prompts {
		if SimpleGitActionPrompt(prompt) {
			t.Fatalf("expected prompt not to match: %q", prompt)
		}
	}
}

func TestGitWorkflowReminderDelegatesWithMainSessionFallback(t *testing.T) {
	decision := ClassifyUserPromptSubmit(Event{"prompt": "幫我創建分支 commit 並且發 PR"})
	if decision.Behavior != Allow {
		t.Fatalf("expected allow, got %q", decision.Behavior)
	}
	for _, want := range []string{
		"standing delegation request",
		"`git-commit` role",
		"discover the provider's callable spawn/delegation tool first",
		"fallback only when no callable route exists after discovery",
		"runtime policy blocks delegation before it starts",
		"continue in the main session",
		"report that fallback",
		"Once delegation starts, the main session must not run mutating git commands",
	} {
		if !strings.Contains(decision.Reason, want) {
			t.Fatalf("expected reminder to contain %q, got %q", want, decision.Reason)
		}
	}
	for _, forbidden := range []string{
		"stop and report that blocker",
		"fail closed",
	} {
		if strings.Contains(decision.Reason, forbidden) {
			t.Fatalf("expected reminder not to contain %q, got %q", forbidden, decision.Reason)
		}
	}
}
