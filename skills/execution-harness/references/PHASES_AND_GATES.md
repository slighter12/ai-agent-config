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
- `lifecycle`: optional only when the harness detects a concrete capture signal; hand facts to `project-lifecycle`, which classifies the target.
- `handoff`: package state for another role or session.

## Gate Routing

- Current diff or code risk review -> `code-review`.
- Missing verification evidence -> responsible task owner with `policy-testing` guidance; the harness only schedules the selected gate.
- Commit readiness or PR flow requested by the user -> `conventional-git-flow`.
- Accepted decision, implementation pivot, status or documentation drift, capture-worthy handoff note, loop active state, discussion record, or reusable workflow lesson -> pass facts as a concrete candidate to `project-lifecycle`.
- Approved shared skill change -> `skill-creator`.
- MCP/CLI tool-surface design -> `mcp-builder-go`.
- Handoff packaging without a capture signal -> `HANDOFF_AND_STATE.md`.
- No remaining gate -> close the phase with stated evidence.
- Use runtime-resolved independent review or challenge roles only when the gate needs independent
  agent execution and such a role exists.

## Phase Discipline

- Each active phase should have one owner.
- Each gate should name the required evidence before continuing.
- Each handoff should include only the context needed by the next owner.
- A failed gate should route back to the smallest responsible owner.
