# Phases And Gates

Use phases as coordination vocabulary, not as mandatory ceremony. Skip phases that do not reduce risk for the active task.

## Common Phases

- `frame`: define objective, constraints, done criteria, and blockers.
- `inspect`: gather local evidence before deciding.
- `plan`: choose owners, order, and acceptance evidence.
- `execute`: make bounded changes or delegate bounded work.
- `verify`: select and gather evidence required by risk.
- `review`: run bounded sanity, risk review, or challenge pass when needed.
- `commit`: prepare staged files, metadata closure, message, and PR flow when requested.
- `lifecycle`: classify phase-end capture signals such as accepted decisions, implementation
  pivots, status or documentation drift, capture-worthy handoff notes, loop active state,
  discussion records, or reusable workflow lessons.
- `handoff`: package state for another role or session.
- `capture`: hand off lifecycle capture candidates to `project-lifecycle`.

## Gate Routing

- Use `policy-testing` for verification depth and gate selection.
- Use `code-review` for code review modes and bounded diff sanity gates.
- Use `conventional-git-flow` for commit, branch, push, PR, and dependency/version closure.
- Use `project-lifecycle` for phase-end lifecycle signal classification: accepted decisions,
  implementation pivots, status or documentation drift, capture-worthy handoff notes, loop active
  state, discussion records, and reusable workflow lessons.
- Use `skill-creator` for approved shared skill changes after explicit user approval.
- Use `mcp-builder-go` for MCP/CLI tool-surface design.
- Use reviewer or oracle roles only when a review or challenge gate needs independent agent execution.

## Phase Discipline

- Each active phase should have one owner.
- Each gate should name the required evidence before continuing.
- Each handoff should include only the context needed by the next owner.
- A failed gate should route back to the smallest responsible owner.
