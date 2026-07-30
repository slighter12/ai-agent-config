---
name: codebase-design
description: Analyze or design code ownership, dependencies, data flow, and change seams before implementation. Use when a planning or architecture workflow needs bounded codebase design judgment. Avoid when the task is already a small local edit.
metadata:
  invocation: model
---

# Codebase Design

Start from the behavior being changed. Identify the current owner, dependency direction, state boundary, and the narrowest seam that can carry the change.

Prefer cohesive ownership, explicit data flow, stable contracts, and deletion of accidental indirection. Introduce an abstraction only when at least two concrete pressures require it now.

Complete when each behavior in scope has a current owner, dependency direction, state boundary, proposed seam, migration risk, and focused validation point; list any remaining unknowns. Return the proposed shape, affected modules, and invariants. Implementation begins only when the enclosing request authorizes edits.
