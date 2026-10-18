---
name: ask-matt
description: Ask which skill or flow fits your current engineering situation. A router over the skills in this repository.
disable-model-invocation: true
metadata:
  invocation: user
  opencode/autoinvoke: "false"
---

# Ask Matt (Workflow Router)

Before stating what any skill does or recommending that a step be skipped, **read that skill's `SKILL.md` first**: the summaries below are for orientation only.

## The Main Flow: Idea -> Ship

1. **`/grill-with-docs` (or `/grill-me` outside a repo)**:
   * Relentless design-tree interview with Entailment Pruning (no `A -> A' -> A''` repeats) and Cross-Branch Tension probing (`A x B` conflicts). Records domain terms in `GLOSSARY.md` and architectural constraints in `docs/adr/`.
   * *Empirical detour*: If a question cannot be settled by talking, run `/prototype` (for runnable UI/state experiments) or `/research` (for primary-source doc/SDK investigation).
2. **Single-Slice vs. Multi-Slice Branch**:
   * **Fits in one session (Single Slice)**: Skip `/to-spec` and run `/implement` directly in the current window.
   * **Spans multiple slices/sessions**: Run `/to-spec`, which synthesizes both the **Spec (with Agreed Test Seams)** and the **Tracer-Bullet Ticket Graph (`Blocked by` + `Write scope`)** in a single pass.
3. **`/implement` (Unified Builder)**:
   * Pass a single ticket (`/implement 01`) to build one slice in-window with `/tdd` (or `Prove It Works` runtime verification for thin wiring) and `/code-review`.
   * Pass a whole multi-ticket spec (`/implement <spec>`) to orchestrate parallel `builder` subagents in isolated git worktrees across the unblocked frontier onto a single integration branch.
4. **`/code-review` -> `/pr` -> `/retro`**:
   * `/code-review` checks Standards + Spec compliance in parallel.
   * `/pr` formats the pull request body (Visual Summary, Before/After Evidence, One-way/Two-way Door Merge Danger).
   * `/retro` inspects session friction and graduates user style/migration decisions into linters, `CODING_STANDARDS.md`, or `AGENTS.md` / ADRs.

## On-Ramps & Codebase Health

* **Incoming raw bug reports / external PRs**: `/triage`
* **Hard bug or performance regression**: `/diagnosing-bugs` -> `/retro`
* **Multi-session foggy effort**: `/wayfinder` (resolves `wayfinder:<type>` decision tickets, then merges into `/to-spec`)
* **Architectural survey & module deepening**: `/improve-codebase-architecture` + `/codebase-design`

## Phase Boundaries

At the boundary between two phases (e.g., Grilling -> Implementation), consult `PHASE-BOUNDARIES.md` in this directory (`Continue > /clear > /handoff > Subagent > /compact`).
