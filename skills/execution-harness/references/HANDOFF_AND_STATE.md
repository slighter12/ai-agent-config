# Handoff And State

Use lightweight state tracking. Avoid formal ceremony unless it prevents confusion.

## Handoff Fields

Include:

- `owner`: role or skill responsible for the next step.
- `goal`: bounded outcome.
- `inputs`: paths, diffs, links, commands, or constraints needed.
- `acceptance`: evidence that proves the step is done.
- `blockers`: unresolved blockers only.
- `suggested_skills`: narrow skills the next owner should consider, or `none`.
- `redaction_notes`: sensitive data removed or still requiring care.
- `brief`: concise context the next owner needs before acting.
- `return`: expected output shape.

## Agent Selection Fields

When delegation is proposed or used, include `agent_selection`. Resolve concrete agents at runtime from tool metadata, local configuration, or explicit user instruction. Do not maintain a fixed roster in this skill.

Always name the orchestration owner:

- `orchestration_owner`: `none`, `main_agent`, or the selected runtime orchestrator role.
- `orchestration_reason`: why coordination stays local or is delegated.

When `orchestration_owner` is a runtime role, include it in `delegated_agents` with purpose, scope, model, cost, and spawn timing disclosure.

For each selected delegated agent, include:

- `agent_id_or_role`: selected runtime agent or role name.
- `resolved_model_or_inheritance`: concrete model, `inherits parent`, or `unknown`.
- `model_source`: `tool_metadata`, `local_config`, `user_instruction`, `inherited`, or `unknown`.
- `cost_note`: concise cost or cost-opacity note.
- `purpose`: phase outcome this agent owns.
- `scope`: allowed files, commands, or read/write boundary.
- `spawn_timing`: when to delegate, including ordering gates.
- `selection_reason`: why this role is the lowest-cost capable fit.

If no delegated agents are selected, still write `agent_selection`; use `orchestration_owner: main_agent` for a local harness plan or `orchestration_owner: none` for a single-agent fallback.

Example:

```yaml
agent_selection:
  orchestration_owner: main_agent
  orchestration_reason: two bounded specialist lanes do not need delegated scheduling
  delegated_agents:
    - agent_id_or_role: runtime-selected-implementation-role
      resolved_model_or_inheritance: local-configured-model
      model_source: local_config
      cost_note: lower-cost implementation role than parent
      purpose: implement one bounded code slice
      scope: src/example/*
      spawn_timing: after interface contract is fixed
      selection_reason: narrow write scope and implementation-focused role
```

## State Tracking

Keep one short source of truth for:

- Current phase.
- Active owner.
- Assumptions.
- Open blockers.
- Verification status.
- Next action.

For long-running, repeated, or multi-session phases, also track:

- `work_unit`: smallest verifiable unit.
- `evaluator_contract`: evidence that proves the phase gate and where it is surfaced.
- `verified_progress`: accepted progress, rejected evidence, and evidence pointers.
- `stop_or_escalate`: retry, turn, blocker, tool-failure, budget, and human-review stops.

This state guides the next action; it does not prove final completion by itself.

## Lifecycle Capture Candidates

At completion or a phase boundary, create a `project-lifecycle` handoff when there is an accepted
decision, implementation pivot, status or documentation drift, capture-worthy handoff note, loop
active state, discussion record, or reusable workflow candidate worth classifying.

Use `project-lifecycle` and its `references/CAPTURE_GATE.md` target list to classify the candidate.
Shared skill updates require explicit user approval and `skill-creator` review.

## Phase Boundary Candidates

At phase completion, route each remaining gate to the narrowest owner. Use `project-lifecycle` only
for capture-worthy lifecycle signals, not as the silent owner for standalone review, verification,
commit readiness, or handoff packaging gates.

- Accepted decision or status docs sync -> `project-lifecycle`.
- Implementation pivot, deferred scope, or rejected approach -> `project-lifecycle`.
- Documentation drift -> `project-lifecycle`.
- Current diff or code risk review -> `code-review`.
- Missing verification evidence -> active task owner with `policy-testing` guidance.
- Commit readiness or PR flow -> `conventional-git-flow` when the user asks for git actions.
- Handoff packaging without a capture signal -> bounded handoff fields.
- Capture-worthy handoff note -> `project-lifecycle`.
- Loop active state or discussion record -> `project-lifecycle`.
- Reusable workflow or skill learning -> `project-lifecycle`.
- No action -> close the phase with stated evidence.
