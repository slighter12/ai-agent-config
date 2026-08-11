# Conditional Full Review

Read this reference only after the main skill has selected an explicit optional mode. It extends the default `Standards` and `Spec` passes; it does not replace them.

- An explicit **correctness review** adds the `Correctness` axis.
- An explicit **full review** or **release readiness** adds both `Correctness` and `Release readiness`.
- `Security` remains conditional on the sensitive diff triggers in the main skill and uses `SECURITY_CHECKLIST.md`.

Use provider-native reviewers in parallel when their isolation or parallelism justifies the coordination cost. If that facility is unavailable, run the same passes sequentially in the current session. Preserve the selected axes and keep their evidence and findings separate.

## Correctness

Assess whether the changed behavior works in its runtime context and preserves supported behavior. Focus on concrete evidence rather than style preferences.

Review these criteria:

- **Defects** — incorrect control flow, state transitions, data handling, error handling, boundary behavior, concurrency, timing, resource cleanup, or failure recovery.
- **Regressions** — changed behavior that breaks existing callers, workflows, persisted data, CLI/API contracts, or previously supported happy and failure paths.
- **Missing validation** — absent or inadequate tests, assertions, input validation, error checks, or runtime safeguards where the diff and surrounding code show they are needed.
- **Compatibility** — public APIs, schemas, wire formats, configuration, versions, platforms, and deployment/runtime assumptions across supported producers and consumers; include forward and backward compatibility where relevant.

For every finding, provide severity, file and location, concrete evidence, impact, and a feasible correction. Distinguish an observed defect from a risk that remains unverified. Keep findings scoped to this axis; if the same hunk also violates Standards or Spec, report the independent reason under that axis rather than merging findings.

### Correctness pass prompt

> Run the Correctness pass on the selected diff and its runtime context. Check defects, regressions, missing validation, and compatibility across normal, boundary, failure, upgrade, and rollback-relevant paths. Trace affected callers, data, and configuration far enough to support each claim. Report each finding with severity, location, evidence, impact, and feasible correction; label unverified risks separately. End with the Correctness finding count and worst issue.

The Correctness pass is complete only after each changed behavior has been checked against its callers, state/data paths, failure paths, and supported compatibility surface, or the unverified gap is recorded as residual risk.

## Release readiness

Assess whether the change can be accepted and released based on evidence, not merely whether the diff appears defect-free. Examine the spec, repository instructions, deployment context, and available build/test/release artifacts.

Review every applicable criterion:

- **Acceptance** — each acceptance criterion maps to implemented behavior and verification evidence; record exclusions and open gaps.
- **Build** — required build, packaging, type-check, lint, or artifact-generation checks pass for the supported target(s).
- **Tests** — required unit, integration, end-to-end, smoke, upgrade, compatibility, and failure-path checks have results; distinguish not run from passed.
- **Migrations** — schema/data migrations are ordered, safe, repeatable as required, compatible with the deployment sequence, and have a data-preserving rollback or recovery plan.
- **Configuration** — defaults, environment variables, secrets, feature flags, permissions, manifests, and versioned settings are complete and compatible.
- **Rollout** — rollout sequencing, dependency ordering, canary/progressive strategy, capacity, availability, and operational ownership are defined.
- **Rollback** — code, configuration, and data rollback boundaries are explicit; irreversible changes have a recovery plan and the procedure is testable.
- **Observability** — health checks, logs, metrics, traces, dashboards, alerts, and relevant SLO/error signals can detect and diagnose failure.
- **Operator docs** — release notes, migration and rollback steps, runbooks, support contacts, and on-call actions are available and current.
- **Blockers** — unresolved release-blocking findings, acceptance gaps, failed checks, unsafe migrations, missing required evidence, or unowned operational risks are listed explicitly.

For each criterion, record the evidence source (command output, test/build artifact, document path, or verified runtime observation), the result, and any remaining uncertainty. Missing or unverified required evidence is a release blocker or residual risk; do not infer success from the absence of a reported defect. A Release readiness pass may conclude **ready** only when every applicable required criterion has concrete evidence and no blocker. Otherwise conclude **not ready** or **undetermined**, naming the missing evidence and blockers.

### Release readiness pass prompt

> Run the Release readiness pass on the selected diff, its originating spec, and the available project/release evidence. Check acceptance, build, tests, migrations, configuration, rollout, rollback, observability, operator documentation, and blockers. For every applicable criterion, cite concrete evidence and mark missing or unverified items. Separate verified evidence, open risks, and blockers. Conclude **ready** only if every required criterion is evidenced and no blocker remains; otherwise conclude **not ready** or **undetermined**. End with the Release readiness finding count and worst blocker.

The Release readiness pass is complete only after all applicable criteria are recorded with evidence or an explicit gap, and every blocker has an owner or a clear reason it prevents readiness.
