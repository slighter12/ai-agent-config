---
name: ask-matt
description: Ask which skill or flow fits your current engineering situation. A router over the skills in this repository.
disable-model-invocation: true
metadata:
  invocation: user
  opencode/autoinvoke: "false"
---

# Ask Matt (Workflow Router)

Before stating what any skill does or recommending that a step be skipped, **read that skill's `SKILL.md` first**: the summaries below are for orientation only. The full skill list, with each skill's purpose and source, is in [references/CATALOG.md](references/CATALOG.md).

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

### Context hygiene

Keep steps 1–2 in **one unbroken context window** (don't compact or clear until after `/to-spec` has written the tickets) so the grilling, spec, and tickets all build on the same thinking. Each `/implement` run then starts fresh, working from the ticket or spec. Run `/retro` in the session it's looking back on, before you clear; after clearing, point it at that session's log instead.

The limit on this is the **[smart zone](https://www.aihero.dev/ai-coding-dictionary/smart-zone)**: the part of the window (roughly the first half, as `AGENTS.md` §3 defines it) within which the model still reasons sharply. If a session approaches it before `/to-spec` has written the tickets, don't push on degraded; `/compact` at the nearest phase boundary and carry on (see Phase Boundaries).

## Phase Boundaries

A **phase** is a chunk of work inside a session: the grilling, the implementation, the QA. At the **boundary** between two of them you have five options, and picking between them is the fuzziest decision in this whole map:

- **Continue**: stay put. Costs nothing, loses nothing.
- **`/clear`**: empty the window, when nothing here matters to what's next.
- **`/handoff`** writes a portable markdown file. Narrow: only for a **new harness**, a **new directory**, a **colleague**, or forking a side task **mid-phase**. What it buys is portability.
- **Subagent**: send a tightly-scoped task to its own window and get a report back.
- **`/compact`** compresses this context and seeds a fresh session with it. The **default**, at the bottom of the tree rather than the first reach.

Read [PHASE-BOUNDARIES.md](PHASE-BOUNDARIES.md) for the ordered tree: the five questions, the reasoning behind each branch, and why the primary-source cost makes **Continue** the one to rule out first. Make the decision **at** a boundary; mid-phase, continue or split the rest into subagents.

## Standalone

Off the main flow entirely.

- **`/grill-me`**: the same relentless interview as `/grill-with-docs`, but **stateless**: it saves nothing locally and builds no `GLOSSARY.md`. Reach for it when you are **not working in a working directory** (sharpening a plan, a design, a piece of writing, anything with no repo under it). If you are in a working directory, use `/grill-with-docs` instead: it runs the same interview and leaves a paper trail, so it is strictly the better one.
- **`grilling`** is the interview primitive itself: rounds, the frontier, facts are the agent's job and decisions are yours. `/grill-me` and `/grill-with-docs` are the two named ways in, and `/triage`, `/wayfinder` and `improve-codebase-architecture` all run it internally. Reach for it directly only when you want the interview with no wrapper around it.
- **`prototype`** is a small, throwaway program that answers one design question: does this state model feel right, or what should this UI look like. Throwaway is a constraint on how the code is written, not a promise to destroy it: the answer folds into the real code, and the prototype itself is kept as a **primary source** on a `prototype/<name>` branch out of main, pointed at from the implementation issue. It's the empirical detour under `/grill-with-docs`, but reach for it any time a design question is hard to settle on paper.
- **`research`**: delegate reading legwork to a **background agent**: it investigates a question against **primary sources**, then leaves a cited Markdown result, written into the repo only when the request authorizes it. Keep working while it reads. The result is something to take *into* the main flow at `/grill-with-docs`, since research feeds the thinking rather than replacing it.
- **`/to-questionnaire`** comes in when the thing blocking you isn't in your head or the codebase but in **someone else's**, and it writes them a questionnaire to fill in. It's the inverse of `/grill-me`: instead of interviewing you about the subject, it interviews you about the **send** (who it's going to, what you need back) and aims the questions at the gap. What comes back is material for `/grill-with-docs` or `/to-spec`.
- **`wizard`** is for the steps only a **human** can take: provisioning infrastructure, setting up credentials or CI secrets, clicking through an unfamiliar third-party dashboard, running a one-off migration or cutover. It generates an interactive bash script that opens each URL, captures each value, and writes it into `.env` and GitHub secrets, so the procedure stops being something you re-explain to an agent every time. Model-invoked, so the agent reaches for it the moment it hits a wall only you can pass. If the agent could just do it itself, it should; this is for where a human is genuinely in the loop.
- **`/wait-what`** is the corrective for a message that didn't land. Use it mid-conversation, inside any other skill, and the agent re-pitches what it just said with the context you were missing, in plain English, using the `GLOSSARY.md` vocabulary. It works after the fact; `/grill-with-docs` is the upfront cure, because a shared language agreed early is what stops the jargon arriving at all.
- **`/teach`**: learn a concept over multiple sessions, using the current directory as a stateful workspace.
- **`writing-for-agents`** is the reference for writing documents agents consume: skills, AGENTS.md, pointed-at docs.

## Precondition

**`/setup-matt-pocock-skills`**: run before your first engineering flow to configure the issue tracker, triage labels, and doc layout the other skills assume. Custom issue trackers also work.
