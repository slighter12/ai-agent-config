# AI Agent Configuration

A shared, versioned AI configuration repo for local assistants.

## Reference

This repository is inspired by [oh-my-opencode-slim](https://github.com/alvinunreal/oh-my-opencode-slim), then adapted for this workspace's multi-agent and skill-first workflow.

## Scope

This repo is intentionally **skills-first** and **agent-file-first**:

- `skills/`: reusable domain policies and references
- `config/codex-agents/`: Codex runtime role files
- `plugins/`: repo-local Codex plugins
- `.agents/plugins/`: repo-local Codex marketplace metadata
- `agents/`: human-readable role templates and shared role docs
- `claude-agents/`: Claude-oriented role definitions
- `AGENTS.md`: baseline behavior for Codex/OpenCode/Gemini
- `CLAUDE.md`: baseline behavior for Claude Code

`rules/` is no longer used in this repo.

## Repository Layout

```text
ai-config/
├── AGENTS.md
├── CLAUDE.md
├── install.sh
├── .agents/
│   └── plugins/
│       └── marketplace.json
├── claude-agents/
│   ├── orchestrator.md
│   ├── explorer.md
│   ├── builder.md
│   ├── reviewer.md
│   ├── oracle.md
│   ├── codex-collaboration-loop.md
│   ├── librarian.md
│   └── test-runner.md
├── config/
│   └── codex-agents/
│       ├── orchestrator.toml
│       ├── explorer.toml
│       ├── builder.toml
│       ├── reviewer.toml
│       ├── oracle.toml
│       ├── librarian.toml
│       └── test-runner.toml
├── agents/
│   ├── orchestrator.md
│   ├── explorer.md
│   ├── builder.md
│   ├── reviewer.md
│   ├── oracle.md
│   ├── librarian.md
│   └── test-runner.md
├── plugins/
│   └── ai-agent-config-guardrails/
│       ├── .codex-plugin/
│       └── hooks/
│           ├── bin/
│           ├── codex/
│           └── lib/
├── scripts/
│   ├── setup_codex_config.sh
│   ├── setup_codex_agents.sh
│   ├── setup_claude_agents.sh
│   └── setup_codex_shell.sh
├── skills/
│   ├── code-review/
│   ├── conventional-git-flow/
│   ├── diagnose/
│   ├── execution-harness/
│   ├── goal-context/
│   ├── antigravity-design-bridge/
│   ├── project-lifecycle/
│   ├── implement-change/
│   ├── mcp-builder-go/
│   ├── planning-grill/
│   ├── policy-api/
│   ├── policy-core/
│   ├── policy-frontend/
│   ├── policy-go/
│   ├── policy-infra/
│   ├── policy-rust/
│   ├── policy-security/
│   ├── policy-testing/
│   ├── skill-creator/
│   └── verification-driven-change/
```

## Install

Go is required. The shell scripts are thin compatibility wrappers around the Go installer in
`hooks-go/cmd/agent-config`.

```bash
cd ~/git/ai-config
./install.sh
```

The installer creates symlinks (without overwriting existing files):

- Codex (`$CODEX_HOME` or `~/.codex`, if present):
  - `AGENTS.md -> ~/git/ai-config/AGENTS.md`
  - Also runs the Go-backed `setup-codex-config` flow to:
    - upsert the `workspace-git` permission profile in `~/.codex/config.toml`
    - set repo-managed top-level `default_permissions = "workspace-git"` back to `:workspace`
    - back up `config.toml` only when the script changes it
    - fail that setup step if legacy `sandbox_mode` settings would conflict with permission profiles
    - let `install.sh` continue with agents and shell setup if this config step is skipped or fails
  - Also runs the Go-backed `setup-codex-agents` flow to:
    - link shared `skills/* -> ~/.agents/skills/*`
    - generate managed `~/.codex/agents/*.toml` files from `config/codex-agents/*.toml` templates
      and `config/codex-agents/role-manifest.json`
    - warn if a legacy `~/.codex/agents -> ~/git/ai-config/agents` symlink still exists
  - Also runs the Go-backed `setup-codex-shell` flow to:
    - link `shell/codex-profile-auto.zsh -> ~/.codex/shell/codex-profile-auto.zsh`
    - add a guarded source line to `~/.zshrc`
- Claude (`$CLAUDE_HOME` or `~/.claude`, if present):
  - `skills/* -> ~/git/ai-config/skills/*`
  - `CLAUDE.md -> ~/git/ai-config/CLAUDE.md`
  - Also runs the Go-backed `setup-claude-agents` flow to:
    - link `claude-agents/*.md -> ~/.claude/agents/*.md`
    - migrate the repo-managed legacy `~/.claude/agents -> ~/git/ai-config/claude-agents` symlink into a directory
- OpenCode (`$OPENCODE_HOME` or `~/.config/opencode`, if present):
  - `skills/* -> ~/git/ai-config/skills/*`
  - `AGENTS.md -> ~/git/ai-config/AGENTS.md`
- Antigravity CLI (`$ANTIGRAVITY_HOME` or `$GEMINI_HOME` or `~/.gemini`, if present):
  - `antigravity-cli/skills/* -> ~/git/ai-config/skills/*`
  - `GEMINI.md -> ~/git/ai-config/AGENTS.md`

### Optional: RTK command-output compression

RTK is an optional external CLI that compresses verbose command output before it reaches an
AI assistant context. Install it with Homebrew instead of bundling it into `install.sh`, because
it changes system package-manager state rather than this repo's symlinked configuration.

```bash
brew install rtk
rtk --version
rtk gain
```

Expected verification:

- `command -v rtk` resolves to a Homebrew-managed binary such as `/opt/homebrew/bin/rtk`.
- `rtk --version` prints the installed version.
- `rtk gain` prints token-savings stats, or `No tracking data yet.` on a fresh install.

For Codex instructions, initialize RTK explicitly after installation:

```bash
rtk init -g --codex
```

This repo keeps Codex global instructions repo-managed through `~/.codex/AGENTS.md ->
~/git/ai-config/AGENTS.md` and Claude global instructions through `~/.claude/CLAUDE.md ->
~/git/ai-config/CLAUDE.md`. If an RTK init command rewrites either file as a regular file, run
`./install.sh` again or restore the symlink manually. Do not include generated RTK files from this
repo's shared baseline files; Codex and Claude RTK behavior is owned by hook layers. Generated
`~/.codex/RTK.md` and `~/.claude/RTK.md` files may remain as RTK reference files.

RTK binary updates are handled by Homebrew. RTK instruction/template updates are not automatic after
`brew upgrade`; rerun `rtk init -g --codex` only when you intentionally want to refresh the generated
`~/.codex/RTK.md`.

Codex hook integration does not rewrite Bash commands. Earlier versions used RTK as a
`PreToolUse` Bash rewrite source, but command rewriting could change shell semantics such as
`rg --files` becoming a compact grep invocation. RTK may still be used explicitly by typing
`rtk ...`; Claude uses its live `rtk hook claude` configuration when enabled.

Codex hooks are distributed as repo-local plugins. Register the repo marketplace once, then install
or reinstall affected plugins when hook metadata changes:

```bash
codex plugin marketplace add "$(pwd)"
codex plugin add ai-agent-config-guardrails@ai-agent-config
codex plugin add ai-agent-git-routing@ai-agent-config
```

The plugin is self-contained. Do not symlink a repo hook directory into `~/.codex/hooks`; remove
any legacy `~/.codex/hooks` or `~/.codex/hooks.json` symlink before relying on plugin hooks.

## Verify Installation

```bash
ls -la ~/.claude/skills ~/.claude/agents ~/.claude/CLAUDE.md
ls -la ~/.agents/skills
ls -la ~/.codex/agents ~/.codex/shell/codex-profile-auto.zsh ~/.codex/AGENTS.md
codex plugin list --marketplace ai-agent-config
ls -la ~/.config/opencode/skills ~/.config/opencode/AGENTS.md
ls -la ~/.gemini/antigravity-cli/skills ~/.gemini/GEMINI.md
```

If `CODEX_HOME`, `CLAUDE_HOME`, `OPENCODE_HOME`, `ANTIGRAVITY_HOME`, or `GEMINI_HOME` is set, check that directory instead of the default path.

## Usage

### 1) Baseline behavior

Use `AGENTS.md` for Codex/OpenCode/Gemini and `CLAUDE.md` for Claude Code.

### 2) Codex project profiles

`shell/codex-profile-auto.zsh` lets zsh sessions apply repo-specific Codex profiles without
placing config files inside project repositories. When `codex` runs in a git repo, the helper
normalizes the repository directory name into a profile name, then loads
`$CODEX_HOME/<profile>.config.toml` with `--profile` if that file exists.

Project profiles should contain only exceptions over the base `~/.codex/config.toml` policy.
Do not put shared default skill policy in a project profile.

Base skill policy decisions belong in `~/.codex/config.toml`. The `codex:*` helper skills from
the `codex@openai-codex` plugin are disabled in base because they are Claude Code Codex plugin
internals for `codex:codex-rescue` / `codex-companion`, not general Codex session workflows:
`codex:codex-cli-runtime`, `codex:codex-result-handling`, and `codex:gpt-5-4-prompting`.
If a future Claude/Codex-plugin entrypoint needs them, re-enable them through a dedicated
profile or explicit project exception rather than returning them to base.

Artifact skills for `.docx`, `.pptx`, and `.xlsx` work are also disabled in base because they
are large, specialized workflows that are not part of normal coding sessions:
`documents:documents`, `presentations:Presentations`, and `spreadsheets:Spreadsheets`.
Enable them only through a dedicated profile or project exception when a document, deck, or
workbook task actually needs them.

Example: running `codex` inside `Project-D` automatically uses
`~/.codex/project-d.config.toml`. Explicit `--profile` arguments are left unchanged, and
non-runtime management commands such as `codex doctor` are not auto-profiled.

### 3) Skill usage

Call skills by name in prompts, for example:

- `$planning-grill`
- `$implement-change`
- `$diagnose`
- `$verification-driven-change`
- `$policy-go`
- `$policy-frontend`
- `$policy-rust`
- `$policy-security`
- `$code-review`
- `$conventional-git-flow`
- `$execution-harness`
- `$project-lifecycle`
- `$antigravity-design-bridge`
- `$skill-creator`

The former lifecycle skills `phase-closeout`, `sync-decision-docs`, and `session-retrospective`
were merged into `project-lifecycle` and removed as standalone skill entrypoints.

This repo uses `~/.agents/skills` as the canonical shared Codex skill root. `~/.codex/skills`
may still contain system, plugin, legacy, or personal skills. If shared repo skills appear in both
locations, prefer removing the duplicate `~/.codex/skills/<skill>` symlinks manually after confirming
they point back to this repo; the installer will not delete existing user files.

Skill creation follows a source-plus-surfaces model:

- Shared/global source of truth lives in this repo at `skills/<skill-name>/`, then `./install.sh`
  links provider surfaces such as `~/.agents/skills`, `~/.claude/skills`, and the configured
  Gemini/Antigravity skills directory.
- Project-local skills should still use one portable source plus provider surfaces. Prefer
  `<project>/.agents/skills/<skill-name>/` as the portable source and link
  `<project>/.claude/skills/<skill-name>` when Claude project discovery is needed.
- Do not default to `.codex/skills` as a duplicate mirror. Create it only when a project or
  provider requires that surface.
- When creating or updating a skill, report the source of truth, surfaces created or skipped,
  validation status, and any install/discovery gaps.

### 4) Routing model

Use this hierarchy when deciding what should activate:

- Rules: baseline behavior, language, safety, and output constraints in `AGENTS.md` or `CLAUDE.md`.
- Task skills: reusable workflows for the user's active task, either explicit with `$skill` or implicit through descriptions. Examples: `implement-change`, `diagnose`, `verification-driven-change`, `planning-grill`, `code-review`, `project-lifecycle`, and `conventional-git-flow`.
- Policy/reference skills: supporting constraints for deeper details. Examples: `policy-api`, `policy-security`, `policy-testing`, `policy-go`, `policy-rust`, `policy-frontend`, and `policy-infra`.
- Agents: bounded role isolation or parallel work when the split is concrete and worth the cost.
- Hooks: deterministic guardrails only; hooks do not replace sandboxing, permissions, rules, skills, or agents.

The default posture is semi-automatic routing: the main assistant should say which skill or agent it is using
when the match is clear, suggest heavier multi-agent routes when useful, and keep high-cost challenge passes
such as `oracle` explicit.

Do not rely on multiple policy skills being automatically loaded together for correctness. A task skill should
carry the minimum workflow guardrails needed for the task, while policy/reference skills add deeper detail when
their topic is primary or materially affects the task. `execution-harness` is orchestration-only: use it for
multi-phase, multi-agent, handoff-heavy, lifecycle-gated, or verification-gated work, not as a replacement for
implementation, diagnosis, design clarification, review, testing strategy, git workflow, or lifecycle capture
skills.

### 5) Multi-agent workflow

Codex runtime roles in this repo use one primary path:

- source of truth: `config/codex-agents/*.toml` templates
- installed location: generated `~/.codex/agents/*.toml`
- discovery model: Codex loads those standalone role files directly

`agents/*.md` remains in the repo as human-readable role templates and shared documentation.
It is not the runtime source for Codex custom roles.

Claude also includes `claude-agents/codex-collaboration-loop.md`, a workflow agent for cases
where Claude should coordinate with Codex across planning consensus, execution delivery, or
end-to-end delivery. It covers drafting a plan, Codex review, Claude triage, Codex implementation,
independent Codex review, accepted fixes, final adversarial review, and Claude's final acceptance
check. Use it when prompts include requirements such as "give this to Codex for review",
"split this for Codex to implement", "Codex should execute this", "run Codex adversarial review",
or "you and Codex should reach consensus before asking me".

Claude agent files in `claude-agents/*.md` are runtime definitions, not just human-readable role
notes. Each file must keep YAML frontmatter with `name` and `description` so Claude Code can list
and invoke it from `~/.claude/agents`.

For this standalone-file layout, you do not need matching `[agents.<name>]` entries in
`~/.codex/config.toml`. Keep `[agents]` for global agent runtime settings such as thread or
depth limits. Use `[agents.<name>]` only as an advanced option when you explicitly want a
config-declared role that points at a separate `config_file`.

Codex role files in `config/codex-agents/*.toml` can include:

- `model`: model id used by the role
- `model_reasoning_effort`: `low|medium|high|xhigh`
- `service_tier`: optional request tier such as `fast`
- `default_permissions`: optional permission profile such as `:read-only` for read-only roles

Important compatibility note:

- Agent role files are parsed as standalone config layers by current Codex builds.
- MCP server definitions and `[[skills.config]]` entries are generated from `config/codex-agents/role-manifest.json`; do not add them directly to role templates.
- Permission profiles do not compose with legacy `sandbox_mode`; do not add `sandbox_mode` to Codex role files when using `default_permissions`.

Default Codex profile in this repo:

- `orchestrator`: `gpt-5.5` + `high` + `fast`
- `explorer`: `gpt-5.5` + `medium` + `fast`
- `builder`: `gpt-5.3-codex` + `medium`
- `git-commit`: `gpt-5.3-codex-spark` + `low` + `workspace-git`
- `reviewer`: `gpt-5.4` + `high` + `fast` + `:read-only`
- `oracle`: `gpt-5.5` + `high` + `:read-only`
- `librarian`: `gpt-5.4` + `medium` + `fast`
- `test-runner`: `gpt-5.3-codex-spark` + `low`

Why this mix:

- `gpt-5.3-codex` stays on implementation-heavy work where prior results have been strong.
- `gpt-5.4` handles daily review where official Codex subagent examples favor a high-effort reviewer profile.
- `gpt-5.5` is reserved for high-precision coordination and manual challenge passes.
- `gpt-5.3-codex-spark` stays on narrow deterministic validation and simple git workflow execution where speed matters, with strict prompts to preserve commit quality.

For simple git actions, direct requests to create a branch, commit, push, or open a PR for
already-existing changes are independent explicit standing authorization to delegate execution to the
Spark-backed `git-commit` role. The main session should inspect enough context to produce a compact
execution handoff packet: repo path, action list, exact file inclusions or explicit all-files
approval, exclusions, commit subject/body decision, branch/remote/base/head values, PR title/body
when relevant, and no amend/rebase/force/no-verify/code-edit constraints. `conventional-git-flow`
owns the reusable git workflow policy; `git-commit.toml` is only the thin Codex execution role. With
a complete handoff, the role should execute directly without re-reading the skill or re-inspecting
diffs for policy decisions. If runtime policy blocks delegation or the handoff is incomplete, the
role fails closed and reports the blocker instead of guessing. On success, the delegate's final
status confirmation is authoritative enough to relay; the main session should close the completed
agent and avoid a second confirmation loop unless the delegate output is incomplete or
contradictory. Hook layers still protect the shell commands that role runs and may add reminders, but
hooks do not select models or delegate workflow.

Simple git delegation uses the `workspace-git` permission profile declared inside the `git-commit`
role file. The Go installer also keeps the same named profile in `~/.codex/config.toml` for manual
selection and compatibility: it backs up the config only when it changes the file, sets the old
repo-managed top-level `default_permissions = "workspace-git"` value back to the legal built-in
`:workspace` default, and upserts the block below. The setup step is best-effort inside
`install.sh`: legacy `sandbox_mode` conflicts are reported without blocking agent, skill, or shell
symlink setup. Run `scripts/setup_codex_config.sh` directly for a focused failure message after
resolving legacy sandbox settings. The profile extends `:workspace`, keeps normal workspace
protections, adds write access to `.git` metadata, and allows GitHub network domains used by
`git push` and `gh pr create`. The generated `~/.codex/agents/git-commit.toml` includes its own
copy of this profile block so commit delegation does not depend on the parent session's
`default_permissions`. `git-commit` does not use role-local `sandbox_mode = "danger-full-access"`.

```toml
default_permissions = ":workspace"

[permissions.workspace-git]
description = "Workspace editing with git metadata writes and GitHub network access."
extends = ":workspace"

[permissions.workspace-git.filesystem.":workspace_roots"]
"." = "write"
".git" = "write"
".codex" = "read"
".agents" = "read"
"**/.env" = "deny"
"**/.env.*" = "deny"
"**/*.env" = "deny"

[permissions.workspace-git.network]
enabled = true

[permissions.workspace-git.network.domains]
"github.com" = "allow"
"api.github.com" = "allow"
"uploads.github.com" = "allow"
```

Examples:

- Spark-backed simple action: "Create a branch.", "Commit this.", "Push this branch.", or
  "Open a PR."
- Main-session first, then Spark git workflow: "Fix the bug, then create a branch, commit, push, and
  open a PR."
- Main-session only until scope is fixed: "Amend the published PR", "resolve the merge conflict", or
  "commit whatever looks right".

Claude agent files use model aliases instead of pinned model IDs:

- `orchestrator`, `reviewer`, `oracle`, and `codex-collaboration-loop`: `opus` + `xhigh`
- `builder`: `sonnet` + `high`
- `explorer`, `librarian`, and `test-runner`: `haiku`

Aliases let Claude Code resolve the current provider default while keeping the intended role tier clear.
Claude read-only roles also define `tools` frontmatter so they do not inherit the full tool surface:

- `explorer`, `reviewer`, and `oracle`: `Read, Glob, Grep`
- `librarian`: `WebSearch, WebFetch, Read, Glob, Grep`

Skills are discovered automatically from repo, user, and system skill directories. For this repo's
managed Codex role files, per-agent MCP and skill enablement is generated from
`config/codex-agents/role-manifest.json`.

Provider-specific role prompts should not duplicate full shared skill policy. For example,
`conventional-git-flow` owns branch, commit, push, and PR workflow rules, while the Codex
`git-commit` role only defines the fast execution model, non-interactive boundaries, and minimal
fallback when the skill is unavailable. Git metadata writes are handled by the parent session's
`workspace-git` permission profile, not by role-local sandbox overrides.

`execution-harness` is an optional process skill for structured coordination across phases, agents, git/workspace state, verification gates, diff sanity, lifecycle gates, and capture candidates. It is orchestrator-suggested and user-approved; it is not a mandatory SDLC and does not replace domain policy skills.

Current role-scoped tool strategy:

- `orchestrator`: coordination MCPs only (`sequential-thinking`, `clickup`), broad skill surface
- `explorer`: local code discovery only, no external MCPs, no creator/installer/design skills, and guarded from self-applying harness or lifecycle workflow
- `builder`: local implementation only, no external MCPs, no system installer skills, and guarded from self-applying harness or lifecycle workflow
- `reviewer`: local review only, no external MCPs, no creator/installer/design skills, and guarded from self-applying harness or lifecycle workflow
- `oracle`: manual high-precision read-only challenge pass, no external MCPs, and guarded from self-applying harness or lifecycle workflow
- `librarian`: web research by default, optional `context7` MCP only when explicitly enabled, with most repo-specific implementation skills disabled
- `test-runner`: no external MCPs, focused on validation with a reduced skill surface

This mirrors the same general principle used by `oh-my-opencode-slim`: keep coordination tools on the coordinator, external knowledge tools on the research role, and keep the code-heavy roles as narrow as possible. This repo deliberately keeps the heavier council-style path manual through `oracle` instead of adding always-on multi-model consensus.

The harness model follows the same principle: the orchestrator may propose `execution-harness` when a task needs explicit phase, agent, git/workspace, verification, lifecycle, or capture coordination, but specialist roles stay narrow and do not dynamically gain skills or MCPs at runtime.

Where Codex per-agent skill configuration is available, specialist roles disable `execution-harness`
and `project-lifecycle` across the standard installed and repo-local skill paths. Other providers may
still see shared skills from their global skill directories, so their role prompts carry explicit
guardrails against self-applying harness or lifecycle workflow.

Antigravity CLI is optional. The `antigravity-design-bridge` skill can ask Antigravity CLI for bounded UI or visual design critique, and may allow low-risk direct edits when target files and scope are explicit. The installer does not install Antigravity CLI, and normal workflows do not depend on it.

Context7 is also optional. Codex role MCP definitions are generated from
`config/codex-agents/role-manifest.json`; a generated MCP is enabled only when the role allows it,
`~/.codex/config.toml` already has a matching `[mcp_servers.<name>]` table, and the configured local
command, such as `bunx`, is present on `PATH`. The installer treats the base config table as the
local opt-in signal; it does not run MCP commands, probe package caches, download packages, or check
remote services.

Codex role MCP definitions and skill toggles are generated from
`config/codex-agents/role-manifest.json`. Edit the manifest when a role needs an MCP allowed or a
skill disabled; do not add `[mcp_servers.*]` or `[[skills.config]]` blocks to the role templates.
Some generated skill entries intentionally cover multiple possible skill roots such as
`{{CODEX_HOME}}/skills`, `{{AGENTS_HOME}}/skills`, and this repo's local `skills/` path. The
`{{CODEX_HOME}}/skills` entries are retained as legacy/duplicate-root safeguards, not as the
preferred shared skill installation path.

The checked-in Codex role files are templates. During `setup-codex-agents`, the Go installer renders
`{{HOME}}`, `{{CODEX_HOME}}`, `{{AGENTS_HOME}}`, and `{{REPO_ROOT}}` into managed
`~/.codex/agents/*.toml` files. The installer overwrites only its own managed generated files or
legacy symlinks that point back to this repo's `config/codex-agents`; user-owned role files are
skipped.

The role prompts are intentionally Codex-first rather than a direct copy of `oh-my-opencode-slim`. In practice that means shorter instructions, less handoff ceremony, and stronger emphasis on bounded delegation.

Recommended task flow for better results than a single long session:

1. `orchestrator` defines done criteria, task split, and validation plan.
2. For complex work, `orchestrator` may suggest `execution-harness` before continuing.
3. `explorer` maps the exact edit surface and regression risks.
4. `builder` implements one bounded slice at a time.
5. `reviewer` checks changed behavior and nearby contracts for regressions.
6. `oracle` runs only when the user explicitly asks for adversarial review or a high-risk challenge pass.
7. `test-runner` runs only the smallest useful checks.
8. `librarian` is used only when external docs, release notes, or version-sensitive facts matter.

This setup works best when the orchestrator keeps tasks narrow and routes work through explicit handoffs instead of letting every role inspect everything.

Use `project-lifecycle` as the canonical lifecycle/capture gate when work reaches a project
decision point, implementation pivot, phase boundary, handoff, status/doc sync point, or
Hermes-style workflow learning review. It may proactively surface capture candidates from
`planning-grill`, `implement-change`, or `execution-harness`, but it must ask before mutating
long-lived docs, shared skills, agent files, shared config, git state, or external systems.
Shared skill changes still require explicit approval and `skill-creator`.

The former standalone `phase-closeout`, `sync-decision-docs`, and `session-retrospective` skills
were removed after their lifecycle/capture responsibilities were merged into `project-lifecycle`.

### 6) Hook plugins

Hooks in this repo are advisory guardrails for repeatable mechanical checks. They do not replace
provider sandboxing, permission prompts, `AGENTS.md`, `CLAUDE.md`, skills, or agent role prompts.

Current plugins:

- `ai-agent-config-guardrails`: Go-backed `PreToolUse` guardrails that block direct file-tool writes
  to `.env*`, block direct file-tool writes to global agent configuration such as
  `~/.codex/config.toml`, `~/.codex/hooks.json`, and `~/.claude/settings.json`.
- `ai-agent-git-routing`: Go-backed `UserPromptSubmit` reminder for direct simple git action prompts
  so the main session sees that branch creation, commit, push, or PR creation should route through
  the `git-commit` role.

Codex uses repo-local plugins instead of owning `~/.codex/hooks.json`. The configured events include
file-tool `PreToolUse` guardrails and a lightweight `UserPromptSubmit` reminder. Open `/hooks` in
Codex after installation to review and trust the plugin hooks. Hook commands call plugin-bundled Go
binaries through the Codex plugin runtime root (`${PLUGIN_ROOT}`) and do not require a
`~/.codex/hooks` symlink.

Static context belongs in `AGENTS.md`, `CLAUDE.md`, or skills. The prompt hook is reinforcement only:
it appends a visible reminder for matching git action prompts, while `AGENTS.md` and
`conventional-git-flow` own the fail-closed routing rule. This repo does not register a
`SessionStart` hook unless future dynamic session context is worth the startup/resume cost.

Claude hook behavior is owned by the active `rtk hook claude` entry in `~/.claude/settings.json`.
This repo no longer installs Claude hook script symlinks.

Build hook binaries before installing or reinstalling the plugins:

```bash
scripts/build_go_hooks.sh
codex plugin add ai-agent-config-guardrails@ai-agent-config
codex plugin add ai-agent-git-routing@ai-agent-config
```

Official compatibility boundaries:

- This repo relies only on Codex `command` hooks and keeps `PermissionRequest` and `PostToolUse`
  unregistered; provider permission prompts own approval decisions, and skills/final responses own
  workflow reminders.
- The Codex config intentionally matches only the local shell and file-edit tool names this repo
  depends on (`Bash`, `apply_patch`, `Edit`, `Write`, and `MultiEdit`). MCP or future tool-specific
  hooks should be added only when there is a concrete guardrail worth enforcing.
- Codex file-tool `PreToolUse` is a guardrail, not a complete enforcement boundary, because
  equivalent work may be possible through another supported tool path or through shell commands that
  obscure touched files.
- Claude hook support is broader, but this repo keeps v0 aligned to the shared command-hook subset to
  avoid provider-specific drift.

Manual hook verification:

```bash
scripts/build_go_hooks.sh
PLUGIN_ROOT="$PWD/plugins/ai-agent-config-guardrails"

"$PLUGIN_ROOT/hooks/bin/agent-config-guardrails" pre-tool-use-file <<'JSON'
{"hook_event_name":"PreToolUse","tool_name":"Write","tool_input":{"file_path":"notes.txt"}}
JSON

"$PLUGIN_ROOT/hooks/bin/agent-config-guardrails" pre-tool-use-file <<'JSON'
{"hook_event_name":"PreToolUse","tool_name":"Write","tool_input":{"file_path":".env.local"}}
JSON

PLUGIN_ROOT="$PWD/plugins/ai-agent-git-routing"
"$PLUGIN_ROOT/hooks/bin/agent-git-routing" <<'JSON'
{"hook_event_name":"UserPromptSubmit","prompt":"幫我創建分支"}
JSON

"$PLUGIN_ROOT/hooks/bin/agent-git-routing" <<'JSON'
{"hook_event_name":"UserPromptSubmit","prompt":"先幫我 commit"}
JSON

"$PLUGIN_ROOT/hooks/bin/agent-git-routing" <<'JSON'
{"hook_event_name":"UserPromptSubmit","prompt":"幫我 push"}
JSON

"$PLUGIN_ROOT/hooks/bin/agent-git-routing" <<'JSON'
{"hook_event_name":"UserPromptSubmit","prompt":"幫我發 PR"}
JSON

"$PLUGIN_ROOT/hooks/bin/agent-git-routing" <<'JSON'
{"hook_event_name":"UserPromptSubmit","prompt":"幫我想 commit message"}
JSON
```

The first command should print nothing and exit `0`. The second should print a deny JSON object for
the `.env*` file-tool write. The next four git-routing commands should each print a `systemMessage`
reminder for the `git-commit` route. The final text-only command should print nothing.

## Update Workflow

```bash
cd ~/git/ai-config
git pull
./install.sh
```

Because symlinks point to this repo, content updates are picked up immediately.

Optional local shell maintenance:

```bash
auto-clean --agent-tools
```

The `--agent-tools` flag updates Bun global packages, updates Cargo-installed tools through
`cargo-update`, and checks RTK Codex status. Install the Cargo updater once with:

```bash
cargo install cargo-update
```

`auto-clean --agent-tools` does not rerun `rtk init -g --codex`; refresh RTK instructions manually
when needed.

## Add New Skill

```bash
cd ~/git/ai-config/skills
mkdir -p policy-example/references
cat > policy-example/SKILL.md <<'DOC'
---
name: policy-example
description: Apply example policy guidance. Use when the user asks for example policy behavior. Avoid when a more specific policy skill applies.
metadata:
  version: "0.1.0"
---

# Example Skill

## Purpose

Apply concise, portable guidance for the example policy.

## Use When

- The request needs example policy behavior.

## Avoid When

- A more specific policy skill applies.

## Version History

- v0.1.0 (YYYY-MM-DD): Initial portable skill draft.

## References

- `references/INDEX.md` - Use when deeper topic navigation is needed.
DOC

cat > policy-example/references/INDEX.md <<'DOC'
# References

Add focused reference files here only when details would make SKILL.md too long or too hard to route.
DOC
```

## Add New Agent

```bash
cd ~/git/ai-config/config/codex-agents
cp reviewer.toml security-auditor.toml
```

Then update `security-auditor.toml` with a unique `name`, `description`, model settings, and any
role-specific instructions. Add MCP allowances or disabled skill groups in
`config/codex-agents/role-manifest.json`.

If you also want a matching human-readable role brief for other tools or repo documentation, add a
separate `agents/security-auditor.md`, but treat that as optional documentation rather than the
Codex runtime source.

## License

MIT
