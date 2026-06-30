---
name: project-lifecycle
description: "Coordinate lifecycle capture for project decisions, implementation pivots, phase boundaries, capture-worthy handoff notes, loop memory, discussion records, status/documentation drift, and workflow learning. Use when the user asks to close out a phase, checkpoint progress, sync decision/status docs, record a decision, record a capture-worthy handoff note, wrap up a session, decide what a loop should remember, or capture reusable workflow/skill lessons; also use when active work reaches a project-level capture candidate. Avoid when the task is still active planning, direct implementation, diagnosis, current-diff review, test-primary work, pure git workflow, pure handoff packaging, pure skill authoring, or broad orchestration with no lifecycle capture need."
license: MIT
compatibility: [codex, claude, gemini]
metadata:
  version: "0.1.3"
---

# Project Lifecycle

## Purpose

Provide one lifecycle/capture gate for project memory. This skill classifies whether a completed
work segment, decision point, implementation pivot, phase boundary, loop run, discussion record,
capture-worthy handoff note, status change, or workflow lesson should be captured, where it belongs,
and what approval is needed before mutation.

For loop memory, this skill classifies state layers. It does not assume provider memory, hidden
context isolation, or automatic cross-session recall. Durable continuity should use repo artifacts,
Goal Briefs, handoff references, result docs, or explicit active-state checkpoints.

## Use When

- The user asks to record, sync, or persist an accepted project decision, product status, roadmap,
  priority note, deferred scope, ADR, `CONTEXT.md`, or other long-lived project source of truth.
- The user says a phase, stage, milestone, planned slice, checkpoint, or before-commit readiness
  point is done and asks what follow-up is needed.
- The user asks to wrap up a session, record a capture-worthy handoff note, compact context, or
  capture workflow, routing, or skill lessons after a concrete work segment.
- The user asks what a loop, automation, or recurring agent workflow should remember between runs.
- The user asks to manually record source material, discussion outcomes, accepted/rejected options,
  or rationale before future automation exists.
- An active `planning-grill` discussion reaches an accepted project-level decision that may matter
  later.
- An active `implement-change` task discovers an implementation pivot, scope change, or product or
  architecture decision that should be surfaced instead of silently disappearing into the final
  answer.
- `execution-harness` reaches a phase, final, capture-worthy handoff note, loop checkpoint, or
  capture gate.

## Avoid When

- Requirements, product direction, or architecture tradeoffs are still being decided; use
  `planning-grill` until there is a capture candidate.
- The user asks for direct implementation of a clear change; use `implement-change` and emit a
  lifecycle candidate only if a pivot or project-level decision emerges.
- The user reports a bug, failure, flake, or regression; use `diagnose`.
- The user asks to review current code changes; use `code-review`.
- Tests or executable evidence are the main product; use `verification-driven-change`.
- The user only asks to commit, branch, push, or open a PR; use `conventional-git-flow`.
- The user asks to create, revise, validate, or package a skill after capture has already been
  approved; use `skill-creator`.
- The task needs multi-agent or multi-phase orchestration before lifecycle capture; use
  `execution-harness`.
- The request is only to package a bounded handoff for another role or session without a
  capture-worthy lifecycle signal; use the active phase owner or `execution-harness`.

## Workflow

1. Identify the lifecycle event: `decision_point`, `implementation_pivot`, `phase_boundary`,
   `loop_run_checkpoint`, `loop_memory_capture`, `discussion_record`,
   `capture_worthy_handoff`, `workflow_learning`, `status_update`, or `documentation_drift`.
2. Confirm the evidence from the current conversation and local state. Do not invent decisions,
   priorities, dates, owners, product status, or implementation status.
3. For discussion records or loop memory, use `references/MANUAL_CAPTURE_PACKET.md` to separate
   source material, accepted decisions, rejected alternatives, active state, long-lived capture, and
   open questions.
4. Classify each candidate with `references/CAPTURE_GATE.md` as `no_capture`,
   `final_answer_note`, `project_decision_doc`, `project_status_doc`, `handoff_note`,
   `active_state_checkpoint`, `discussion_record`, `workflow_lesson`, `project_local_docs`,
   `project_local_skill`, `shared_skill_update`, `new_shared_skill`, or `script_or_helper`.
5. Read only the long-lived docs or skill references needed to validate the proposed capture.
   Prefer existing sources of truth over creating new structures.
