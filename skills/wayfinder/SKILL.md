---
name: wayfinder
description: Determine the next useful step from repository state, specs, tickets, and blocking edges. Use when the user explicitly asks what to do next or where a project stands. Avoid when they already requested a specific implementation action.
disable-model-invocation: true
metadata:
  invocation: user
  opencode/autoinvoke: "false"
---

# Wayfinder

Inspect the tracker declared in `docs/agents/issue-tracker.md`, relevant specs, current git state, and repository evidence. Resolve a local feature root by the tracker's precedence rules; ask for the feature only when multiple roots remain plausible. If tracker configuration is absent, return an explicit `setup-skills` invocation as the next action.

Return a read-only handoff containing current state, blockers, next action, and the exact source of truth. Artifact creation and implementation remain in their owning user-invoked workflows.
