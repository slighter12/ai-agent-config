# Harness Overview

An execution harness is an optional envelope around agent work. It coordinates phases, owners,
workspace state, verification gates, lifecycle gates, and handoffs. It is not a replacement for
narrower domain skills.

## In Scope

- Coordinate multi-phase work.
- Decide which gate is needed before continuing.
- Route work to the narrowest owner skill or agent role.
- Keep assumptions, blockers, status, and next action visible.
- Route lifecycle capture when accepted decisions, implementation pivots, status or documentation
  drift, capture-worthy handoff notes, loop active state, discussion records, or reusable workflow
  lessons need classification.
- Propose learning candidates when a workflow is reusable.

## Out Of Scope

- Implement code directly.
- Perform full code review.
- Decide detailed testing rules.
- Stage, commit, push, or open PRs.
- Create or update skills.
- Dynamically enable skills or MCPs for subagents.

## Ownership Boundaries

- `policy-core` owns baseline behavior and shared vocabulary.
- `policy-testing` owns verification strategy and gate selection.
- `conventional-git-flow` owns branch, commit, push, PR, and dependency/version closure.
- `code-review` owns code review modes, including bounded diff sanity; harness only coordinates review gates.
- `project-lifecycle` owns phase-end lifecycle signal classification, including accepted decisions,
  implementation pivots, status or documentation drift, capture-worthy handoff notes, loop active
  state, discussion records, and reusable workflow lessons.
- `skill-creator` owns approved skill authoring and validation.
- `mcp-builder-go` owns agent-native MCP/CLI tool design.

## Design Stance

Use the smallest useful harness. Do simple work directly. Apply structured execution only when coordination lowers risk more than it adds ceremony.
