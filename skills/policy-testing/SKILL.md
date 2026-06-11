---
name: policy-testing
description: "Apply testing and verification policy: risk levels, unit/integration/e2e boundaries, determinism, fixtures, mocks, validation gates, and coverage confidence. Use when test strategy, verification depth, quality gates, or how to validate behavior is primary. Avoid when the task is mainly test-primary edits, implementation, or diagnosis; use `verification-driven-change`, `implement-change`, or `diagnose` and treat this as supporting policy."
metadata:
  version: "0.1.5"
---

# Policy Guide

Use the referenced policy file for full rules. Keep output aligned with these rules.
Use when risk classification, validation strategy, verification depth, or coverage confidence for a change is primary.
Use as supporting policy when the user asks for TDD, test-first implementation, red-green-refactor, integration tests, regression coverage, characterization tests, contract tests, property checks, or how to validate behavior before changing code; task flow remains owned by `verification-driven-change`, `implement-change`, or `diagnose` unless testing strategy itself is the task.
Also use when `execution-harness` needs a verification gate selection.
Avoid when the task is only language/framework execution defaults; use optional language-specific policy detail (for example `policy-go` or `policy-rust`) or optional `policy-api` detail only when implementation or contract behavior is being verified, while `policy-testing` owns the risk gate and verification scope.

## Version History

- v0.1.0 (2026-05-08): Initial portable policy release with converged testing ownership boundaries for risk gating, verification scope, and cross-policy boundary guidance.
- v0.1.1 (2026-05-13): Add reference-only verification gate and tool shipcheck guidance.
- v0.1.2 (2026-05-15): Add TDD and test-first routing language.
- v0.1.3 (2026-05-18): Reposition testing as supporting policy for implementation and diagnosis task flows.
- v0.1.4 (2026-05-19): Add `verification-driven-change` as the task owner for test-primary and executable-evidence changes.
- v0.1.5 (2026-05-29): Align routing wording with optional-depth policy handoffs.

## References
- `references/INDEX.md` - Use for navigation and file selection.
- `references/TESTING.md`
- `references/VALIDATION_GATES.md`
- `references/TOOL_SHIPCHECK.md`
