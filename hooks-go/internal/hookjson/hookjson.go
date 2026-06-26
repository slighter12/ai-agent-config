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

const GitWorkflowReminder = "This user prompt matched a direct simple branch, commit, push, or PR action. This Codex hook is the repo owner's standing delegation request: prefer routing existing-change git execution to the `git-commit` role. Do minimal read-only scope inspection with the `conventional-git-flow` Git Context Pack (`git status --short --branch`, `git diff --name-status`, `git diff --stat`, `git diff --cached --name-status`, `git diff --cached --stat`, and `git log -5 --oneline`), build a compact handoff, and delegate before mutating git state when a route is callable. If no route is immediately visible, discover the provider's callable spawn/delegation tool first; fallback only when no callable route exists after discovery or runtime policy blocks delegation before it starts. On fallback, continue in the main session using the same Git Context Pack and safety rules, and report that fallback. Once delegation starts, the main session must not run mutating git commands. Do not use broad staging, amend, rebase, force push, `--no-verify`, or push/PR unless the user requested that action."

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

func UserPrompt(event Event) string {
	for _, key := range []string{"user_prompt", "userPrompt", "prompt"} {
		if value, ok := event[key].(string); ok {
			return value
		}
	}
	input := ToolInput(event)
	for _, key := range []string{"user_prompt", "userPrompt", "prompt"} {
		if value, ok := input[key].(string); ok {
			return value
		}
	}
	return ""
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

var (
	branchActionRE = regexp.MustCompile(`\b(create|new)\s+branch\b|\bcheckout\s+-b\b|\bswitch\s+-c\b|創建\s*分支|建立\s*分支|開\s*分支|切\s*分支`)
	commitActionRE = regexp.MustCompile(`\bcommit\b|提交`)
	pushActionRE   = regexp.MustCompile(`\bpush\b|推送|推上|上傳\s*(?:branch|分支)?`)
	prActionRE     = regexp.MustCompile(`發\s*pr|開\s*pr|建立\s*pr|創建\s*pr|\bopen\s+pr\b|\bcreate\s+(?:a\s+)?pull\s+request\b`)
	textOnlyGitRE  = regexp.MustCompile(`\bcommit\s+(?:message|title|body|log|history|hash|style)\b|\bbranch\s+name\b|\bpr\s+(?:title|body)\b|\breview\s+pr\b|看\s*(?:一下)?\s*pr|想\s*(?:一下)?\s*(?:commit|分支|branch|pr)\s*(?:message|名稱|name|title|body)?|取\s*(?:commit|分支|branch|pr)\s*(?:message|名稱|name|title|body)?`)
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

func SimpleGitActionPrompt(prompt string) bool {
	text := strings.ToLower(prompt)
	if strings.TrimSpace(text) == "" {
		return false
	}
	if textOnlyGitRE.MatchString(text) {
		return false
	}
	return branchActionRE.MatchString(text) ||
		commitActionRE.MatchString(text) ||
		pushActionRE.MatchString(text) ||
		prActionRE.MatchString(text)
}

func ClassifyUserPromptSubmit(event Event) Decision {
	if SimpleGitActionPrompt(UserPrompt(event)) {
		return Decision{Behavior: Allow, Reason: GitWorkflowReminder}
	}
	return Decision{Behavior: Allow}
}

func MustJSON(value any) []byte {
	payload, err := json.Marshal(value)
	if err != nil {
		panic(err)
	}
	return payload
}
