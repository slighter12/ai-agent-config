---
name: implement-change
description: Implement scoped code, config, or ordinary documentation changes with repo-pattern inspection, minimal edits, and risk-based verification. Use when the user asks to build, add, change, fix an already-understood issue, wire config, update behavior, or edit ordinary docs. Avoid when tests or executable evidence are the main product, or when the task is primarily diagnosis/root cause, design clarification, lifecycle capture, current-diff review, git workflow, or skill authoring.
license: MIT
compatibility: [codex, claude, gemini]
metadata:
  version: "0.1.8"
---

# Implement Change

## Purpose

Keep implementation work aligned with the existing repo while avoiding broad refactors, speculative abstractions, and missing validation.

## Use When

- The user asks to implement, add, change, fix, wire, update, or remove scoped behavior.
- Requirements are clear enough to start from existing code and repo patterns.
- The task may touch code, config, docs, tests, or generated artifacts.

## Avoid When

- The failure is not understood or reproducible; use `diagnose`.
- Tests, examples, contracts, or executable evidence are the main product or driver of the change; use `verification-driven-change`.
- The user asks for architecture, requirements, or tradeoff clarification; use `planning-grill`.
- The user asks to capture or sync accepted decisions, roadmap/status changes, priority notes, ADRs,
  capture-worthy handoff notes, workflow lessons, or repo context docs as a long-lived source of
  truth; use `project-lifecycle`.
- The user asks to review current changes; use `code-review`.
- The user asks to commit, push, branch, or open a PR; use `conventional-git-flow`.
- The task is to create or update a skill; use `skill-creator`.

## Workflow

1. Inspect the relevant entry points, nearby patterns, tests, configs, and repo rules before editing.
2. State or infer the smallest scoped behavior change; ask only if ambiguity materially affects correctness.
3. Apply minimal edits that fit existing ownership boundaries and style. Do not rename public APIs, add dependencies, or refactor broadly unless required.
4. Check cross-cutting risks while implementing:
   - API contract or response shape impact: preserve stable status codes, machine-readable error codes, request IDs, and documented response shape.
   - Security, auth, secrets, PII, or logging exposure: never hardcode or log secrets, tokens, credentials, PII, or sensitive request/response bodies.
   - Auth, authorization, attacker-controlled input, or crypto impact: preserve existing auth/authz boundaries, validate untrusted input, and use existing vetted auth/crypto mechanisms instead of inventing new ones.
   - Testing or verification depth needed for the risk level; for test-first work, use one behavior slice at a public interface and existing test conventions.
   - Runtime/config/deploy impact when env vars, ports, services, external dependencies, or runtime requirements change.
5. Before finalizing the approach, load the corresponding `policy-*` skill when that policy area materially affects correctness, safety, or acceptance. This is a hard gate for API contracts, security/auth/secrets/PII, verification strategy, runtime/deploy/config behavior, and language/framework ownership boundaries. Do not load unrelated policy skills for routine edits where the built-in guardrails fully cover the risk.
6. If implementation reveals a product/architecture decision, rejected approach, deferred scope,
   scope-changing pivot, status or documentation drift, capture-worthy handoff note, loop active
   state, discussion record, or reusable workflow lesson that future work would need to remember, emit a
   `project-lifecycle` capture candidate before continuing or in the final summary. If none of
   those signals appear, report `lifecycle_capture_candidate: none` rather than asking for lifecycle
   capture on every ordinary completion.
7. Validate with the narrowest useful evidence allowed by the repo policy: focused tests, type checks, smoke checks, review, or manual checklist.
8. Summarize changed behavior, files touched, verification, assumptions, lifecycle capture candidates, and residual risk.

## Tool And Side-Effect Boundaries

- Do not run programs, tests, formatters, or installers unless the user requested execution or the active repo policy requires it.
- Do not commit, stage, push, deploy, migrate, or perform destructive actions unless explicitly requested.
- Preserve unrelated user changes in the worktree.
- Do not expose internal errors, stack traces, SQL, hostnames, paths, or sensitive details in user-facing responses.
- Do not invent new test frameworks, fixture systems, dependencies, public API names, or business logic unless explicitly approved.
- If infra-impact triggers are present but infra updates are skipped, state why infra changes are not required.
- This skill carries the minimum implementation guardrails for routine edits, but the policy gate above is mandatory when a policy area materially affects correctness, safety, or acceptance.

## Output

Return:

- `summary`: the behavior or artifact changed.
- `files_touched`: exact paths.
- `verification`: commands run, skipped checks, or manual checklist.
- `assumptions`: correctness-relevant assumptions.
- `lifecycle_capture_candidate`: signal-driven `project-lifecycle` candidate, such as an
  implementation pivot, accepted decision, deferred scope, status/doc drift, capture-worthy handoff
  note, loop active state, discussion record, reusable workflow lesson, or `none`.
- `residual_risks`: remaining risks or follow-up work.

## Version History

- v0.1.0 (2026-05-18): Initial scoped implementation workflow.
- v0.1.1 (2026-05-18): Add audited minimal API, security, testing, and infra guardrails.
- v0.1.2 (2026-05-19): Route test-primary and executable-evidence changes to `verification-driven-change`.
- v0.1.3 (2026-05-21): Route code review requests to renamed `code-review`.
- v0.1.4 (2026-05-22): Route long-lived decision/status document synchronization to `sync-decision-docs`.
- v0.1.5 (2026-05-28): Emit lifecycle capture candidates for implementation pivots and project-level decisions.
- v0.1.6 (2026-05-29): Add minimal auth, input-validation, and crypto implementation guardrails.
- v0.1.7 (2026-06-11): Add signal-driven lifecycle candidate checks at ordinary task completion.
- v0.1.8 (2026-06-18): Add hard policy gate for material API, security, testing, infra, and language/framework risks.
