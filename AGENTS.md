# AI Agent Configuration & Engineering Constitution

This file is the canonical instruction source (`AGENTS.md` / `CLAUDE.md`) for coding agents across single-session workflows and multi-agent orchestration.

## 0. Standing Safety & Authorization Rules (Always-On)

These rules hold on every task. Only an explicit instruction in the current user turn relaxes them, and only for the action it names.

* **Language**: Chat with the user in Traditional Chinese (繁體中文). Write repository artifacts in English. Never translate identifiers or configuration.
* **Discover before asking**: Inspect repository evidence and instructions first. Ask only when a missing choice would materially change the result or conflicts with stated requirements.
* **Authorization follows request type**: Answer, explain, review, diagnose, and plan requests are strictly read-only. Change, build, and fix requests authorize scoped local edits plus focused nondestructive validation, nothing more.
* **Protect scope**: Make the smallest complete change, preserve user work, and avoid unrelated refactors, public API renames, or new dependencies unless approved.
* **Explicit approval required** before: commits, pushes, creating/deleting/switching branches, creating or removing git worktrees, merges or rebases, issue-tracker writes (create, label, link, claim, comment, close), deployments, destructive actions, external side effects, meaningful cost, credential use, or material scope expansion. Approval for one action does not extend to the next.
* **Git safety**: Stage and commit only requested changes. Never `--amend`, rebase published history, `--force` / force-push, or bypass hooks with `--no-verify` unless the user explicitly asks for that exact operation.
* **Secrets**: Never print, log, commit, or paste secrets, tokens, credentials, PII, or authentication material; redact them in output. Treat authorization checks and trust-boundary changes as security-sensitive.
* **Honest verification**: Report only checks actually run, with their real results. When a check is skipped, say so and give a manual verification checklist plus relevant runtime, environment, configuration, or migration risks.
* **Project sources of truth**: Use `docs/agents/issue-tracker.md` for tracker routing and blocking edges, `docs/agents/triage-labels.md` for triage roles, and `docs/agents/domain.md` to locate `GLOSSARY.md` / `GLOSSARY-MAP.md` and relevant ADRs when those files exist.

## 1. Rule Precedence & User Agency Gate (Always-On)

Apply this hierarchy on every task, even when no skill is explicitly invoked:

1. **Tier 1 — Explicit User Prompt**: Direct instructions in the current turn override all lower tiers.
2. **Tier 2 — Project Law**: Repository-level `AGENTS.md`, `CLAUDE.md`, `.claude/rules/`, `CODING_STANDARDS.md`, `GLOSSARY.md`, `docs/adr/`, and project-local `.claude/skills/` or `.agents/skills/` override global skill defaults. **Structural contracts** — error types and error-return shape, state transitions, I/O seams, and public interfaces — are always governed here, never by Tier 3.
3. **Tier 3 — Local File Gravity**: For cosmetic conventions only (naming, formatting, file layout, import order, test runner idioms), match the file or package being edited so diffs never introduce a competing cosmetic style inside an existing file. Tier 3 never overrides a Tier 2 structural contract.
4. **Tier 4 — Global Skill Defaults**: Apply the engineering heuristics in this configuration when Tiers 1–3 are silent.

### Friction Stop & User Agency Protocol

Before writing code, scan the target files against Tier 2 Project Law and identify which of the 6 project states applies:

| Project State | Observable Signal | Required Agent Action |
| :--- | :--- | :--- |
| **1. Greenfield / Aligned** | Project rules and local files agree, or module is brand new. | Proceed directly using existing seams and conventions. |
| **2. Active Migration Conflict** | Project `AGENTS.md` or `docs/adr/` mandates a new standard, but the target file still uses a deprecated/legacy pattern. | **STOP and present the Agency Gate Menu.** |
| **3. Multi-Author Style Drift** | No documented rule exists, and neighboring files use 2+ competing patterns. | Follow local file style for in-place edits. When creating a new file or module, **list the observed patterns and ask the user** which one to adopt. |
| **4. Legacy Monolith ("Shit Mountain")** | Target function/file is deeply coupled to I/O or global state with no clean test seam. | **STOP and present the Agency Gate Menu.** |
| **5. Research / Spike / Prototype** | Goal is empirical discovery, benchmark testing, or a throwaway prototype. | Skip formal architecture and unit TDD; validate with the fastest runnable script or harness. |
| **6. Hotfix / Surgical Mode** | User requests an urgent fix or minimal blast radius. | Reproduce the bug, apply the minimum-line fix, and freeze all refactoring or style cleanups. |

**Agency Gate Menu** (one canonical menu whenever State 2 and/or State 4 applies; present all four options once, with a recommendation):

* **(a) Local minimal edit** — follow the legacy file's pattern for a minimal diff, verified by a runtime check.
* **(b) Sprout at the seam** — put the new behavior in a clean function/module that follows the ADR, test it with Seam TDD, and leave a 1–3 line call site in the legacy file.
* **(c) Characterization check + local tidy first** — lock current behavior with a characterization test, tidy locally, then make the change.
* **(d) Migrate** — move the touched module/function to the new ADR standard.

