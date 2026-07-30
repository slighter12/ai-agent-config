---
name: code-review
description: Review code or current changes for correctness, regression, security, and spec alignment without editing. Use when the user asks for code review, change review, release readiness, or a security-focused review. Avoid when the user asks to implement fixes rather than assess them.
metadata:
  invocation: model
---

# Code Review

Remain read-only. Before reviewing:

1. Pin the comparison point: explicit range, merge base, staged diff, unstaged diff, or named files.
2. Find the controlling spec or ticket.
3. Find repository standards and instructions.
4. Inspect the diff in its runtime context.

Prioritize concrete defects, regressions, missing validation, compatibility failures, and violated acceptance criteria. Each finding needs severity, location, evidence, impact, and a feasible correction. Do not report style preferences as defects.

For an ordinary review, review directly. Only an explicit request for a full review authorizes launching at most two independent reviewers; reconcile their evidence before answering.

For auth, tokens, secrets, PII, or log exposure, also load [the security checklist](references/SECURITY_CHECKLIST.md).

For release readiness, separate verified evidence from unverified conditions. Check acceptance coverage, required build and test results, compatibility, migrations and configuration, rollout and rollback, observability and operator documentation, and unresolved blockers. Missing required evidence prevents a release-ready conclusion even when no defect was found.

Return findings first, ordered by severity, then open questions and residual risks. If there are no findings, say so and state what was not verified.
