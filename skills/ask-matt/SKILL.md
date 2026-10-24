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

## Output Contract

The answer to `/ask-matt` is a route, never a solution. This holds for every turn of the conversation `/ask-matt` started, not only the first reply. Reply with only:

1. **Position**: the first Routing Question below that applies, with the evidence (files read, ticket or spec found).
2. **Route**: the skill sequence with the exact command to run next (for example `/grill-with-docs` -> `/implement`, or `/implement .scratch/<feature-slug>/issues/01-<slug>.md`).
3. **Why**: one line on why this route fits.

Do not solve the problem, compare implementation options, propose code, or edit files — even when the message describes a concrete task, reads as a command ("fix this", "handle this"), or asks how to fix something. Reading files to find the position is allowed. Stop after the recommendation; the user runs the next command.

When the user answers a question you asked, the reply is still a route: give the next command with the answer in it and stop. The answer settles the route; it never authorizes starting the work here.

## Routing Questions

Ask these in order; the first one that applies picks the route.

1. **Does an on-ramp fit?** A bug, incoming issues, stacking patches, or the other situations under On-Ramps & Codebase Health take that route.
2. **Is the work already specced or ticketed?** Resume where it stands. When a spec already has two or more unresolved child tickets, recommend `/implement <spec>`, not one `/implement` per ticket, however small the tickets are. For a single ticket, recommend `/implement` on that ticket.
3. **Can one session hold the open decisions?** If not — the way to the destination is not visible yet and settling it takes more than one session — run `/wayfinder`; it resolves the decisions and hands off to `/to-spec`. A well-scoped feature never goes here, however many design questions it has: those are grilling questions.
4. **Are the decisions settled?** If not, run `/grill-with-docs` (`/grill-me` outside a repo).
5. **One slice or several?** One slice: `/implement` in the current window, skipping `/to-spec`. Several slices: `/to-spec`, then `/implement <spec>`.

**Project state belongs to `/implement`.** Do not classify the 6 project states or present the Agency Gate Menu here: `/implement` checks the target files against Project Law and runs the gate before writing code. When you see friction (an ADR the target code contradicts, or a coupled legacy file), say in the **Why** line that `/implement` will stop at its Agency Gate.

## The Skills Along the Route

1. **`/grill-with-docs` (or `/grill-me` outside a repo)**:
   * Relentless design-tree interview with Entailment Pruning (no `A -> A' -> A''` repeats) and Cross-Branch Tension probing (`A x B` conflicts). Records domain terms in `GLOSSARY.md` and architectural constraints in `docs/adr/`.
   * *Empirical detour*: If a question cannot be settled by talking, use `prototype` (for runnable UI/state experiments) or `research` (for primary-source doc/SDK investigation).
2. **`/to-spec`**: synthesizes the **Spec (with Agreed Test Seams)** and the **Tracer-Bullet Ticket Graph (`Blocked by` + `Write scope`)** in a single pass.
3. **`/implement` (Unified Builder)**: pass the target; `implement` picks the mode from it.
   * **Mode A** — a single ticket (`/implement .scratch/<feature-slug>/issues/01-<slug>.md`) builds one slice in-window with the `tdd` skill (or `Prove It Works` runtime verification for thin wiring) and the `code-review` skill.
   * **Mode B** — a spec with two or more unresolved child tickets (`/implement <spec>`) plans parallel `builder` subagents in isolated git worktrees onto one integration branch and asks before creating any; if the user declines, it builds the children one at a time in `Blocked by` order.
4. **`code-review` -> `pr` -> `/retro`**:
   * `code-review` checks Standards + Spec compliance in parallel.
   * `pr` formats the pull request body (Visual Summary, Before/After Evidence, One-way/Two-way Door Merge Danger).
   * `/retro` inspects session friction and graduates user style/migration decisions into linters, `CODING_STANDARDS.md`, or `AGENTS.md` / ADRs.

## On-Ramps & Codebase Health

* **Incoming raw bug reports / external PRs**: `/triage`
* **Hard bug or performance regression**: `diagnosing-bugs` -> `/retro`
* **Urgent fix with a known cause**: `/implement` directly, stating that it is a hotfix.
* **Spike or throwaway experiment**: `prototype`
* **Before changing unfamiliar code**: `how` (what it does and where it lives), then `why` (why it is shaped this way, ending in Preserve / Change / Avoid / Risk). Use `research` instead for external docs, SDKs, or specs.
* **Risky or untrusted change before merge**: `blast-radius` (breakage beyond the diff, with the safety fact proven), alongside `code-review` for standards and spec.
* **No scripted way to prove app behavior**: `verify-app` (create or maintain a project-local verification skill).
* **Green but heavier than needed**: `tidy` (behavior-preserving simplification of the recent diff or a named module).
* **Prose before it ships** (README, docs, commit message, PR body): `unslop`.
* **Patches keep stacking on one task**: `/step-back` (test the shared premise, compare a fresh design, then you choose).
* **Architectural survey & module deepening**: `/improve-codebase-architecture` + `codebase-design`

## Phase Boundaries

At the boundary between two phases (e.g., Grilling -> Implementation), consult `PHASE-BOUNDARIES.md` in this directory (`Continue > /clear > /handoff > Subagent > /compact`).
