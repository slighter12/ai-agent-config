# Authoring Standards

Use these standards to write portable, high-signal skills for Codex, Claude, Gemini, and OpenCode.

## Source Priority

1. Provider official docs for supported syntax and invocation behavior.
2. Local repo conventions for language, safety, validation, and output style.
3. Marketplace or community examples only for trigger phrasing and structure ideas.

## Portable Core

- Use `SKILL.md` with YAML frontmatter containing `name`, `description`, and optional standard fields: `license`, `compatibility`, and `metadata`.
- Make `name` a stable hyphen-case identifier.
- Make `description` a routing contract: capability first, then `Use when...`, then `Avoid when...`.
- Use `metadata.version` for machine-readable versioning and `Version History` for human-readable changes.
- Keep `SKILL.md` under 500 lines; move detailed material into `references/`.
- Use `scripts/` only for deterministic helpers that are better executed than regenerated.
- Use `assets/` only for reusable files consumed by outputs.
- Do not add README, changelog, install guide, or broad examples unless the user explicitly needs them.

## Placement Lifecycle

Before creating files, decide whether the skill is shared/global, project-local, or a provider-specific exception. Use `PLACEMENT_AND_INSTALLATION.md` for source-of-truth and provider surface rules.

- Shared/global skills in this repo live under `skills/<skill-name>/` and are installed with `./install.sh`.
- Project-local skills should still have a canonical portable source and provider surfaces; prefer `<project>/.agents/skills/<skill-name>/` as source, but do not treat that source path as proof of Codex project discovery.
- Symlink `<project>/.claude/skills/<skill-name>` when Claude project discovery is needed.
- Do not default to `.codex/skills` as a duplicate mirror. Use it only when a project or provider explicitly requires or verifies that surface.
- A create/update result must say which surfaces were created, linked, already present, skipped, not configured, or unverified.

## Guardrails Versus Policy

Duplicate guardrails, not policy.

- Put the minimum task-critical guardrails directly in a task skill when missing them would cause drift.
- Keep detailed policy, matrices, examples, edge cases, and language/API/security rules in references or narrower policy skills.
- Do not rely on a task skill automatically loading another skill for correctness.
- Do not copy full policy documents into multiple task skills; repeat only short, stable constraints such as "establish a feedback loop before fixing" or "keep edits minimal and aligned with repo patterns."
- When the same guardrail appears in multiple skills, keep it short and generic so future policy changes do not require many synchronized edits.

## Prompt Shape

- Start with the outcome the skill enables.
- Define trigger cases and nearby non-goals before workflow steps.
- Keep workflow steps stable and minimal; do not force step-by-step reasoning when order does not affect correctness.
- Include explicit stop conditions for search, validation, and side effects.
- Define output shape only as strictly as the task requires.
- Use `MUST`, `CRITICAL`, `always`, and `never` only for security, data protection, destructive actions, or hard compatibility rules.

## Predictability And Pruning

- Prefer a model-invoked skill only when the agent must reach it automatically or another skill must route to it; otherwise keep the workflow behind explicit invocation or an existing owner to avoid context load.
- Split a workflow only when it has a distinct trigger, a distinct owner, or later steps cause premature completion of earlier legwork.
- Give important steps checkable completion criteria. A vague "review thoroughly" is weaker than a concrete stop condition.
- Use a stable leading word for repeated concepts when it makes routing or execution sharper.
- Keep one source of truth for each rule. Delete duplicate wording, no-op prose, stale sediment, and sprawl before adding new references.

## Provider Notes

- Codex and Gemini rely on portable `name` and `description` for skill routing; keep optional standard fields small and spec-aligned.
- Claude supports extra frontmatter such as `allowed-tools`, `disable-model-invocation`, and `user-invocable`; add these only in Claude-specific overlays or copies.
- OpenCode discovers standard skills from its own and compatible `.agents` surfaces, lists their names and descriptions, and loads full skill bodies on demand through the native skill tool.
- For one shared skill file, express Claude-like safety behavior in normal Markdown sections instead of provider-specific frontmatter.

## Supported Frontmatter Subset

The local validator intentionally supports a small portable subset:

- Top-level scalar fields: `name`, `description`, `license`, and string-form `compatibility`.
- Top-level inline list: `compatibility: [codex, claude, gemini, opencode]`.
- One nested mapping: `metadata`, with single-line scalar values such as `version` and `author`.
- No block lists, multiline strings (`|` or `>`), tab indentation, anchors, or complex YAML objects.

Use this subset for shared skills. Put richer provider-specific metadata in overlays or packaging manifests instead of `SKILL.md` frontmatter.

## Template Sections

Use this section order unless the skill has a strong reason to differ:

1. `Purpose`
2. `Use When`
3. `Avoid When`
4. `Workflow`
5. `Tool And Side-Effect Boundaries`
6. `Output`
7. `Version History`
8. `References`

## Pre-Release Checklist

- Frontmatter has `name`, `description`, and optional standard fields only.
- Description has capability, trigger cases, and avoid cases.
- `SKILL.md` is under 500 lines.
- Workflow steps have checkable completion criteria where correctness depends on stopping at the right point.
- No-op prose, duplicate rules, and provider-specific frontmatter have been pruned from shared skill files.
- Version history is present for shared or team-maintained skills.
- Detailed docs are linked from `references/` instead of duplicated.
- Scripts and assets are present only when they materially improve reliability.
- Provider-specific behavior is documented as an overlay, not mixed into shared frontmatter.
- Source placement and provider surfaces are reported, including any install/discovery gap.
