---
name: skill-creator
description: Create or update portable CLI skills, including source placement, provider surfaces, validation, packaging, and install/discovery status. Use when the user asks to create, revise, validate, package, modernize, or decide placement for a SKILL.md-based skill for Codex, Claude, Gemini, project-local, or shared multi-provider use, or when `project-lifecycle` has an approved shared skill capture candidate. Avoid when the task is only application code changes, ordinary documentation edits, lifecycle capture triage, or third-party skill installation without authoring changes.
metadata:
  version: "0.2.5"
---

# skill-creator

## Purpose

Create high-signal, portable skills that work well across Codex, Claude, and Gemini by default. Treat skill creation as a lifecycle: choose the source of truth, create or update the portable `SKILL.md`, sync provider discovery surfaces when appropriate, validate, and report install/discovery status.

## Use When

- The user wants to create a new reusable skill.
- The user wants to update an existing skill for cross-provider routing.
- The task involves `SKILL.md`, skill descriptions, bundled references, scripts, assets, packaging, or validation.
- The user asks where a skill should live, how to make it project-local or shared, or why a provider did not discover a skill.
- The user asks whether a workflow belongs in a skill, an agent prompt, or a provider-specific overlay.
- `project-lifecycle` has classified a workflow or skill lesson as an approved `shared_skill_update` or `new_shared_skill`.

## Avoid When

- The request is only a normal code change in an application repo.
- A more specific domain skill should guide the work after the skill structure is already clear.
- The request is to decide whether a session lesson, project decision, handoff, or workflow observation is worth capturing; use `project-lifecycle`.
- The user is asking to install third-party skills rather than author or modify them.

## Cross-Provider Defaults

Use the portable core format unless the user explicitly asks for a provider-specific variant:

- Frontmatter contains `name`, `description`, and optional standard fields: `license`, `compatibility`, and `metadata`.
- `description` states capability, trigger conditions, and nearby avoid cases.
- `SKILL.md` stays short enough to route and load quickly; target fewer than 500 lines.
- `metadata.version` is the recommended machine-readable version for installers, listings, and reload/cache visibility.
- `Version History` is the human-readable change log for meaningful behavior changes.
- Deep guidance lives in `references/`; deterministic helpers live in `scripts/`; reusable output assets live in `assets/`.
- Task skills carry the minimum guardrails needed to avoid drift, but detailed policy stays in references or narrower policy skills.
- Tool, write, deployment, credential, notification, or destructive side effects require explicit boundaries.
- Creation is not complete until the source path, provider surfaces, validation status, and install/discovery gaps are reported.

For CLI-backed workflows, use this skill to decide whether a workflow belongs in a skill, reference, script, or helper. Use `mcp-builder-go` when the task is to design or implement the actual MCP/CLI tool surface and output contract.

## Claude-Only Overlay Rules

Do not put Claude-only keys in the shared template. Add them only in a Claude-specific copy or installation overlay:

- Use `disable-model-invocation: true` for deployment, commit, external-service, destructive, or high-side-effect skills.
- Use `allowed-tools` when a skill should be constrained to a small tool surface.
- Use `user-invocable: false` when a skill should be background-only and not exposed as a slash command.

If a user asks for one file that must run everywhere, keep provider-specific controls out of frontmatter and describe the behavior in `Tool And Side-Effect Boundaries` instead.

## Workflow

1. Identify the intended outcome, target users, provider targets, trigger examples, avoid examples, side effects, and output contract.
2. For approved skill evolution requests, use `references/SKILL_EVOLUTION_GATE.md` to inventory visible skills and prefer the narrowest existing owner before proposing a new skill.
3. Decide placement using `references/PLACEMENT_AND_INSTALLATION.md`: shared/global, project-local, or provider-specific exception.
4. Choose the shared portable template by default. Use provider-specific variants only when a provider feature is necessary for correctness or safety.
5. Decide what minimum guardrails must live in the task skill itself so it remains safe when supporting policy skills are not loaded.
6. Draft or update `SKILL.md` with `Purpose`, `Use When`, `Avoid When`, `Workflow`, `Tool And Side-Effect Boundaries`, `Output`, and `References`.
7. Move long explanations, detailed policy, API notes, examples, and troubleshooting into focused files under `references/`.
8. Add scripts only when they improve deterministic reliability or prevent repeated code generation.
9. Create or report provider discovery surfaces:
   - shared repo skills use this repo's installer to link supported global surfaces;
   - project-local skills use one canonical project source plus provider surfaces; `.agents/skills/<name>` is the preferred portable source convention, not proof of Codex project discovery;
   - link `.claude/skills/<name>` when Claude project discovery is needed, and report Codex project discovery as verified, unverified, not configured, or provider-specific;
   - `.codex/skills` is not the default shared or project-local mirror for this repo; create it only when a project/provider explicitly requires or verifies that surface.
10. Run or recommend `agent-config validate-skill` when validation is requested or required by the packaging workflow.
11. Return a concise summary of changed skill behavior, placement, provider-specific notes, validation, and manual verification checklist.

## Output Contract

Return:

- `summary`: what skill behavior changed or was created.
- `files_touched`: exact paths.
- `placement`: `source_of_truth` and `scope`.
- `provider_surfaces`: paths created, linked, already present, skipped, not configured, or unverified.
- `validation`: checks run or why they were skipped.
- `install_status`: whether shared install, project-local symlinks, provider packaging, or discovery checks were performed or skipped.
- `provider_notes`: Codex, Claude, and Gemini compatibility implications.
- `assumptions`: correctness-relevant assumptions only.
- `manual_verification`: focused checklist for routing, validation, and packaging.

## Version History

- v0.2.0 (2026-05-04): Switch to portable core authoring with provider-specific overlay guidance.
- v0.2.1 (2026-05-04): Keep machine-readable versions in `metadata.version`.
- v0.2.2 (2026-05-13): Add conservative skill evolution and CLI-backed skill reference routing.
- v0.2.3 (2026-05-18): Add guardrail duplication guidance for self-contained task skills.
- v0.2.4 (2026-05-21): Add placement, provider surface, and install/discovery lifecycle guidance.
- v0.2.5 (2026-05-29): Restore self-contained skill evolution gate guidance without sibling skill reference dependencies.

## References

- `references/AUTHORING_STANDARDS.md` - Portable authoring rules and provider-specific overlay guidance.
- `references/SKILL_TEMPLATE.md` - Shared template for new skills.
- `references/workflows.md` - Workflow structuring patterns.
- `references/output-patterns.md` - Output contract patterns.
- `references/SKILL_EVOLUTION_GATE.md` - Self-contained gate for approved skill evolution and new skill admission.
- `references/CLI_BACKED_SKILLS.md` - When deterministic helpers, scripts, or CLIs belong behind a skill.
- `references/PLACEMENT_AND_INSTALLATION.md` - Source-of-truth, provider surface, install, and discovery rules.
