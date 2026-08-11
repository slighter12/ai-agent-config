# Local Provider Mechanics

This repository keeps one flat, editable skill source and links it into Codex, Claude Code, OpenCode, and Gemini/Antigravity discovery roots.

- Every skill declares `metadata.invocation: user|model`.
- User-invoked skills also set `disable-model-invocation: true`, `metadata["opencode/autoinvoke"]: "false"`, and `agents/openai.yaml` with `policy.allow_implicit_invocation: false`.
- Use the repository's `agent-config` CLI to initialize, validate, and package skills; use absolute skill paths because commands run from `hooks-go`.
- Keep source directories flat as `skills/<skill-name>`; the installer owns provider symlinks.
