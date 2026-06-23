package main

import (
	"fmt"
	"io"
	"os"

	"ai-agent-config/hooks/internal/hookjson"
)

func main() {
	raw, err := io.ReadAll(os.Stdin)
	if err != nil {
		os.Exit(0)
	}
	event, _ := hookjson.ParseEvent(raw)
	decision := hookjson.ClassifyUserPromptSubmit(event)
	if decision.Reason == "" {
		return
	}
	fmt.Println(string(hookjson.MustJSON(map[string]any{
		"hookSpecificOutput": map[string]any{
			"hookEventName":     "UserPromptSubmit",
			"additionalContext": decision.Reason,
		},
	})))
}
