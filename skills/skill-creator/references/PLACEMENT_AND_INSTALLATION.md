# Placement And Installation

Use this reference when creating or updating a skill so authoring, install, and provider discovery stay aligned.

## Lifecycle

Skill creation is not complete when `SKILL.md` exists. Complete the lifecycle:

1. Decide scope: shared/global, project-local, or provider-specific exception.
2. Create or update the canonical source of truth.
3. Create, link, or report provider discovery surfaces.
4. Validate the portable skill.
5. Report source, surfaces, install/discovery status, and any skipped surfaces.

## Shared Or Global Skills

For this repo, shared personal skills use:

- Source of truth: `ai-agent-config/skills/<skill-name>/`.
- Codex shared surface: `~/.agents/skills/<skill-name>`.
- Claude shared surface: `~/.claude/skills/<skill-name>`.
- Gemini/Antigravity shared surface: configured Gemini/Antigravity skills directory, currently reached through `~/.gemini/antigravity-cli/skills`.

Use this repo's `./install.sh` to create or refresh global symlinks. Do not treat `~/.codex/skills` as the default shared source or mirror for this repo; it may contain system, plugin, legacy, or personal skills.

## Project-Local Skills

Use the same lifecycle as shared skills, but keep the source inside the target project:

- Preferred portable source: `<project>/.agents/skills/<skill-name>/`. This is a source convention, not proof that every provider will discover the skill automatically.
- Claude project surface: `<project>/.claude/skills/<skill-name>` symlinked to the portable source when project-local Claude discovery is needed.
- Codex project surface: report Codex project discovery as verified, unverified, not configured, or provider-specific. Do not claim the skill is Codex-discoverable solely because `.agents/skills` exists. Create `<project>/.codex/skills/<skill-name>` only when the project or provider explicitly requires or verifies that surface.
- Gemini/Antigravity project surface: use the project or extension packaging convention when one exists; otherwise report that Gemini discovery is not configured instead of inventing a hidden copy.

Use symlinks when possible so there is one portable source. Copy only when a provider requires provider-specific frontmatter or packaging.

## Provider-Specific Exceptions

Create provider-specific copies only when a provider feature is required for correctness or safety:

- Claude-only frontmatter such as `allowed-tools`, `disable-model-invocation`, or `user-invocable`.
- Gemini extension packaging rather than a raw portable skill directory.
- Codex-specific plugin or legacy skill location required by a project.

When making a provider-specific copy, state which file is canonical and how drift should be avoided.

## Reporting Requirements

Every create/update result should include:

- `placement`: `source_of_truth` and `scope` (`shared/global`, `project-local`, or `provider-specific exception`).
- `provider_surfaces`: paths created, linked, already present, skipped, not configured, or unverified.
- `validation`: validator and manual checks run or skipped.
- `install_status`: whether `./install.sh`, project-local symlinks, provider packaging, or discovery checks were performed or skipped.

If surfaces were not created, say the skill is authored but not yet discoverable for that provider.
