---
name: code-review
description: "Review code or current git changes through sanity, project-profile-aware full, security, architecture-diff-risk, or release-readiness modes without using provider-native review mode. Use when the user asks to review current changes, current diff, staged or unstaged changes, a code path, full/thorough/firmware code review, security review, auth/token/secret/PII/log leak handling, deploy readiness, go-live blockers, release risk, or whether changes are reasonable, coherent, safe to keep, or ready for follow-up. Ask one routing question when the user only says review this and the target could be code, design, git metadata, release readiness, or session/process review. Avoid when the task is design or architecture direction, whole-codebase architecture discovery, implementation fixes, commit/branch/PR metadata, session retrospective, external issue triage, or non-code artifact review."
license: MIT
compatibility: [codex, claude, gemini]
metadata:
  version: "0.2.10"
---

# Code Review

## Purpose

Provide an independent code review workflow for current diffs, focused code paths, full reviews, security-sensitive reviews, and release-readiness checks. The goal is to produce evidence-backed findings while keeping fixes, git operations, deployment actions, planning, and provider-native review commands separate.

## Use When

- The user asks whether current changes are reasonable.
- The user says "review current changes" or similar.
- The user asks for a full, thorough, strict, or adversarial code review.
- The user asks for a security review or mentions auth, authz, tokens, secrets, credentials, PII, crypto, logging leaks, or trust boundaries.
- The user asks for a release audit, deploy-readiness review, go-live check, launch blockers, rollback risk, or whether a change is safe to ship.
- The user asks to review a current refactor diff for architecture risk.
- The user wants a fresh look at a diff from another session.
- The user wants to detect scope creep, accidental files, inconsistent docs, or missing validation.
- The user wants a project-profile-aware review, including application, library/SDK, infrastructure, firmware/embedded, or security-sensitive code.
- The task should ignore unrelated conversation assumptions and rely on repository evidence.
- The user provides a focus such as affected functionality, intended bug fix, or expected non-goal.
- `execution-harness` requests a bounded diff sanity gate.

## Avoid When

- The user wants implementation fixes rather than review findings.
- The user asks to commit, push, open a PR, or create release notes.
- The user asks to review commit messages, branch names, PR titles, or PR bodies; use `conventional-git-flow`.
- The user asks to review session lessons, workflow learnings, lifecycle capture, or
  capture-worthy handoff notes; use `project-lifecycle`.
- The user asks to review design direction, product direction, architecture options, or whole-codebase architecture shape; use `planning-grill`.
- The user asks for multi-agent review, explicit phase gates, or orchestration; use `execution-harness`.
- The request is non-code document, spreadsheet, or presentation review.

## Workflow

1. Select a review mode:
   - `sanity`: default for current diff, staged or unstaged changes, and "is this coherent/reasonable" prompts.
   - `full`: explicit full, thorough, strict, adversarial, or broad code review.
   - `security`: explicit security review or security-sensitive terms such as auth, token, secret, PII, crypto, logging leak, or trust boundary.
   - `architecture-diff-risk`: current diff or refactor diff architecture risk only.
   - `release-readiness`: explicit release audit, deploy readiness, go-live blocker, rollback risk, or safe-to-ship review.
2. If the prompt only says "review" and the target is unclear, ask one question to distinguish current diff, code path, design/architecture direction, git metadata, release readiness, or session/process review.
3. Infer the project profile from repo and prompt facts before setting review depth: application, library/SDK, infrastructure, firmware/embedded, security-sensitive, docs/skill repo, or mixed.
4. If the review appears high-impact, cross-layer, security-sensitive, firmware/safety-critical, too large for one context, or likely to need independent perspectives, ask whether to escalate through `execution-harness` or multi-agent review.
5. Establish a clean review frame using current git state, relevant diffs or code paths, nearby files, repo rules, and active policy references.
6. Before finalizing the verdict, load the corresponding `policy-*` skill when that policy area materially affects correctness, safety, or acceptance. This is a hard gate for API contracts, security/auth/secrets/PII, verification strategy, runtime/deploy/config behavior, and language/framework ownership boundaries. Do not load unrelated policy skills for routine review surfaces where the review frame fully covers the risk.
7. Inspect enough evidence for the selected mode and project profile; keep `sanity` bounded, but allow `full` and `security` to widen scope when the risk justifies it.
8. Report only material, evidence-backed findings. If no material issue is found, say so and list residual risks or skipped checks.
9. Do not modify files unless the user explicitly asks for fixes after the review.

## User-Supplied Focus

Treat focus text in the same user request as part of the review frame. Examples:

- "whether this affects xxx"
- "this is mainly to fix xxx"
- "check whether the change is still limited to xxx"

Use focus text to guide inspection, but do not treat it as proof. If the focus conflicts with the diff, report the mismatch as a finding or use `unclear intent`.

## Mode Guidance

