# Handoff And State

Use lightweight state tracking. Avoid formal ceremony unless it prevents confusion.

## Handoff Fields

Required:

- `owner`: role or skill responsible for the next step.
- `goal`: bounded outcome.
- `inputs`: paths, diffs, links, commands, or constraints needed.
- `acceptance`: evidence that proves the step is done.
- `blockers`: unresolved blockers only.
- `return`: expected output shape.

Add `suggested_skills` or `brief` only when they change the next owner's work. Before every handoff, inspect all fields for secrets, credentials, PII, private URLs, or raw request/response bodies. If any field may contain them, redact it before handoff and include `redaction_notes`; notes may name only the field and data type, never the raw value.

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

If no delegated agents are selected, still write `agent_selection`; use `orchestration_owner: main_agent` for an active local harness plan or `orchestration_owner: none` only for a non-activated single-agent fallback.

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

For long-running, repeated, or multi-session phases, carry `work_unit`, `evaluator_contract`, `verified_progress`, `claim_boundary`, `stop_or_escalate`, `metric_type`, and `context_layers` from the phase plan. Do not restate their semantics in the handoff.

This state guides the next action; it does not prove final completion by itself.

## Lifecycle Capture Candidates

At completion or a phase boundary, create a fact-only `project-lifecycle` handoff when a concrete capture signal is present. `project-lifecycle` owns target classification. Shared skill updates still require explicit user approval and `skill-creator` review.

## Phase Boundary Candidates

Use `PHASES_AND_GATES.md` as the canonical routing table before selecting a remaining gate. This reference only packages handoff state; use `project-lifecycle` only for a concrete capture signal.
