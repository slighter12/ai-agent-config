# SKILL.md Template

Use this portable template when creating a shared skill for Codex, Claude, Gemini, and OpenCode.

```md
---
name: example-skill
description: Perform [capability]. Use when the user asks for [specific trigger contexts, file types, workflows, or intents]. Avoid when [nearby non-goals or cases better handled by another skill].
metadata:
  version: "0.1.0"
---

# Example Skill

## Purpose

Enable the agent to [outcome], using only the context and resources needed for this task.

## Use When

- The request involves [trigger 1].
- The user needs [trigger 2].
- The task touches [file type, workflow, domain, or tool].

## Avoid When

- The request is only [nearby non-goal].
- Another skill is more specific: [skill-name].
- The task requires provider-specific behavior not covered by this shared skill.

## Workflow

1. Confirm the goal, inputs, constraints, and expected output.
2. Inspect only the relevant files, docs, or resources.
3. Include the minimum task-critical guardrails here; do not rely on another skill being loaded for correctness.
4. Use bundled scripts only when they improve deterministic reliability or avoid repeated code generation.
5. Load reference files only when their topic is directly needed.
6. Stop when the requested outcome is complete, blocked by missing input, or further work would be speculative.

## Tool And Side-Effect Boundaries

- Prefer read-only inspection until the task clearly requires edits or execution.
- Do not run side-effectful commands unless the user asked for them or the active policy allows them.
- Do not create new dependencies, files, or broad refactors unless the task explicitly requires them.
- For destructive, deployment, commit, notification, credential, or external-service workflows, require explicit user confirmation.

## Output

Return:

- `summary`: what was done or recommended.
- `files_touched`: exact paths, if any.
- `assumptions`: only correctness-relevant assumptions.
- `manual_verification`: checklist when execution was skipped or not required.

## Version History

- v0.1.0 (YYYY-MM-DD): Initial portable skill draft.

## References

- `references/INDEX.md` - Use when deeper topic navigation is needed.
```

## Claude-Only Overlay

Use these only in Claude-specific copies or installation overlays, not in the portable shared template:

```yaml
allowed-tools: Read, Grep, Glob
disable-model-invocation: true
user-invocable: false
```

## Quick Checklist

- Name uses lowercase letters, digits, and hyphens only.
- Placement scope is explicit: shared/global, project-local, or provider-specific exception.
- Source of truth, scope, provider surfaces, validation, and install/discovery status are reported.
- Frontmatter includes `name`, `description`, and optional standard fields only.
- Description explains capability, `Use when`, and `Avoid when`.
- Body defines workflow, side-effect boundaries, output, and references.
- Task-critical guardrails are present without duplicating full policy documents.
- Machine-readable version lives in `metadata.version`; human-readable changes live in `Version History`.
- Detailed docs live in `references/`, not duplicated in `SKILL.md`.
- Provider-only controls are kept out of the portable template.
- The skill is validated, and any provider that cannot discover it yet is called out.
