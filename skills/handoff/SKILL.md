---
name: handoff
description: Create a durable cross-session handoff from current repository and tracker state. Use when the user explicitly asks to pause, transfer, or resume work later. Avoid when a normal final response is enough.
disable-model-invocation: true
metadata:
  invocation: user
  opencode/autoinvoke: "false"
---

# Handoff

Capture the objective, completed work, current branch and worktree state, decisions, evidence run, open risks, blockers, and exact next action.

Write to the configured tracker or project handoff location. Keep it factual and reconstructable, use that location as the sole planning source, and list only checks actually run.
