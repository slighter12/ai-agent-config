# Local Provider Mechanics

This repository keeps one flat, editable skill source and links it into Codex, Claude Code, OpenCode, and Gemini/Antigravity discovery roots.

- Every skill declares `metadata.invocation: user|model`.
- User-invoked skills also set `disable-model-invocation: true`, `metadata["opencode/autoinvoke"]: "false"`, and `agents/openai.yaml` with `policy.allow_implicit_invocation: false`.
- Use the repository's `agent-config` CLI to initialize, validate, and package skills; use absolute skill paths because commands run from `hooks-go`.
- Keep source directories flat as `skills/<skill-name>`; the installer owns provider symlinks.
- Relative links such as `[FILE.md](FILE.md)` point only at files inside the same skill directory. To call another skill, write "use the `<name>` skill": Claude Code and OpenCode invoke the `Skill` tool with `"<name>"`, and Codex opens the `SKILL.md` path its skill list shows for `<name>`. The called skill must be model-invoked; for a user-invoked skill, tell the user to run `/<name>` instead. Never reach into another skill's directory with a parent-directory relative path; a contract test in `hooks-go/integration/catalog_contracts_test.go` rejects any.
- When retiring or merging a skill, add an entry to `config/retired-skills.json` (schema v2: `schemaVersion`, `updatedAt`, `retiredSkills[]` with `name`, `retiredIn`, `replacedBy`, `reason`) and record the divergence in `skills/UPSTREAM.md` without naming the retired skill.
