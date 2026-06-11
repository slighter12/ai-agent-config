package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestRunBashRewritesWithFakeRTK(t *testing.T) {
	tempDir := t.TempDir()
	rtkPath := filepath.Join(tempDir, "rtk")
	if err := os.WriteFile(rtkPath, []byte("#!/usr/bin/env sh\nif [ \"$1\" = rewrite ]; then printf '%s\\n' 'rtk git status'; exit 3; fi\nexit 0\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	cmd := exec.Command("go", "run", ".", "pre-tool-use-bash")
	cmd.Dir = "."
	cmd.Env = append(os.Environ(), "PATH="+tempDir+string(os.PathListSeparator)+os.Getenv("PATH"))
	cmd.Stdin = strings.NewReader(`{"tool_name":"Bash","tool_input":{"command":"git status"}}`)
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("command failed: %v\n%s", err, output)
	}
	if !strings.Contains(string(output), `"command":"rtk git status"`) {
		t.Fatalf("missing rewrite output: %s", output)
	}
}
