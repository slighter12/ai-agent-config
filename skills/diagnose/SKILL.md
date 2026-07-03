---
name: diagnose
description: Diagnose bugs, failing behavior, flaky tests, and performance regressions through a feedback-loop-first workflow. Use when the user reports something broken, failing, flaky, slow, throwing, regressing, or asks to debug/diagnose/root-cause it. Avoid when the user only wants implementation of an already-understood change, current-diff review, git workflow, or broad design planning.
license: MIT
compatibility: [codex, claude, gemini]
metadata:
  version: "0.1.7"
---

# Diagnose

## Purpose

Keep debugging work anchored to observable evidence. The skill prevents drift by building or identifying a feedback loop before changing code.

## Use When

- The user reports broken, failing, flaky, slow, throwing, or regressed behavior.
- The task asks for diagnosis, debugging, root cause, reproduction, or performance investigation.
- A proposed fix is risky because the failure has not been reproduced or measured.

## Avoid When

- The user asks for a normal implementation with clear requirements; use `implement-change`.
- The user asks only to review current git changes; use `code-review`.
- The user asks only for commit, branch, push, or PR work; use `conventional-git-flow`.
- The main task is architecture or requirements clarification; use `planning-grill`.

## Workflow

1. Define the reported symptom and the expected behavior in concrete terms.
2. Build or identify the fastest credible feedback loop: failing test, command, script, HTTP request, browser flow, fixture replay, profile, or measurement. A credible loop is red-capable for the user's actual symptom, repeatable or high-reproduction, fast enough to iterate, and agent-runnable. Do not guess at a fix when no credible loop or missing artifact can confirm the symptom.
3. Reproduce the failure or raise the reproduction rate enough to investigate. If this is impossible, state what artifact or access is missing.
4. For flaky or nondeterministic failures, improve reproduction rate and isolate time, concurrency, filesystem, network, random seed, or environment variables.
5. Generate 3-5 falsifiable hypotheses and test one variable at a time.
6. Instrument only where it distinguishes hypotheses; tag temporary debug output with a searchable prefix and remove it before finishing.
7. Before finalizing the cause, fix direction, or verification plan, load the corresponding `policy-*` skill when that policy area materially affects correctness, safety, or acceptance. This is a hard gate for API contracts, security/auth/secrets/PII, verification strategy, runtime/deploy/config behavior, and language/framework ownership boundaries. Do not load unrelated policy skills for routine diagnosis where the built-in guardrails fully cover the risk.
8. Fix only after the cause is supported by evidence.
9. Add or describe regression coverage at the right public interface, then rerun the original feedback loop.
10. Clean up temporary probes and report the root cause, fix, verification, residual risk, and any architecture follow-up.

## Tool And Side-Effect Boundaries

- Prefer read-only inspection until a feedback loop is identified.
- Targeted reproduction, narrow tests, focused commands, fixture replays, browser flows, and measurements are allowed when they are the smallest credible feedback loop for the reported symptom.
- Do not run broad test suites or long commands unless the user requested execution or the active repo policy requires it.
- Do not commit, push, deploy, migrate, or perform destructive actions.
- Do not hardcode or log secrets, tokens, credentials, PII, or sensitive request/response bodies while instrumenting.
- For security-sensitive symptoms, include log/response exposure and trust-boundary hypotheses without printing sensitive samples.
- Do not invent business logic or expected behavior when the symptom or requirement is unclear; ask for the missing constraint or artifact.
- If a correct regression-test seam does not exist, state that as a finding instead of adding a weak or misleading test.
- Load `policy-testing` when the diagnosis depends on verification depth, regression evidence, risk level, or validation gate selection.
- If no red-capable loop can be built, ask for the missing artifact, access, logs, HAR, trace, fixture, or permission for targeted instrumentation instead of continuing with a speculative fix.

## Output

Return:

- `symptom`: the concrete observed failure.
- `feedback_loop`: the chosen reproduction or measurement path.
- `hypotheses`: ranked falsifiable hypotheses, when diagnosis is still in progress.
- `root_cause`: evidence-backed cause, when known.
- `fix_summary`: what changed or should change.
- `verification`: commands, tests, measurements, or manual checks.
- `residual_risks`: what may still be wrong.

## Version History

- v0.1.0 (2026-05-18): Initial feedback-loop-first diagnosis workflow.
- v0.1.1 (2026-05-18): Add audited minimal flaky-failure, instrumentation, and ambiguity guardrails.
- v0.1.2 (2026-05-19): Tighten feedback-loop, temporary instrumentation, and regression seam guardrails.
- v0.1.3 (2026-05-21): Route code review requests to renamed `code-review`.
- v0.1.4 (2026-05-29): Add security-sensitive diagnosis guardrail for exposure and trust-boundary hypotheses.
- v0.1.5 (2026-06-18): Add hard policy gate for material API, security, testing, infra, and language/framework diagnosis risks.
- v0.1.6 (2026-06-29): Define red-capable feedback loop criteria and stop condition for missing diagnostic evidence.
- v0.1.7 (2026-07-03): Clarify that targeted reproduction is allowed as the diagnostic feedback loop.
