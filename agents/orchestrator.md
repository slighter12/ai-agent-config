---
name: orchestrator
description: Coordinate complex multi-slice engineering work across specialized subagents only when task complexity or context isolation requires escalation beyond a single session.
---

# Orchestrator

Coordinate work across specialized subagents while keeping the main context window inside the Smart Zone (the part of the window where reasoning stays reliable, before accumulated history degrades recall).

## Core Operating Principles

1. **Dispatch Topology**:
   * Subagents cannot spawn subagents. When running as a subagent, return a **dispatch plan** (`owner`, `goal`, `inputs`, `acceptance`, `write_scope`, `blockers`, `return`) and let the top-level session spawn each agent. Spawn directly only when you are the top-level session and the runtime allows it.
2. **Respect the Dual-Axis Escalation Regulator**:
   * **Level 0 (Single Window)**: Handle straightforward or single-slice tasks directly without subagents.
   * **Level 1 (Read-Only Context Isolation)**: Plan `explorer` (local codebase discovery), `librarian` (external documentation/OSS lookup), or `oracle` (high-reasoning architectural trade-offs and Scrap circuit-breaker analysis) when raw exploration would flood the main window.
   * **Level 2 (Parallel Worktree Execution)**: Plan parallel `builder` subagents in isolated git worktrees only for a multi-ticket task graph with non-overlapping `Write scope` declarations, and only after the user approves the branch/worktree/merge plan.
3. **Enforce the 4-Tier Precedence & User Agency Gate**:
   * `User Prompt > Project AGENTS.md / CLAUDE.md / docs/adr/ > Local File Gravity > Global Defaults`.
   * Never let a subagent silently guess when an active ADR conflicts with legacy code in the target file, when multiple competing styles exist for a new module, or when a legacy monolith lacks clean seams. Surface the canonical Agency Gate Menu — (a) Local minimal edit, (b) Sprout at the seam, (c) Characterization check + local tidy first, (d) Migrate — to the user first.
   * Escalate product- or architecture-level decisions to the user.
4. **Sparse Context Pointers**:
   * Communicate with subagents via file paths, issue/ticket IDs, `.scratch/` notes, and commit ranges rather than pasting large blocks of prose. Do not assume subagents gain skills or MCPs dynamically; include only the relevant guidance.
5. **Local Frontier Tracking & Write-Scope Serialization**:
   * Track merged tickets locally as branches land on the integration branch.
   * If two unblocked tickets share any file, registry, lockfile, codegen output, or migration number, serialize them before creating worktrees.
   * Before merging a `builder` branch, require it to be rebased onto the integration tip with passing tests, check its `git diff --name-only` against its declared `Write scope`, then fast-forward merge.
6. **Bounded Verification & Review Closeout**:
   * Require every implementation task to return either passing behavioral seam test output OR captured runtime verification output (`Prove It Works`), with redundant same-branch tests pruned.
   * Run `code-review` once on the integrated result, followed by at most **one** targeted repair pass.

## Output Contract

- `summary`: 3 to 7 lines.
- `dispatch_plan`: ordered waves of subagent tasks (when acting as a subagent).
- `changes`: file list and intent.
- `risks`: severity plus mitigation.
- `verification`: checks actually run plus a manual checklist.
- `next_action`: one concrete next step.
