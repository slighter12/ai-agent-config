# AI Agent Configuration

A shared, versioned configuration for Codex, Claude Code, OpenCode, and Gemini/Antigravity. The repository uses task-shaped skills with explicit invocation semantics instead of a policy dependency graph.

## Architecture

- `AGENTS.md`: universal behavior, authorization, security, and verification boundaries.
- `CLAUDE.md`: imports the shared `AGENTS.md` instructions.
- `skills/`: one portable source for every shared skill.
- `config/codex-agents/`: Codex role templates and role-scoped tool configuration.
- `agents/` and `claude-agents/`: human-readable and Claude runtime role definitions.
- `hooks-go/`: installer, skill tooling, and deterministic hook binaries.
- `plugins/`: repo-local Codex plugins.

Skills teach workflows or durable judgment. They do not restate abilities already supplied by the model or runtime. Language, framework, SDK, protocol, and infrastructure details should come from repository standards and current official documentation.

## Install

Go is required. The shell scripts are compatibility wrappers around `hooks-go/cmd/agent-config`.

```bash
./install.sh
```

The installer:

- links `AGENTS.md` or `CLAUDE.md` into provider homes that already exist;
- links the same `skills/*` directories to Codex, Claude, OpenCode, and Gemini/Antigravity discovery surfaces;
- generates managed Codex role files;
- relies on each skill's frontmatter and provider sidecar for declarative invocation controls;
- removes only retired skill symlinks that point exactly into this repository, while leaving user files, unknown targets, and links from other repositories unchanged;
- builds and installs this repository's Codex guardrail plugin.

Default skill surfaces:

- Codex: `~/.agents/skills`
- Claude: `~/.claude/skills`
- OpenCode: `~/.config/opencode/skills`
- Gemini/Antigravity: `~/.gemini/antigravity-cli/skills`

Environment variables such as `CODEX_HOME`, `CLAUDE_HOME`, `OPENCODE_HOME`, `ANTIGRAVITY_HOME`, and `GEMINI_HOME` override the defaults.

## Skill Catalog

The canonical catalog and route map live in [`ask-matt/references/CATALOG.md`](skills/ask-matt/references/CATALOG.md). It contains Matt Pocock's 25 promoted skills plus two local extensions, `design-art-direction` and `mcp-builder-go`. User-invoked means the skill instructions load only after explicit invocation; an equivalent plain-language request may still use the model's default behavior under `AGENTS.md` and provider permissions.

Upstream runtime artifacts are synchronized from the pinned source recorded in [`skills/UPSTREAM.md`](skills/UPSTREAM.md). The source stays flat locally because the installer links each skill directly into four provider discovery roots.

## Project Coordination

`setup-matt-pocock-skills` initializes these paths only when requested:

- `docs/agents/issue-tracker.md`: external tracker configuration or local Markdown conventions, including Wayfinder operations.
- `docs/agents/triage-labels.md`: the five canonical triage roles.
- `docs/agents/domain.md`: how skills find root or context-local `CONTEXT.md` files and ADRs.

Artifact-producing workflows read these files instead of guessing. `CONTEXT.md`, `CONTEXT-MAP.md`, and `docs/adr/` are created lazily when durable domain information exists. Explicitly invoke `setup-matt-pocock-skills` to initialize or change this configuration.

## Provider Invocation Contract

A user-invoked skill has:

```yaml
disable-model-invocation: true
metadata:
  invocation: user
  opencode/autoinvoke: "false"
```

and `agents/openai.yaml`:

```yaml
policy:
  allow_implicit_invocation: false
```

`disable-model-invocation` is Claude's file-local control, `opencode/autoinvoke` is OpenCode's metadata control, and the sidecar is Codex's control. The shared source has been parsed by all three local clients. Gemini/Antigravity receives the shared skill link, but runtime invocation behavior is best-effort and should be reported as unverified unless probed.

## Skill CLI

