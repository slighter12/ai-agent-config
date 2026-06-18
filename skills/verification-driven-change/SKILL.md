---
name: verification-driven-change
description: Implement or revise tests and verification loops when tests or executable evidence are the main product or driver of the change. Use when the user asks to add integration tests, add regression coverage, drive a change from a failing test, use TDD/test-first/red-green-refactor, write BDD/acceptance examples, contract/API/e2e tests, characterization tests, or property/invariant checks. Avoid when verification is only a supporting step for ordinary implementation, when the task is pure test strategy without edits, when root cause is unknown and diagnosis is primary, current-diff review, git workflow, or skill authoring.
license: MIT
compatibility: [codex, claude, gemini]
metadata:
  version: "0.1.2"
---

# Verification-Driven Change

## Purpose

Keep test and verification work anchored to executable evidence instead of broad testing advice. This skill owns changes where tests, examples, contracts, or checks are the primary artifact or the driver for implementation.

## Use When

- The user asks to add, revise, or expand integration, contract, API, e2e, regression, characterization, property, invariant, or acceptance tests.
- The user explicitly asks for TDD, test-first, red-green-refactor, BDD, acceptance-driven work, or a failing test before implementation.
- The task is to protect a refactor or legacy change by locking current behavior first.
- The requested output is a verification loop, executable evidence, or test artifact, not only a strategy discussion.

## Avoid When

- The user asks for ordinary implementation and verification is only a supporting step; use `implement-change`.
- The failure or root cause is still unknown and diagnosis is primary; use `diagnose`.
- The user only asks how to choose unit, integration, e2e, fixtures, mocks, or coverage depth without edits; use `policy-testing` or `planning-grill`.
- The user asks to review current git changes; use `code-review`.
- The user asks for commit, branch, push, or PR work; use `conventional-git-flow`.
- The task is to create or update a skill; use `skill-creator`.

## Workflow

1. Identify the verification intent: coverage-add, regression-driven, test-first, contract-driven, characterization-driven, property/invariant-driven, acceptance-driven, or integration-first.
2. Inspect existing test conventions, public behavior boundaries, fixtures, and commands before creating or changing tests.
3. Choose the smallest public behavior boundary that can provide meaningful evidence. Do not test implementation details just because they are easy to reach.
4. Advance one behavior slice at a time:
   - For test-first or regression-driven work, create or describe the failing evidence before implementation.
   - For characterization work, lock current behavior before refactoring.
   - For contract or API work, preserve documented status codes, response shapes, error codes, and compatibility expectations.
   - For property or invariant work, state the invariant before choosing examples or generators.
5. Before finalizing the verification approach, load `policy-testing` for risk level, fixture/mock, determinism, coverage, and validation-gate decisions. Also load the corresponding `policy-*` skill when API contracts, security/auth/secrets/PII, runtime/deploy/config behavior, or language/framework ownership boundaries materially affect correctness, safety, or acceptance.
6. Implement only the scoped production changes needed by the selected verification loop.
7. Run or describe the narrowest useful verification allowed by repo policy, then report evidence, skipped checks, assumptions, and residual risk.

## Guardrails

- Do not treat every request to add tests as TDD. TDD requires an explicit test-first or red-green workflow.
- Do not horizontal-slice a large batch of imagined tests before proving one behavior slice.
- Do not invent a new test framework, runner, fixture system, dependency, or environment unless explicitly approved.
- Do not weaken assertions or broaden fixtures just to make a failing check pass.
- Do not hardcode or log secrets, tokens, credentials, PII, or sensitive request/response bodies in tests.
- If a correct verification seam does not exist, state that as a finding instead of adding a weak or misleading test.
- Load `policy-testing` for detailed risk level, fixture, mock, determinism, coverage, and validation-gate rules.

## Tool And Side-Effect Boundaries

- Do not run tests, programs, formatters, installers, or external services unless the user requested execution or the active repo policy requires it.
- Do not commit, stage, push, deploy, migrate, or perform destructive actions unless explicitly requested.
- Preserve unrelated user changes in the worktree.
- Keep generated test data synthetic and non-sensitive.

## Output

Return:

- `verification_intent`: the selected verification mode.
- `behavior_boundary`: the public behavior or contract under test.
- `change_summary`: test and production changes made or recommended.
- `verification`: commands run, expected evidence, skipped checks, or manual checklist.
- `assumptions`: correctness-relevant assumptions.
- `residual_risks`: remaining uncertainty or follow-up seams.

## Version History

- v0.1.0 (2026-05-19): Initial verification-driven task workflow for test and executable-evidence changes.
- v0.1.1 (2026-05-21): Route current-change review requests to renamed `code-review`.
- v0.1.2 (2026-06-18): Add hard policy gate for testing and material API, security, infra, and language/framework verification risks.
