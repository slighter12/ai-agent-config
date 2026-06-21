# Capture Gate

Use this gate when a user or active task identifies a project decision, implementation pivot, phase
boundary, loop memory point, discussion record, capture-worthy handoff note, status change,
documentation drift, or workflow lesson that may need capture.

The gate is conservative: surface useful candidates early, but do not write anything without
approval.

## Capture Targets

Classify each candidate as one of:

- `no_capture`: one-off observation, weak evidence, unstable trigger, or already covered.
- `final_answer_note`: useful reminder for the current response, but not worth repo changes.
- `project_decision_doc`: accepted decision, rationale, tradeoff, ADR, domain note, or `CONTEXT.md`
  update.
- `project_status_doc`: roadmap, priority note, implementation status, deferred scope, or follow-up
  status update.
- `handoff_note`: capture-worthy context package for another session, role, phase, or human handoff.
  Plain handoff packaging without a long-lived capture signal belongs to the active phase owner or
  `execution-harness`, not lifecycle capture.
- `active_state_checkpoint`: short-lived loop or phase state needed for the next run, not a permanent
  source of truth.
- `discussion_record`: source-grounded summary of reviewed material, accepted decisions, rejected
  alternatives, open questions, and follow-up owners.
- `workflow_lesson`: reusable workflow, routing, or agent behavior lesson observed in this work.
- `project_local_docs`: project-specific rule or practice that belongs in README, AGENTS.md, or
  equivalent docs.
- `project_local_skill`: reusable only inside one project or workspace.
- `shared_skill_update`: cross-project guidance that fits the narrowest existing shared skill.
- `new_shared_skill`: distinct cross-project workflow with its own trigger, avoid cases, output
  contract, and side-effect boundaries.
- `script_or_helper`: deterministic behavior safer as code than prompt instructions.

## Decision And Status Docs

Propose long-lived docs capture when all are true:

- The decision or status change is accepted enough to persist, or clearly marked open/deferred.
- The target source of truth exists or is obvious from project conventions.
- Future work would be confused without the rationale, status, deferred scope, or domain language.
- The proposed wording does not invent owners, dates, priorities, product facts, or implementation
  completion.

Use ADR-style capture only when at least one is true:

- The decision is hard to reverse.
- There is a real tradeoff with rejected alternatives.
- The outcome would be surprising without recorded rationale.

Do not use decision docs as a scratchpad for active debate.

## Loop Memory And Discussion Records

Use an active state checkpoint when the next run needs continuity: current phase, owner, blockers,
verification status, and next action. Do not commit raw transcripts or full run logs by default.

Use a discussion record when source material or decisions would be hard to reconstruct later. A good
record names the sources, separates accepted decisions from rejected alternatives, and routes any
long-lived follow-up through the capture targets above.

For loop state, separate layers:

- Active state: current phase, owner, accepted progress, rejected evidence, blockers, verification
  status, evidence pointers, and next check.
- Durable capture: accepted decisions, status changes, pivots, documentation drift, handoff notes,
  or workflow lessons.
- Not captured by default: raw transcripts, full logs, provider memory assumptions, and prose-only
  confidence.

An active state checkpoint does not prove completion. Completion claims must be verified against
authoritative current state or routed to human review when no deterministic validator, frozen
metric, evidence-review protocol, or calibrated rubric exists.

## Workflow Learning And Skill Evolution

Only propose shared skill or agent behavior changes when all are true:

- The lesson is reusable across projects.
- The trigger and avoid cases are stable.
- The narrowest owner skill or role is clear, or a new skill is justified.
- Side effects, permissions, and provider compatibility are explicit.
- The candidate includes routing examples and manual validation.
- The user explicitly approves writing shared skill, agent, config, or installer files.

Before creating a new skill:

1. Inventory visible skills and their descriptions.
2. Prefer updating the narrowest existing owner.
3. Prefer references for long guidance.
4. Prefer scripts only for deterministic behavior that prevents repeated model reasoning.
5. Reject candidates that are project-specific, unstable, already covered, or not yet validated.

## Rejection Reasons

Use concise reasons:

- `too_project_specific`
- `no_stable_trigger`
- `already_covered`
- `belongs_in_project_docs`
- `not_enough_evidence`
- `side_effect_boundary_unclear`
- `provider_boundary_unclear`
- `active_decision_not_settled`
