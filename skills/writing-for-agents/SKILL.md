---
name: writing-for-agents
description: Guidelines for authoring high-leverage, low-variance agent documents (SKILL.md, AGENTS.md, CLAUDE.md, and context pointer docs). Use when creating, editing, or auditing skills and steering files.
metadata:
  invocation: model
---

# Writing for Agents

Write instructions that minimize context load, survive the 1–2 skill auto-activation budget, and execute deterministically across both high-reasoning and fast code models.

## 1. Respect the Two Budgets & Activation Mechanics

* **Context Load (Always-On Cost)**: Every line in `AGENTS.md` / `CLAUDE.md` and every word in a skill's `description` frontmatter is paid on **every turn**. Keep `AGENTS.md` / `CLAUDE.md` under 150–200 lines; push reference material and stage procedures into skills or pointer docs.
* **Single-Skill Self-Sufficiency**: Because coding harnesses auto-activate at most 1–2 skills per prompt based on `description` similarity, never split a single workflow's safety invariants across 5 micro-skills. Each stage skill must carry its own compact Pre-Flight Gate.
* **Front-Load Context Pointers**: In `description` fields and `AGENTS.md` pointers, state **what** the document does and the **distinct branches/triggers** for reaching it. Quote YAML `description` values whenever they contain a colon (`:`).

## 2. Translate Abstract Philosophy into Mechanical Heuristics

Models under high context pressure fail on abstract philosophical prose ("write maintainable, decoupled tests") and succeed on **syntactic, countable triggers**:

* *Bad (Abstract)*: "Avoid testing implementation details or over-mocking."
* *Good (Mechanical)*: "If a newly extracted helper has only `1` caller in production code and exists only so the test can import it, inline it back into its caller. If a test configures `> 2` mocks or only asserts `mock.calledWith` without asserting a returned value or state change, delete the unit test and run runtime verification."

## 3. Positive Steering & Checkable Completion Criteria

* **Prompt the Positive Target**: State the exact behavior to perform rather than listing negations ("do not..."), reserving hard prohibitions only for explicit redlines paired with their positive replacement.
* **Checkable Completion Criteria**: End every procedural step on an observable state (e.g., `"Phase 1 is complete when one command deterministically reproduces the failure"`, `"Stop and wait for user confirmation before editing files"`).
* **Prune No-Ops**: Delete any sentence the model already obeys by default (e.g., `"Write clean, bug-free code"`).
