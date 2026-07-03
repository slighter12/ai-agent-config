# Harness Diff Gate

Use this reference when `execution-harness` routes to `code-review` as a bounded diff sanity gate.

## Purpose

Answer whether the current diff is coherent enough to continue, hand off, or prepare for commit.

## Appropriate Timing

- After a phase completes.
- Before commit readiness checks.
- When taking over an existing diff from another session.
- When staged and unstaged changes need coherence checks.
- When dependency, lockfile, generated, or version metadata may be missing from the logical change.

## Still Read-Only

This gate must remain read-only:

- Do not modify files.
- Do not stage files.
- Do not commit, push, or create PRs.
- Do not run tests or programs.
- Do not invoke provider-native review commands.

## Checks

Look for:

- Scope creep.
- Unrelated file churn.
- Missing dependency/version metadata closure.
- Docs, templates, scripts, or config inconsistencies.
- Staged versus unstaged mismatch.
- Missing validation story for the apparent risk level.
- Obvious contract or behavior regressions visible from the diff.

## Boundaries

- Do not become a full code review when this reference is being used as a bounded harness diff gate.
- Do not replace independent review or challenge gates.
- Do not decide commit strategy; use `conventional-git-flow`.
- Do not decide verification depth; use `policy-testing`.