- For `sanity`, use `references/REVIEW_FRAME.md` and keep inspection diff-centered.
- For `full`, use project-profile guidance from `references/REVIEW_FRAME.md` before deciding breadth. When an originating spec, issue, or documented standard is in scope, review both axes: whether the code follows repo standards and whether it satisfies the stated spec.
- For `security`, load `policy-security` as the disclosure and trust-boundary authority; load API, frontend, infra, language, or testing policy when that surface is materially in scope.
  At minimum, inspect auth/authz boundary changes, secrets/tokens/PII in code/logs/tests/responses, unsafe user-facing error disclosure, custom crypto, and validation of attacker-controlled input.
- For `architecture-diff-risk`, review only architecture risk visible in the current diff or targeted code path. Whole-codebase architecture discovery belongs to `planning-grill`.
- For `release-readiness`, use `references/RELEASE_READINESS.md`. Keep the review read-only and focused on release blockers, confirmation items, deployment order, rollback, migrations, config/env, queues, cache, assets, and public contracts.

## Optional Challenge Pass

Run a challenge-oriented pass when the selected mode is `full`, when the user explicitly asks for it, or when the project profile/risk warrants it. The challenge pass is still read-only and must not invoke provider-native review commands unless the user explicitly requests those commands.

In the challenge pass:

1. Challenge the inferred intent and scope.
2. Look for hidden regressions, contract changes, validation blind spots, and risky assumptions.
3. Keep only evidence-backed material findings.
4. If no issue is found, write `adversarial_findings: none`.

## Tool And Side-Effect Boundaries

- Prefer read-only inspection.
- Do not invoke provider-native review commands or slash review commands.
- Do not run tests or programs unless the active repo policy or user request requires execution.
- Do not deploy, tag, publish, run migrations, rotate secrets, clear caches, or change remote infrastructure.
- Do not stage, commit, push, or create PRs.
- Do not rely on hidden conversation context from another session as evidence.
- Do not silently spawn subagents; ask or route to `execution-harness` when escalation is needed and not already requested.

## Output

Return:

- `mode`: sanity, full, security, architecture-diff-risk, or release-readiness.
- `release_conclusion`: only for release-readiness mode; `BLOCKED`, `NEEDS_CONFIRMATION`, or `NO_BLOCKER_FOUND`.
- `reviewed_range`: only for release-readiness mode; selected PR, branch, commit range, tag range, or fallback range.
- `release_findings`: only for release-readiness mode; P0/P1/P2 release blockers and confirmation items.
- `deployment_order`: only for release-readiness mode when ordering matters.
- `unable_to_verify`: only for release-readiness mode; neutral limits or blocked confirmations.
- `project_profile`: inferred project type and review depth implications.
- `verdict`: reasonable, needs changes, or unclear intent.
- `inferred_intent`: what the diff appears to be doing.
- `findings`: material issues only, with file references.
- `scope_notes`: accidental or unrelated changes, if any.
- `validation_notes`: checks run, skipped, or recommended.
- `adversarial_findings`: only when an optional challenge pass was requested.
- `escalation_notes`: whether multi-agent or harness escalation was requested, skipped, or recommended.
- `residual_risks`: what may still be wrong after review.

## Version History

- v0.1.0 (2026-05-04): Initial portable diff-centered sanity review skill.
- v0.1.1 (2026-05-04): Add user focus handling and optional challenge pass rules.
- v0.1.2 (2026-05-04): Route ambiguous current-change review requests to bounded sanity review.
- v0.1.3 (2026-05-13): Add harness diff gate reference while preserving bounded read-only behavior.
- v0.2.0 (2026-05-21): Rename to code-review and add sanity, full, security, architecture-diff-risk, project-profile, and escalation routing.
- v0.2.1 (2026-05-21): Align output reference with the review field order.
- v0.2.2 (2026-05-29): Add self-contained minimum checklist for security review mode.
- v0.2.3 (2026-05-29): Align security review routing wording with optional-depth policy handoffs.
- v0.2.4 (2026-06-11): Route capture-worthy handoff notes to project lifecycle while keeping
  ordinary handoff packaging out of code review.
- v0.2.5 (2026-06-18): Add hard policy gate before review verdicts for material API, security, testing, infra, and language/framework risks.
- v0.2.6 (2026-06-29): Add standards/spec axis guidance for full reviews without default subagent escalation.
- v0.2.7 (2026-06-29): Add read-only release-readiness review mode.
- v0.2.8 (2026-06-29): Align release-readiness output contract across references.
- v0.2.9 (2026-07-03): Fix portable review evidence paths, staged diff inspection, output reference wording, and git-flow boundary guidance.
- v0.2.10 (2026-07-03): Generalize harness diff gate wording away from local runtime role names.

## References

- `references/INDEX.md` - Navigation for review frame and output rules.
- `references/RELEASE_READINESS.md` - Read-only release, deploy-readiness, and go-live risk review.
- `references/HARNESS_DIFF_GATE.md` - Optional execution-harness diff gate boundaries.
