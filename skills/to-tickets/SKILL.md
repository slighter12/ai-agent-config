---
name: to-tickets
description: Split an approved spec into executable tracer-bullet tickets with explicit blocking edges. Use when the user explicitly asks to create tickets or an implementation sequence. Avoid when the spec is not approved or work is too small to split.
disable-model-invocation: true
metadata:
  invocation: user
  opencode/autoinvoke: "false"
---

# To Tickets

Read the approved spec and repository seams. Create the smallest end-to-end tracer bullet first, then only the tickets needed to reach acceptance.

Each ticket must contain outcome, scope, acceptance evidence, dependencies, and `blocks`/`blocked by` edges. Avoid phase-only tickets, speculative abstractions, and triage labels.

Write through the tracker declared in `docs/agents/issue-tracker.md`. For local Markdown, resolve and reuse the feature root by that file's precedence rules before reading `spec.md` or writing `issues/`. If the tracker is missing, unconfigured, ambiguous, or collides with a different feature, stop artifact creation and return the exact missing choice or `setup-skills` invocation.
