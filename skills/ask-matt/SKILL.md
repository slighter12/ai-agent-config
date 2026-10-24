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

1. **Scale**: one slice, several slices, or an effort too big for one session whose way is still foggy. Decide the scale before the project state; the project state never changes it. A foggy multi-session effort routes to `/wayfinder` first, even in State 2 or State 4: the Agency Gate applies later, when a ticket is built.
2. **Project state**: which of the 6 project states in `AGENTS.md` the work is in, with the evidence (files read, ticket or spec found).
3. **Route**: the skill sequence with the exact command to run next (for example `/grill-with-docs` -> `/implement`, or `/implement .scratch/<feature-slug>/issues/01-<slug>.md`).
4. **Why**: one line on why this route fits.

Do not solve the problem, compare implementation options, propose code, or edit files — even when the message describes a concrete task, reads as a command ("fix this", "handle this"), or asks how to fix something. Reading files to detect the project state is allowed. Stop after the recommendation; the user runs the next command.

When the user answers a question or picks a menu option you offered (for example an Agency Gate choice), the reply is still a route: give the next command with that choice in it (for example `/implement <target> (Agency Gate: b)`) and stop. The choice settles the route; it never authorizes starting the work here.

## The Main Flow: Idea -> Ship

Too big for one session and still foggy (for example a migration across several services with open decisions)? Start with `/wayfinder` instead of step 1; it resolves the decisions and hands off to `/to-spec`.

1. **`/grill-with-docs` (or `/grill-me` outside a repo)**:
   * Relentless design-tree interview with Entailment Pruning (no `A -> A' -> A''` repeats) and Cross-Branch Tension probing (`A x B` conflicts). Records domain terms in `GLOSSARY.md` and architectural constraints in `docs/adr/`.
   * *Empirical detour*: If a question cannot be settled by talking, run `prototype` (for runnable UI/state experiments) or `research` (for primary-source doc/SDK investigation).
2. **Single-Slice vs. Multi-Slice Branch**:
   * **Fits in one session (Single Slice)**: Skip `/to-spec` and run `/implement` directly in the current window.
   * **Spans multiple slices/sessions**: Run `/to-spec`, which synthesizes both the **Spec (with Agreed Test Seams)** and the **Tracer-Bullet Ticket Graph (`Blocked by` + `Write scope`)** in a single pass.
3. **`/implement` (Unified Builder)**: pass the target; `implement` picks the mode from it.
   * **Mode A** — a single ticket (`/implement .scratch/<feature-slug>/issues/01-<slug>.md`) builds one slice in-window with the `tdd` skill (or `Prove It Works` runtime verification for thin wiring) and the `code-review` skill.
   * **Mode B** — when a spec already has two or more unresolved child tickets, recommend `/implement <spec>`, not one `/implement` per ticket, however small the tickets are. It plans parallel `builder` subagents in isolated git worktrees onto one integration branch and asks before creating any; if the user declines, it builds the children one at a time in `Blocked by` order.
4. **`code-review` -> `pr` -> `/retro`**:
   * `code-review` checks Standards + Spec compliance in parallel.
   * `pr` formats the pull request body (Visual Summary, Before/After Evidence, One-way/Two-way Door Merge Danger).
   * `/retro` inspects session friction and graduates user style/migration decisions into linters, `CODING_STANDARDS.md`, or `AGENTS.md` / ADRs.

## On-Ramps & Codebase Health

* **Incoming raw bug reports / external PRs**: `/triage`
* **Hard bug or performance regression**: `diagnosing-bugs` -> `/retro`
* **Before changing unfamiliar code**: `how` (what it does and where it lives), then `why` (why it is shaped this way, ending in Preserve / Change / Avoid / Risk). Use `research` instead for external docs, SDKs, or specs.
* **Risky or untrusted change before merge**: `blast-radius` (breakage beyond the diff, with the safety fact proven), alongside `code-review` for standards and spec.
* **No scripted way to prove app behavior**: `verify-app` (create or maintain a project-local verification skill).
* **Green but heavier than needed**: `tidy` (behavior-preserving simplification of the recent diff or a named module).
* **Prose before it ships** (README, docs, commit message, PR body): `unslop`.
* **Patches keep stacking on one task**: `/step-back` (test the shared premise, compare a fresh design, then you choose).
* **Multi-session foggy effort**: `/wayfinder` (resolves decision tickets, then hands off to `/to-spec`)
* **Architectural survey & module deepening**: `/improve-codebase-architecture` + `codebase-design`

## Phase Boundaries

At the boundary between two phases (e.g., Grilling -> Implementation), consult `PHASE-BOUNDARIES.md` in this directory (`Continue > /clear > /handoff > Subagent > /compact`).