6. Propose concrete updates or the next owner. Route approved skill authoring to `skill-creator`,
   git actions to `conventional-git-flow`, review to `code-review`, and ordinary docs edits to
   `implement-change`.
7. Ask for explicit approval before mutating long-lived docs, shared skills, agent files, shared
   config, git state, or external systems.
8. After approved edits, make the smallest reversible change and summarize what changed, what stayed
   open, and any remaining confirmation needed.

## Loop State Rules

- `active_state_checkpoint` is short-lived continuity for the next run: current phase, owner,
  accepted progress, rejected evidence, blockers, verification status, evidence pointers, and next
  check.
- Long-lived capture is for accepted decisions, status, pivots, documentation drift, reusable
  workflow lessons, and capture-worthy handoff notes.
- Raw transcripts, full logs, and verbose run history are not committed by default. Store focused
  excerpts or paths only when they are needed as evidence.
- Progress ledgers and active-state notes guide future work but do not prove completion. Completion
  claims must return to authoritative current state.
- If a loop has no deterministic validator, frozen metric, evidence-review protocol, or calibrated
  rubric, lifecycle capture should preserve the review packet and the human decision needed instead
  of recording the agent's self-certification as accepted.

## Decision Capture Rules

- Use ADR-style capture only for accepted decisions that are hard to reverse, carry a real tradeoff,
  set an architecture or product direction, or would surprise future work without rationale.
- Keep ADR-style capture minimal: decision, context, options considered, tradeoff, status, date, and
  follow-up owner when known.
- Do not turn active debate, temporary experiments, or stylistic preferences into long-lived
  decisions until the user accepts them.

## Tool And Side-Effect Boundaries

- Do not silently modify long-lived docs, shared skills, agent config, git state, or external
  systems. Default to capture candidates and proposed updates.
- Do not treat every final response as a lifecycle gate. Use this skill when there is an explicit
  wrap-up/closeout request or a concrete decision/pivot/lesson with reusable value.
- For ordinary task completion, surface lifecycle candidates only when there is a concrete signal:
  accepted decision, pivot, status or documentation drift, capture-worthy handoff note, loop active
  state, discussion record, or reusable workflow lesson.
- Do not turn active planning into accepted docs; unresolved items stay open or candidate-only.
- Do not use stale docs to override an explicit accepted decision from the active conversation; flag
  conflicts instead.
- Do not make `CONTEXT.md` a scratchpad. Record stable project context, resolved domain language,
  accepted constraints, and decisions that would otherwise be surprising.
- Do not turn active loop state into a permanent transcript ledger by default. Preserve only the
  summary needed for next-run continuity unless the user explicitly asks for an audit trail.
- Do not describe provider memory as required or available unless the current runtime explicitly
  provides it. Checked-in docs and explicit artifacts are the portable source of truth.
- Propose ADR capture only when the decision is hard to reverse, has a real tradeoff, or would be
  surprising without rationale.

## Output

Return:

- `summary`: lifecycle event and capture scope.
- `evidence`: conversation or local facts supporting the candidate.
- `capture_candidates`: candidates, classifications, and skip reasons.
- `recommended_owner`: the next skill or role, or `none`.
- `proposed_updates`: concrete file-level updates before mutation, or updates applied after approval.
- `manual_capture_packet`: source-grounded discussion or loop-memory packet, when relevant.
- `approval_needed`: exact edits or side effects needing approval, or `none`.
- `files_touched`: paths changed, if edits were approved.
- `manual_verification`: checklist for confirming lifecycle capture is correct.
- `residual_risks`: unresolved decisions, stale docs, source-of-truth conflicts, or weak evidence.

## Version History

- v0.1.3 (2026-06-30): Add minimal ADR-style capture thresholds without introducing a template system.
- v0.1.2 (2026-06-21): Clarify loop state layers, active-state evidence pointers, human-review
  capture, and no default reliance on provider memory or raw transcripts.
- v0.1.1 (2026-06-11): Add loop memory, discussion record, and manual capture packet guidance.
- v0.1.0 (2026-05-28): Initial lifecycle/capture gate consolidating phase closeout, decision-doc
  sync, handoff, status capture, and workflow learning triage.

## References

- `references/CAPTURE_GATE.md` - Capture classification, document thresholds, and skill evolution
  routing.
- `references/MANUAL_CAPTURE_PACKET.md` - Manual source-grounded capture packet for discussions,
  loop memory, and pre-automation workflow decisions.
