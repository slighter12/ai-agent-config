---
name: implement
description: Implement an approved local change and run the focused validation needed for confidence. Use when the user explicitly asks to implement an understood spec, ticket, or fix. Avoid when diagnosis, requirements discovery, review, commit, push, or deployment is the primary request.
disable-model-invocation: true
metadata:
  invocation: user
  opencode/autoinvoke: "false"
---

# Implement

Read the controlling spec or ticket, repository instructions, and nearby patterns. Make the smallest complete change, preserving unrelated user work and public contracts.

Run focused, nondestructive validation required by the change. Prefer an existing narrow test or check; do not broaden into a full suite without cause.

Completion is the local change, focused verification, and handoff. Review, commit, push, deployment, and external communication begin only under a separate explicit user request.

Report changed files, evidence actually run, skipped checks, assumptions, and runtime/config/migration risks.
