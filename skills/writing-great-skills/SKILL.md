---
name: writing-great-skills
description: Create, revise, validate, or package a concise portable CLI skill with explicit invocation semantics. Use when the user explicitly asks to author or improve a SKILL.md-based skill. Avoid when they only want to use or install an existing skill.
disable-model-invocation: true
metadata:
  invocation: user
  opencode/autoinvoke: "false"
---

# Writing Great Skills

Write the smallest instruction set that teaches durable judgment the model does not already have. Put routing in the description, keep one shared body, and add references only for conditional detail.

Choose `metadata.invocation: user` or `model`. For user skills add `disable-model-invocation: true`, `opencode/autoinvoke: "false"`, and `agents/openai.yaml` with `policy.allow_implicit_invocation: false`.

Audit the skill before validation:

- Front-load a compact leading word that anchors the intended trigger or execution behavior.
- Keep one distinct trigger per branch. Inline what every branch needs; disclose branch-only detail behind a precise reference pointer.
- End each step with a checkable completion criterion, exhaustive where premature completion would be costly.
- Keep each meaning in one authoritative place. Apply the no-op test sentence by sentence and delete instructions that do not change default behavior.
- State the target behavior positively. Keep negation only for hard safety boundaries and pair it with the safe action.

Resolve `<config-repo>` from this skill's canonical source path. Use absolute paths for skill and output arguments with the deterministic CLI:

- `go -C <config-repo>/hooks-go run ./cmd/agent-config --repo-root <config-repo> init-skill <name> --path <absolute-dir>`
- `go -C <config-repo>/hooks-go run ./cmd/agent-config --repo-root <config-repo> validate-skill <absolute-skill-dir>`
- `go -C <config-repo>/hooks-go run ./cmd/agent-config --repo-root <config-repo> validate-skills <absolute-skills-dir>`
- `go -C <config-repo>/hooks-go run ./cmd/agent-config --repo-root <config-repo> package-skill <absolute-skill-dir> [<absolute-output-dir>]`

Load [local provider placement](references/LOCAL_PROVIDER_PLACEMENT.md) only when installing or integrating a skill across providers.
