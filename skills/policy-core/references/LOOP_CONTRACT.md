# Loop Contract

Use this reference when a long-running agent loop, goal mode, repeated phase, or cross-session
handoff needs a shared completion contract. This file defines provider-neutral standards, not a
provider runtime implementation.

## Core Rule

A loop is complete only when authoritative current-state evidence proves the accepted scope. The
agent must not redefine success around partial progress, prior intent, progress notes, or a narrow
check that does not cover the claim.

Every loop contract must define:

- `work_unit`: the smallest unit that can be verified without interpreting the whole project.
- `evaluator_contract`: what evidence the checker can see and how completion is judged.
- `verified_progress`: accepted work units and evidence pointers, not prose-only confidence.
- `claim_boundary`: accepted scope, non-claims, and rejected evidence.
- `stop_or_escalate`: retry, turn, budget, blocker, tool-failure, and human-review stop rules.
- `metric_type`: deterministic validator, frozen metric, evidence review, subjective rubric, or
  human review required.
- `context_layers`: conversation context, active state artifact, durable source artifact, lifecycle
  capture, and raw logs/transcripts.

## Metric Type

Prefer objective evidence when it exists:

- `deterministic_validator`: command, test, schema check, static inspection, or other replayable
  check with clear pass/fail semantics.
- `frozen_metric`: an uneditable metric or benchmark outside the agent's work boundary.
- `evidence_review`: a human or reviewer checks named artifacts against explicit criteria.
- `subjective_rubric`: quality or taste criteria with examples, negative examples, and calibration
  status.
- `human_review_required`: no objective validator or calibrated rubric exists, so the loop must
  produce a review packet and stop for human judgment.

LLM judgment is supporting evidence only unless the goal explicitly accepts LLM-as-judge. When no
quantified or calibrated standard exists, the agent must not mark the goal complete by self-review.

## Claim Boundary

Every long-running goal or repeated loop must state:

- `Accepted Scope`: what this stage is allowed to claim.
- `Non-Claims`: nearby stronger claims the work does not establish.
- `Rejected Evidence`: evidence that looks relevant but must not count, with the reason.

Example:

```md
Accepted Scope: scaffolded synthetic compositional recombination under the documented training gate.
Non-Claims: broad generalization, real-world transfer, or unconstrained compositional reasoning.
Rejected Evidence: previous "generalization training" wording and broad pass labels, because they
do not isolate the scaffolded synthetic recombination condition.
```

## Context Layers

Do not treat goal text or memory as proof.

- Conversation context: useful for locating work, but not durable evidence.
- Active state artifact: compact next-run continuity such as current phase, owner, accepted/rejected
  progress, blockers, and next check.
- Durable source artifact: Goal Brief, result document, checked-in docs, logs, or evidence excerpts
  that can be inspected later.
- Lifecycle capture: accepted project decisions, status, pivots, or reusable workflow lessons.
- Raw transcript/log: not committed by default; store only when explicitly required.

Progress ledgers and active state artifacts can guide the next turn, but final completion audit must
return to authoritative current state.

## Provider Notes

- Codex goal mode persists a thread goal, injects continuation steering, and relies on the model to
  call `update_goal` when complete or blocked. Contracts must guard against model self-marking by
  requiring objective evidence or human review.
- Claude `/goal` includes an evaluator flow, but evaluator-visible evidence is limited to surfaced
  conversation context. Contracts must require verification evidence to be reported in conversation.
- Other providers must not be assumed to have a goal evaluator, cross-session memory, isolation, or
  tool access unless the runtime explicitly provides it.

## Ownership

- `planning-grill`: unresolved roadmap, product direction, architecture tradeoff, or success
  criteria.
- `goal-context`: manual Goal Brief and launch contract authoring.
- `execution-harness`: phase, owner, gate, and evidence coordination for long-running or repeated
  work.
- `project-lifecycle`: active state classification, durable capture, discussion records, and
  workflow lessons.
- `policy-testing`: verification strategy and gate confidence.
