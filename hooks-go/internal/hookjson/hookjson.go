package hookjson

import (
	"encoding/json"
	"path/filepath"
	"regexp"
	"strings"
)

const (
	Allow = "allow"
	Deny  = "deny"
)

var globalConfigPaths = []string{
	".codex/config.toml",
	".codex/hooks.json",
	".claude/settings.json",
}

type Event map[string]any

type Decision struct {
	Behavior string
	Reason   string
}

func ParseEvent(raw []byte) (Event, bool) {
	if len(strings.TrimSpace(string(raw))) == 0 {
		return Event{}, true
	}
	var event Event
	if err := json.Unmarshal(raw, &event); err != nil {
		return Event{"_hook_parse_error": true}, false
	}
	if event == nil {
		return Event{}, true
	}
	return event, true
}

func ToolName(event Event) string {
	for _, key := range []string{"tool_name", "toolName", "tool"} {
		if value, ok := event[key].(string); ok {
			return value
		}
	}
	return ""
}

func ToolInput(event Event) map[string]any {
	for _, key := range []string{"tool_input", "toolInput", "input"} {
		if value, ok := event[key].(map[string]any); ok {
			return value
		}
	}
	return map[string]any{}
}

func CommandText(event Event) string {
	input := ToolInput(event)
	for _, key := range []string{"command", "cmd"} {
		if value, ok := input[key].(string); ok {
			return value
		}
	}
	if value, ok := event["command"].(string); ok {
		return value
	}
	return ""
}

func ClassifyFileTool(event Event) Decision {
	for _, rawPath := range TouchedPaths(event) {
		name := filepath.Base(rawPath)
		if name == ".env" || strings.HasPrefix(name, ".env.") {
			return Decision{Behavior: Deny, Reason: "Blocked direct edit/write to `.env*` file."}
		}
		normalized := strings.ReplaceAll(rawPath, "\\", "/")
		for _, path := range globalConfigPaths {
			if strings.HasSuffix(normalized, path) {
				return Decision{Behavior: Deny, Reason: "Blocked direct edit/write to global agent configuration."}
			}
		}
	}
	return Decision{Behavior: Allow}
}

func ClassifyPreToolUse(event Event) Decision {
	if parseErr, _ := event["_hook_parse_error"].(bool); parseErr {
		return Decision{Behavior: Deny, Reason: "Hook input was not valid JSON; blocked conservatively."}
	}
	switch ToolName(event) {
	case "Write", "Edit", "MultiEdit", "apply_patch":
		return ClassifyFileTool(event)
	default:
		return Decision{Behavior: Allow}
	}
}

func TouchedPaths(event Event) []string {
	input := ToolInput(event)
	seen := map[string]struct{}{}
	addPath := func(value any) {
		text, ok := value.(string)
		if !ok {
			return
		}
		text = strings.TrimSpace(text)
		if text == "" {
			return
		}
		seen[text] = struct{}{}
	}
	for _, key := range []string{"file_path", "filePath", "path", "target_file", "targetFile"} {
		addPath(input[key])
	}
	for _, text := range WalkStrings(input) {
		for _, match := range patchFileRE.FindAllStringSubmatch(text, -1) {
			addPath(match[1])
		}
		for _, match := range patchMoveRE.FindAllStringSubmatch(text, -1) {
			addPath(match[1])
		}
	}
	paths := make([]string, 0, len(seen))
	for path := range seen {
		paths = append(paths, path)
	}
	return paths
}

var (
	patchFileRE = regexp.MustCompile(`(?m)^\*\*\* (?:Add|Update|Delete) File: (.+)$`)
	patchMoveRE = regexp.MustCompile(`(?m)^\*\*\* Move to: (.+)$`)
)

func WalkStrings(value any) []string {
	switch typed := value.(type) {
	case string:
		return []string{typed}
	case []any:
		var result []string
		for _, item := range typed {
			result = append(result, WalkStrings(item)...)
		}
		return result
	case map[string]any:
		var result []string
		for _, item := range typed {
			result = append(result, WalkStrings(item)...)
		}
		return result
	default:
		return nil
	}
}

func MustJSON(value any) []byte {
	payload, err := json.Marshal(value)
	if err != nil {
		panic(err)
	}
	return payload
}
