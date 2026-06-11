# skill-creator

Create and modernize portable `SKILL.md`-based skills for Codex, Claude, Gemini, or shared multi-provider use.

## What It Does

`skill-creator` helps design skills around a shared portable core:

- `SKILL.md` frontmatter uses `name`, `description`, and optional standard fields: `license`, `compatibility`, and `metadata`.
- `description` acts as the routing contract with capability, `Use when`, and `Avoid when`.
- The main skill body stays concise and outcome-first.
- `metadata.version` provides a machine-readable version; `Version History` documents human-readable changes.
- Detailed guidance moves to `references/`.
- Provider-specific controls, especially Claude-only invocation controls, stay in overlays or provider-specific copies.
- Creation includes placement, provider discovery surfaces, validation, and install/discovery status.

## When To Use

Use this skill when you need to:

- Create a new reusable skill.
- Revise an existing skill for better routing.
- Validate shared skill structure and description quality.
- Decide whether a skill should be shared/global, project-local, or provider-specific.
- Decide whether provider-specific behavior belongs in a shared skill or an overlay.
- Package a skill after validation.

Avoid this skill for ordinary application code edits or domain work that should be handled by a more specific policy or implementation skill.

## Portable Template

New skills should follow this shape:

```md
---
name: example-skill
description: Perform [capability]. Use when the user asks for [specific trigger contexts]. Avoid when [nearby non-goals].
metadata:
  version: "0.1.0"
---

# Example Skill

## Purpose

## Use When

## Avoid When

## Workflow

## Tool And Side-Effect Boundaries

## Output

## Version History

## References
```

See `references/SKILL_TEMPLATE.md` for the full template.

## Scripts

Create a shared skill source in this repo:

```bash
skills/skill-creator/scripts/init_skill.py my-new-skill --path skills
```

Create a project-local portable source:

```bash
skills/skill-creator/scripts/init_skill.py my-new-skill --path .agents/skills
```

Validate a skill:

```bash
skills/skill-creator/scripts/quick_validate.py skills/my-new-skill
```

Package a skill:

```bash
skills/skill-creator/scripts/package_skill.py skills/my-new-skill ./dist
```

## Provider Strategy

- Codex: keep the shared `name` and `description` concise because routing starts there. For this repo, shared Codex skills are installed through `~/.agents/skills`; do not mirror to `.codex/skills` unless a project proves it needs that surface. For project-local skills, treat `.agents/skills` as the portable source convention and explicitly verify or report Codex discovery.
- Claude: use overlays or provider-specific copies for `allowed-tools`, `disable-model-invocation`, and `user-invocable`. Project-local Claude skills should expose `.claude/skills/<skill-name>`, usually as a symlink to the portable source.
- Gemini: keep the shared skill compatible with Agent Skills structure and avoid provider-specific shared frontmatter. Use extension packaging when Gemini requires package-level metadata.

## Placement Lifecycle

Skill creation is complete only after the result reports:

- `placement`: `source_of_truth` and `scope`;
- `provider_surfaces`: paths created, linked, already present, skipped, not configured, or unverified;
- `validation`: validator and manual checks run or skipped;
- `install_status`: install, symlink, packaging, and discovery-check status.

See `references/PLACEMENT_AND_INSTALLATION.md`.

## Validation Rules

`quick_validate.py` checks:

- `SKILL.md` exists and has valid YAML frontmatter.
- Frontmatter contains `name`, `description`, and optional standard fields only.
- Name is lowercase hyphen-case and 64 characters or fewer.
- Description is not a placeholder, includes `Use when`, and includes `Avoid when`.
- `SKILL.md` is 500 lines or fewer.

Validation does not prove provider discovery. Check or report provider surfaces separately.

## References

- `references/AUTHORING_STANDARDS.md`
- `references/SKILL_TEMPLATE.md`
- `references/workflows.md`
- `references/output-patterns.md`
- `references/PLACEMENT_AND_INSTALLATION.md`