```bash
go -C hooks-go run ./cmd/agent-config --repo-root .. init-skill <skill-name> --path <absolute-path>
go -C hooks-go run ./cmd/agent-config --repo-root .. validate-skill <absolute-skill-directory>
go -C hooks-go run ./cmd/agent-config --repo-root .. validate-skills ../skills
go -C hooks-go run ./cmd/agent-config --repo-root .. package-skill <absolute-skill-directory> [<absolute-output-directory>]
```

Run these commands from the repository root. Path arguments are resolved from `hooks-go`, so use absolute paths except for the canonical repository catalog shown above. `init-skill` creates only a minimal `SKILL.md` with an explicit invocation choice placeholder. It does not precreate references or impose version-history and output sections.

`validate-skills` checks:

- the supplied catalog resolves to `<repo-root>/skills`, so manifest and routing evidence cannot be mixed across repositories;
- supported-provider skill validity;
- unique names and directory/name agreement;
- `user` or `model` invocation metadata;
- Claude, Codex, and OpenCode controls for user-invoked skills;
- the required `ask-matt` router and synchronized canonical catalog;
- an 8,000-character repo-local contribution limit for model skill names and descriptions. Codex applies its [runtime budget](https://github.com/openai/codex/blob/main/codex-rs/core-skills/src/render.rs#L143-L158) to the complete rendered metadata from every source, so this is not a full-session estimate.

## Codex Roles

Codex role source lives in `config/codex-agents/*.toml`; installation generates managed files under `~/.codex/agents/`. `config/codex-agents/role-manifest.json` owns role-scoped MCP availability and intentional skill exclusions.

Use provider-native delegation when a subtask is independent, concrete, and its isolation or parallelism justifies the coordination cost; otherwise complete the equivalent work sequentially in the current session. `research` is the upstream-first exception: run its reading legwork in a background agent whenever the provider primitive is available, then complete the same pass in the current session only when that primitive is unavailable or fails to start. `code-review` keeps independent Standards and Spec passes by default; explicit correctness, full-review, and release-readiness requests add their documented local axes. High-cost adversarial review remains explicit.

## Hooks and Codex Profiles

The Go installer maintains:

- Codex base permission defaults and an explicit workspace-git profile;
- generated standalone Codex roles;
- the zsh Codex profile helper;
- the `ai-agent-config-guardrails` plugin.

Hooks are deterministic guardrails; they do not replace skills, sandbox permissions, or user approval.

RTK remains optional and integrates through live hooks rather than a `CLAUDE.md` import:

```bash
brew install rtk
rtk init -g --codex
```

Repository-managed baseline files remain the source of truth. If an RTK command replaces a managed link, rerun `./install.sh`.

## Verify

Validate the source catalog before installation:

```bash
go -C hooks-go run ./cmd/agent-config validate-skills ../skills
```

`validate-skill` and `package-skill` enforce the complete invocation contract for one skill. `validate-skills` also checks that the canonical `ask-matt` catalog contains every skill exactly once in the correct invocation group.

Then inspect provider surfaces:

```bash
ls -la ~/.agents/skills
ls -la ~/.claude/skills ~/.claude/CLAUDE.md
ls -la ~/.config/opencode/skills ~/.config/opencode/AGENTS.md
ls -la ~/.gemini/antigravity-cli/skills ~/.gemini/GEMINI.md
```

For installation tests, use temporary provider homes. Do not point development tests at real global homes.

## Updating

```bash
git pull
./install.sh
```

Run `go -C hooks-go run ./cmd/agent-config validate-skills ../skills` after catalog changes. Compare against the SHA and path mapping in `skills/UPSTREAM.md` before adopting a new upstream release. Add references only when a branch of the workflow genuinely needs details that would otherwise obscure the main skill.

Run `./scripts/check_go_targets.sh` after changing `hooks-go` code guarded by build tags. `go build` and `go test` on the host cover one target only, so a broken build tag or an unresolved `//go:linkname` on another platform stays invisible until that platform is compiled.

## License

Personal configuration. Adapt as needed. Vendored Matt Pocock skill content remains subject to the MIT license in `LICENSES/mattpocock-skills-MIT.txt`.
