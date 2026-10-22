---
name: tidy
description: Simplify recently changed code or a named module for clarity while preserving behavior exactly, verifying after every step. Use when asked to simplify or clean up code, or for an agreed tidy step after tests are green.
metadata:
  invocation: model
---

# Tidy

Make code easier to read, change, and debug without changing what it does. The goal is not fewer lines; it is code a new reader understands faster than the original. `code-review` reports findings without editing; `improve-codebase-architecture` surveys the whole repository; `tidy` edits, but only the recent diff or the module the user names, and only in ways that preserve behavior.

## Boundaries

- Edit only on a change request: the user asked to simplify or clean up code, or agreed to a tidy step. Otherwise report candidates without editing.
- Scope: the changes since a fixed point the user names (or the uncommitted diff), or one named module. No drive-by refactors outside it.
- Before editing, identify the project state in `AGENTS.md` §1. In State 6 (Hotfix / Surgical), do not tidy; report that and stop. In State 2 or State 4, stop and present the Agency Gate Menu from `AGENTS.md` §1 instead of tidying.
- Preserve behavior exactly: inputs, outputs, side effects and their order, error types and messages, and edge cases. Never change public interfaces, error-return shape, state transitions, or I/O seams; those are Tier 2 structural contracts.
- Add no dependencies. Do not touch performance-critical code unless a measurement shows the simpler version is no slower.
- Leave changes uncommitted. Commits, branches, and PRs need their own approval.

## 1. Understand Before Touching

For each candidate, answer: what is its responsibility, what calls it and what does it call, what are its edge cases and error paths, which tests define its behavior, and why it might have been written this way. If you cannot answer, read more first. Before deleting or restructuring logic that no test or comment explains, or any defensive guard (a null check, retry, timeout, clamp, or flag), use the `why` skill (`git log -S`, `git blame`, `docs/adr/`, linked issues) to confirm what constraint it originally guarded. Keep comments that explain an algorithmic, concurrency or ordering, or performance reason.

## 2. Establish the Baseline

Run the tests that cover the scope and record the command and result. If they fail before you start, stop and report. If a candidate lacks behavioral coverage, pause and ask the user whether to add a characterization test first (Agency Gate option (c)) or drop the candidate from this pass.

## 3. Find Candidates

Use the `codebase-design` skill and apply its §3 rules (the Deletion Test, Subtract Before You Add, the 30-Second Test, and Structure over branches) to the scope. Typical signals: deep nesting, long functions with mixed responsibilities, nested ternaries, boolean flag arguments, repeated conditionals, a decision repeated in several places, vague or misleading names, duplicated logic, dead code, and wrappers that add nothing.

Prefer clarity over compactness: explicit control flow beats a dense expression that needs a mental pause. Match the conventions of the file you are editing, not outside preferences.

## 4. Apply One Change at a Time

For each simplification: make the change, rerun the baseline checks, and keep it only when they pass. Revert a change whose checks fail or whose preservation you cannot show; do not stack a second change on an unverified first. Subtract before you add.

## 5. Report

- **Changed**: each simplification with its location and the check that passed after it.
- **Skipped**: candidates you left alone, with the reason (no coverage, unclear purpose, performance-critical, out of scope).
- **Verdict**: Verified, Not verified, or Inconclusive for "behavior preserved", with the final command and result.
