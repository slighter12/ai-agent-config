---
name: policy-core
description: "Apply cross-policy reference guidance for response structure, language rules, execution constraints, scope limits, and shared vocabulary. Use when cross-cutting policy alignment or precedence across multiple policies is the main question, or when multiple policy skills govern the same decision point and ownership/precedence is unclear. Avoid when a task skill such as `implement-change`, `diagnose`, `planning-grill`, or `code-review` can own the workflow and policy boundaries are clear."
metadata:
  version: "0.1.6"
---

# Policy Guide

Use the referenced policy files for rules. CORE.md is the source of truth when loaded, and the repo-root AGENTS.md provides a minimal baseline when CORE is not applied. Keep heavy workflow coordination in `execution-harness`, not in core.

Trigger `policy-core` automatically when the same decision point is governed by multiple policies, ownership or precedence is unresolved, and choosing the wrong owner could affect safety, API contracts, verification confidence, runtime behavior, or repo-wide baseline rules. Do not trigger it for ordinary policy overlap when local policy wording already names the owner or handoff.

When naming workflow owners, distinguish ownership from loaded availability. Do not say a sibling
skill is unavailable merely because it has not been loaded into the current prompt; name it as the
routing owner unless the runtime or skill registry explicitly reports it missing.

Loop contract ownership summary:

- `planning-grill` owns unresolved roadmap, product direction, architecture tradeoffs, and success
  criteria before a goal or harness contract is ready.
- `goal-context` owns manual Goal Brief and `/goal` launch contract authoring. It does not own
  durable memory or active execution.
- `execution-harness` owns phase, owner, gate, evidence, and handoff coordination for long-running
  or repeated work. It does not write Goal Briefs or `/goal` launch contexts.
- `project-lifecycle` owns active-state classification, durable capture, discussion records, and
  workflow lessons.
- `policy-core` owns shared vocabulary and precedence only, not the artifacts above.

## Version History

- v0.1.0 (2026-05-08): Initial portable policy release with the converged cross-policy core boundary matrix, precedence, and global response/execution constraints.
- v0.1.1 (2026-05-13): Add shared workflow vocabulary and execution-harness boundary guidance.
- v0.1.2 (2026-05-18): Clarify core as cross-policy reference guidance, not a primary task workflow.
- v0.1.3 (2026-05-21): Update task workflow boundary references for `code-review`.
- v0.1.4 (2026-05-29): Define detected precedence/conflict auto-trigger without making core a general policy dependency.
- v0.1.5 (2026-06-21): Add provider-neutral loop contract vocabulary for goal, harness, and
  lifecycle alignment without defining provider runtime behavior.
- v0.1.6 (2026-06-21): Add self-contained loop ownership summary and loaded-availability boundary
  to prevent sibling routing owners from being described as unavailable.

## References

- `references/INDEX.md` - Use for navigation and file selection.
- `references/CORE.md`
- `references/POLICY_BOUNDARIES.md` - Use when multiple policy skills govern the same decision point and ownership or precedence is unclear.
- `references/LOOP_CONTRACT.md` - Use for long-running goal, repeated loop, surfaced evidence,
  claim boundary, metric type, and active-state vocabulary.
