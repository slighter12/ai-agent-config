package hookjson

import (
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
