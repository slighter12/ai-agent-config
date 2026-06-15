package main

import (
	"fmt"
	"io"
	"os"

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
	_ = event
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