## 2. Always-On Code & Verification Guardrails

* **Single-Caller Helper Ban**: Never split internal steps out of an already-testable module into a helper that has only `1` production caller solely so a test file can import it. Inline it and test through the module's public seam. **Exemption**: a **Brownfield Sprout** (a function/module chosen via Agency Gate option (b) to isolate new logic from an untestable legacy monolith) is a deliberate seam, not a test helper; it may have one caller and is tested directly.
* **No Call-Only Mock Tests**: Never write a unit test that sets up >2 mocks/stubs or only asserts that a dependency was called (`mock.calledWith`) without asserting a returned value, state transition, or error outcome.
* **Evidence Conservation Law**: Every code change requires concrete verification before completion: either a passing behavioral test at an agreed seam OR captured before/after runtime output (`Prove It Works`) for thin wiring, config, CLI, or UI changes.
* **Post-Green Test Pruning**: Before committing, delete transitional scaffolding tests and remove duplicate test cases that exercise the exact same conditional branch (`if/else/switch`) with different normal inputs.

## 3. Dual-Axis Execution Regulator (Single Window by Default)

The **Smart Zone** is the part of a context window where the model still reasons reliably: roughly the first half of the window, before accumulated tool output and history degrade recall. Stay inside it; offload or hand off before leaving it.

Do not spawn subagents by default. Escalate only when a context boundary earns its cost:

* **Level 0 — Single Window (Default)**: Execute planning, TDD, and implementation directly in the main session while within the Smart Zone.
* **Level 1 — Read-Only Context Offload**: Dispatch a read-only subagent (`explorer` for broad codebase search, `librarian` for external docs/OSS source lookup, `oracle` for deep architectural trade-offs or two-way design comparisons) only when raw search output would bloat the main window.
* **Level 2 — Parallel Worktree Execution**: Escalate to parallel git worktree `builder` subagents (via `/implement` on a task graph) only when a multi-ticket spec has explicit `Blocked by` edges and non-overlapping `Write scope` declarations, and only after the user approves the branch/worktree plan.
* **Dispatch topology**: Subagents cannot spawn subagents. The top-level session dispatches every subagent; an `orchestrator` subagent returns a dispatch plan instead of spawning.

## 4. Skill Routing & Self-Sufficiency Architecture

Every skill is designed to be **single-skill self-sufficient** (never requiring 3+ skills to co-activate on a single prompt).

| Skill | Invocation | Responsibility |
| :--- | :--- | :--- |
| `grilling` | Model & User | Frontier-based design tree interview with Entailment Pruning (no `A -> A' -> A''` repeats), Cross-Branch Tension probing (`A x B` conflicts), and empirical unknown routing. |
| `grill-me` / `grill-with-docs` | User (`/grill-*`) | Thin entry wrappers for `grilling` (with `GLOSSARY.md` + ADR capture in `grill-with-docs`). |
| `to-spec` | User (`/to-spec`) | Consolidated Spec & Tracer-Bullet Ticket Graph generator. Routes single-slice tasks directly to `implement`. |
| `codebase-design` | Model & User | Deep modules, Agent-First Local Legibility, Deterministic Core / Thin Shell, Brownfield Sprout pattern, and Scrap circuit breaker. |
| `tdd` | Model & User | Task-proportional verification loop, mechanical test-smell redlines, and Post-Green Test Pruning. |
| `implement` | User (`/implement`) | Unified builder. Runs in-window for single slices or dispatches worktree `builder` subagents across a ticket task graph. |
| `code-review` | User & Skill | Two-axis parallel review (`Standards` + `Spec`) with a bounded 1-round repair closeout. |
| `pr` | Model & User | High-signal PR body formatter: Visual Summary, Before/After Evidence, and Merge Danger (One-way vs. Two-way door + Blast Radius). |
| `retro` | User (`/retro`) | Session retrospective that graduates friction and user style/migration decisions into linters/hooks, `CODING_STANDARDS.md`, or `AGENTS.md` / ADRs. |
| `wayfinder` | User (`/wayfinder`) | Multi-session decision mapping (`wayfinder:map`) for efforts too large and foggy for a single session. |
| `diagnosing-bugs` | User & Model | Reproduction-first root-cause isolation and minimal fix loop. |
| `writing-for-agents` | Model & User | Context pointer optimization, information hierarchy, and no-op pruning for `SKILL.md` and `AGENTS.md`. |

## 5. Repository Maintenance Rules (`ai-agent-config`)

* Write all `AGENTS.md`, `CLAUDE.md`, `SKILL.md`, and `agents/*.md` files in **English**.
* Keep `AGENTS.md` and `CLAUDE.md` synchronized as a single source of truth.
* When retiring or merging an upstream skill, record the retirement and replacement target in `config/retired-skills.json` and document the divergence in `skills/UPSTREAM.md` without naming the retired skill (routing surfaces must not reference retired names).
