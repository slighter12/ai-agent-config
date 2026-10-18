---
name: retro
description: Conduct a post-session retrospective to improve the repo's agent environment, graduate user style/migration decisions into project rules or checks, and prune no-op instructions.
disable-model-invocation: true
metadata:
  invocation: user
  opencode/autoinvoke: "false"
---

# Session Retrospective (`/retro`)

Inspect the current coding session (or a specified session log) to improve the **environment the agent works in**, not the feature code itself.

Load and follow the `writing-for-agents` skill (`../writing-for-agents/SKILL.md`) before proposing any edits to `AGENTS.md`, `CLAUDE.md`, `CODING_STANDARDS.md`, or skills.

## Process

### 1. Trace Real Session Friction & User Decisions
Scan the session history for concrete moments where:
* The agent had to pause at the **Pre-Flight Agency Gate** to ask the user how to resolve an ADR migration conflict, multi-author style drift, or legacy monolith coupling.
* The agent wasted tool calls hunting for files, build/test commands, or domain terms.
* The agent made a mistake (or wrote redundant/coupled tests) that a deterministic check or reviewer rule could have prevented.

Discard generic best practices that cannot be traced to a specific moment in the session.

### 2. Classify Findings into the Right Layer (5 Layers of Correction)

Route every proposed improvement to the lowest layer that can enforce it, keeping `AGENTS.md` / `CLAUDE.md` lean:

| Session Finding | Where the Fix Must Land |
| :--- | :--- |
| **Mechanical / Syntactic Violation**<br>(Banned API, import boundary, formatting, missing test/lint run) | **Layer 1/2 — Deterministic Check**: Propose a linter rule, type constraint, pre-commit hook, or CI check (or wire an existing unwired script). A check can fail; prose cannot. |
| **User Decision on Style Drift or Active Migration**<br>(e.g., "Use Pattern B for new modules; use Sprout for legacy package X") | **Layer 2/3 — Project Law Graduation**: Propose a concise entry in `docs/adr/`, `CODING_STANDARDS.md`, or a 1-line pointer in the project's `AGENTS.md` / `CLAUDE.md` so future sessions inherit the decision automatically without re-asking. |
| **Judgement-Call Code/Review Standard**<br>(Nuanced convention no linter can catch) | **Reviewer Context (`CODING_STANDARDS.md`)**: Place it where `code-review` reads it (low context pressure), not in `AGENTS.md` (which taxes every implementation turn). |
| **Slow File / Symbol Discovery** | **Navigation Pointer**: Add a 1-line context pointer in `AGENTS.md` or `GLOSSARY.md` pointing to the non-obvious entrypoint. |
| **Missing Domain Term or Ambiguous Noun** | **Domain Glossary (`GLOSSARY.md`)**: Record the canonical term and definition. |
| **Bloated or No-Op Lines in `AGENTS.md` / `CLAUDE.md`** | **Pruning**: Propose deleting instructions the model already obeys by default or moving long procedures out into on-demand skills. |

### 3. Present Ranked Proposals & Wait for Approval

Present the candidate environment improvements ordered by severity (`Moment in Session -> Proposed Layer & Target File -> Exact Diff/Addition`).

**Hard Stop**: Do not edit any file or install any hook until the user explicitly selects which candidates to apply.
