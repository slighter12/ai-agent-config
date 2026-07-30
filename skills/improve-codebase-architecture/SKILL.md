---
name: improve-codebase-architecture
description: Improve an existing codebase's dependency shape through a bounded architecture change. Use when the user explicitly asks to reduce coupling, clarify boundaries, or repair architecture. Avoid when they only want an assessment or a feature implementation.
disable-model-invocation: true
metadata:
  invocation: user
  opencode/autoinvoke: "false"
---

# Improve Codebase Architecture

Use `codebase-design` once to map the current dependency direction, ownership, data flow, and painful seam. Reuse that map and state the invariant the change should restore.

Make the smallest architecture edit that removes the concrete pressure. Prefer deleting indirection, moving ownership, or exposing one clear seam over introducing a framework.

Preserve behavior and public APIs unless the user approved a migration. Run focused validation and report residual coupling.
