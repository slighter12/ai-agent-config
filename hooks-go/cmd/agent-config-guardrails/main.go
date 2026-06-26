package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"ai-agent-config/hooks/internal/hookjson"
)

func main() {
	if len(os.Args) < 2 {
		os.Exit(0)
	}
	raw, err := io.ReadAll(os.Stdin)
	if err != nil {
		os.Exit(0)
	}
	event, _ := hookjson.ParseEvent(raw)

	switch os.Args[1] {
	case "pre-tool-use-bash":
		runBash(event)
	case "pre-tool-use-file":
		runFile(event)
	default:
		os.Exit(0)
	}
}

func runBash(event hookjson.Event) {
	if hookjson.ToolName(event) != "Bash" {
		return
	}
	commandText := hookjson.CommandText(event)
	trimmedCommand := strings.TrimSpace(commandText)
	if trimmedCommand == "" || strings.HasPrefix(trimmedCommand, "rtk ") {
		return
	}
	rtkPath, err := exec.LookPath("rtk")
	if err != nil {
		return
	}
	if trimmedCommand == "rg" || strings.HasPrefix(trimmedCommand, "rg ") {
		writeBashRewrite("rtk " + trimmedCommand)
		return
	}
	cmd := exec.Command(rtkPath, "rewrite", commandText)
	output, err := cmd.Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); !ok || exitErr.ExitCode() != 3 {
			return
		}
	}
	rewritten := strings.TrimRight(string(output), "\r\n")
	if rewritten == "" || rewritten == commandText {
		return
	}
	writeBashRewrite(rewritten)
}

func writeBashRewrite(command string) {
	write(map[string]any{
		"hookSpecificOutput": map[string]any{
			"hookEventName":      "PreToolUse",
			"permissionDecision": "allow",
			"updatedInput": map[string]any{
				"command": command,
			},
		},
	})
}

func runFile(event hookjson.Event) {
	decision := hookjson.ClassifyPreToolUse(event)
	if decision.Behavior != hookjson.Deny {
		return
	}
	write(map[string]any{
		"hookSpecificOutput": map[string]any{
			"hookEventName":            "PreToolUse",
			"permissionDecision":       "deny",
			"permissionDecisionReason": decision.Reason,
		},
	})
}

func write(value any) {
	fmt.Println(string(hookjson.MustJSON(value)))
}
