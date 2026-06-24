package main

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestRunBashRewritesWithFakeRTK(t *testing.T) {
	tempDir := t.TempDir()
	writeFakeRTK(t, tempDir)
	output := runBashHook(t, tempDir, `{"tool_name":"Bash","tool_input":{"command":"git status"}}`)

	var got struct {
		HookSpecificOutput struct {
			PermissionDecision string `json:"permissionDecision"`
			UpdatedInput       struct {
				Command string `json:"command"`
			} `json:"updatedInput"`
		} `json:"hookSpecificOutput"`
	}
	if err := json.Unmarshal([]byte(output), &got); err != nil {
		t.Fatalf("output is not JSON: %v\n%s", err, output)
	}
	if got.HookSpecificOutput.PermissionDecision != "allow" {
		t.Fatalf("expected allow decision, got %q", got.HookSpecificOutput.PermissionDecision)
	}
	if got.HookSpecificOutput.UpdatedInput.Command != "rtk git status" {
		t.Fatalf("expected rewritten command, got %q", got.HookSpecificOutput.UpdatedInput.Command)
	}
}

func TestRunBashSkipsAlreadyRTKCommand(t *testing.T) {
	tempDir := t.TempDir()
	writeFakeRTK(t, tempDir)
	output := runBashHook(t, tempDir, `{"tool_name":"Bash","tool_input":{"command":"rtk git status"}}`)
	if strings.TrimSpace(output) != "" {
		t.Fatalf("expected no rewrite output, got: %s", output)
	}
}

func TestRunBashDoesNothingWhenRTKMissing(t *testing.T) {
	output := runBashHookWithPath(t, t.TempDir(), `{"tool_name":"Bash","tool_input":{"command":"git status"}}`)
	if strings.TrimSpace(output) != "" {
		t.Fatalf("expected no rewrite output, got: %s", output)
	}
}

func writeFakeRTK(t *testing.T, dir string) {
	t.Helper()
	rtkPath := filepath.Join(dir, "rtk")
	script := "#!/usr/bin/env sh\nif [ \"$1\" = rewrite ] && [ \"$2\" = \"git status\" ]; then printf '%s\\n' 'rtk git status'; exit 3; fi\nexit 0\n"
	if err := os.WriteFile(rtkPath, []byte(script), 0o755); err != nil {
		t.Fatal(err)
	}
}

func runBashHook(t *testing.T, pathDir string, stdin string) string {
	t.Helper()
	return runBashHookWithPath(t, pathDir+string(os.PathListSeparator)+os.Getenv("PATH"), stdin)
}

func runBashHookWithPath(t *testing.T, path string, stdin string) string {
	t.Helper()
	cmd := exec.Command("go", "run", ".", "pre-tool-use-bash")
	cmd.Dir = "."
	cmd.Env = append(os.Environ(), "PATH="+path)
	cmd.Stdin = strings.NewReader(stdin)
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("command failed: %v\n%s", err, output)
	}
	return string(output)
}
