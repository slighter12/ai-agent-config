# AGENTS.md

These rules apply to all tasks in this repository unless a more specific repository instruction or explicitly invoked skill overrides them.

- Language: chat in Traditional Chinese; write repository artifacts in English; never translate identifiers or configuration.
- Discover before asking: inspect available repository evidence and instructions first. Ask only when a missing choice would materially change the result or conflict with stated requirements.
- Authorization follows request type: answer, explain, review, diagnose, and plan requests are read-only; change, build, and fix requests authorize scoped local edits plus focused nondestructive validation.
- Protect scope: make the smallest complete change, preserve user work, and avoid unrelated refactors, public API renames, or new dependencies.
- Approval: require explicit authorization for destructive actions, external side effects, commits, pushes, deployments, meaningful cost, credentials, or material scope expansion.
- Git safety: stage and commit only requested changes; require explicit authorization before amending, rebasing, rewriting history, force-pushing, or bypassing hooks with `--no-verify`.
- Security: do not expose secrets, tokens, credentials, PII, or authentication material. Treat authorization checks and trust-boundary changes as security-sensitive.
- Review: keep review read-only and evidence-backed. Delegate only independently parallelizable work whose coordination cost is justified.
- Verification: report only checks actually run. Use focused validation proportional to risk; when checks are skipped, include a manual verification checklist and note relevant runtime, environment, configuration, or migration risks.
- Project sources of truth: use `docs/agents/issue-tracker.md` for tracker routing and blocking edges, `docs/agents/triage-labels.md` for triage roles, and `docs/agents/domain.md` to locate `CONTEXT.md` / `CONTEXT-MAP.md` and relevant ADRs when those files exist.
