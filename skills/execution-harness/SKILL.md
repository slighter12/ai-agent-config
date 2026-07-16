---
name: execution-harness
description: Coordinate optional orchestration across phases, agents, git/workspace state, verification gates, diff sanity gates, lifecycle gates, and capture candidates. Use when the user asks for a harness, multi-agent delivery, explicit phase gates, long-running handoff, or cross-phase coordination. Avoid when the request is mainly implementation, diagnosis, design clarification, current-diff review, testing strategy, git workflow, or standalone lifecycle capture.
license: MIT
compatibility: [codex, claude, gemini]
metadata:
  version: "0.1.13"
---

# Execution Harness

## Purpose

Provide an optional execution envelope for complex agent work. This skill coordinates phases and owners without replacing task skills that own implementation, diagnosis, design clarification, review, testing strategy, git workflow, or lifecycle capture. Harness activation and runtime-orchestrator delegation are separate decisions.

For long-running, repeated, or multi-session work, this skill adopts the provider-neutral loop
contract vocabulary from `policy-core` at the phase level. It does not write `/goal` launch
contracts and does not depend on any provider-specific evaluator, memory, or context isolation.

When routing work to another skill, distinguish routing ownership from loaded availability. Do not
say a sibling skill is unavailable merely because it has not been loaded into the current prompt;
name it as the owner unless the runtime or skill registry explicitly reports it missing.

## Use When

- The user explicitly asks for a harness, structured workflow, phase plan, delivery envelope, or multi-agent execution model.
- A task spans multiple phases, agents, sessions, or handoffs.
- Git/workspace state, verification gates, diff sanity, or recovery points need explicit coordination.
- The task needs a clear owner map and phase-by-phase evidence before continuing.
- A phase boundary needs explicit gates across review, verification, commit readiness, or handoff.
- A phase or final boundary needs signal-driven lifecycle routing for accepted decisions,
  implementation pivots, status or documentation drift, capture-worthy handoff notes, loop active state,
  discussion records, or reusable workflow lessons.
- A completed harness-managed task may reveal a lifecycle capture candidate worth handing off to `project-lifecycle`.

## Avoid When

- A simple task can be completed directly with no extra coordination.
- The request is only to implement code or config; use `implement-change` or the relevant implementation owner.
- The request is to diagnose a bug, flaky behavior, or performance regression; use `diagnose`.
- The request is to clarify design, architecture, requirements, or tradeoffs; use `planning-grill`.
- The request is only to review code or current changes; use `code-review`. Use this skill only when multi-agent review, explicit gates, or orchestration are requested or approved.
- The request is only to commit, push, or open a PR; use `conventional-git-flow`.
- The request is only to decide testing scope; use `policy-testing`.
- The request is only to create or update a skill; use `skill-creator`.
- The request is only to classify lifecycle signals from a completed phase, sync docs, record a
  capture-worthy handoff note, or capture workflow learning without broader orchestration; use
  `project-lifecycle`.

## Activation Model

- Prefer `orchestrator-suggested / user-approved` activation.
- Suggest this skill when a task needs explicit coordination across phases, agents, git/workspace state, verification gates, or capture.
- Apply it directly only when the user explicitly asks for harness-style workflow or has already approved using it for the active task.
- Do not self-apply it for ordinary implementation, review, testing, or git workflows.

## Mandatory Behavior When Activated

- For an active harness, choose `main_agent` or a runtime orchestrator role. Use `none` only when reporting a non-activated single-agent fallback.
- Keep orchestration local by default; a request for a harness does not itself delegate a runtime orchestrator.
- State that the response is an active harness plan.
- Do not perform implementation work in the same response unless the user explicitly requests execution after the harness plan.
- Keep every selected phase visible, including phases owned by the main agent instead of a delegated agent.
- Every phase must name an owner, scope, expected evidence, and acceptance gate.
- Long-running or repeated phases must name `work_unit`, `evaluator_contract`, `verified_progress`, `claim_boundary`, `stop_or_escalate`, `metric_type`, and `context_layers`; `policy-core` owns their detailed vocabulary.
- Every delegated owner must appear in both `phase_plan` and `agent_selection`.

## Workflow

1. Frame the objective, constraints, done criteria, and known blockers.
2. Select only the phases and gates that reduce risk for this task.
3. Choose the orchestration owner:
   - `none` only when declining the harness in favor of direct single-agent work.
   - `main_agent` for an active harness with small or medium coordination, including one or two specialist lanes.
   - Runtime orchestrator role when available and the user explicitly asks for a team or orchestrator, or the task needs three or more independent specialist lanes, sustained cross-session scheduling, or arbitration between specialist outputs.
4. Route each gate to the narrowest owner skill or agent role; if orchestration was delegated, route the plan through that orchestrator before specialist handoffs.
5. Resolve delegated agent names and model facts at runtime from tool metadata, local configuration, or explicit user instruction.
6. Keep one lightweight source of truth for assumptions, status, owner, and next action.
7. For any long-running, repeated, or multi-session phase, define every loop-contract field:
   - `work_unit`: the smallest verifiable unit for that phase.
   - `evaluator_contract`: what evidence proves the phase gate and where it must be surfaced.
   - `verified_progress`: where accepted and rejected progress evidence pointers are tracked.
   - `claim_boundary`
   - `stop_or_escalate`: retry, turn, blocker, tool-failure, budget, and human-review stops.
   - `metric_type`
   - `context_layers`
   Use deterministic validators or frozen metrics when available. If no objective or calibrated standard exists, route the phase to a human-review packet instead of self-certifying completion.
8. Pass only the relevant guidance into handoffs; do not assume subagents can dynamically gain skills or MCPs.
9. Coordinate the verification gate selected by the active repo policy or `policy-testing`; do not choose verification depth here.
10. At phase or final completion, detect whether a concrete lifecycle signal is present. If it is, hand the facts to `project-lifecycle` for classification; otherwise continue the remaining gate or close the phase.
11. Keep reusable learning as a `project-lifecycle` capture candidate; otherwise set `capture_candidate` to `none`.

## Agent Selection Disclosure

- Do not keep a fixed list of agent names in this skill.
- Resolve available agents at runtime from tool metadata, local configuration, or explicit user instruction.
- Select a runtime orchestrator role only when the selected runtime exposes one; otherwise keep orchestration with `main_agent`.
- When `orchestration_owner` is a runtime role, disclose it as a delegated agent with model, cost, purpose, scope, and spawn timing.
- Prefer the lowest-cost capable delegated role for each phase.
- Treat inherited-model or unknown-cost agents as cost-opaque and disclose that before selecting them.
- If model or cost cannot be resolved, write `unknown` or `inherits parent`; do not guess.

## Tool And Side-Effect Boundaries

- This skill coordinates capability routing; it does not perform runtime capability injection.
- Do not assume an agent can use a skill or MCP that its role config disables.
- Do not assume provider memory, a goal evaluator, hidden context isolation, or unsurfaced tool
  evidence. Phase gates must rely on surfaced evidence and authoritative current state.
- Do not generate Goal Briefs or `/goal` launch contexts; route that work to `goal-context` after
  roadmap and success criteria are settled.
- Git staging, committing, pushing, branch creation, PR creation, destructive commands, and external-service side effects still require explicit approval.
- Testing or command execution remains governed by the active repo policy and `policy-testing`.
- Lifecycle capture signals found during phase closeout or session closeout are owned by
  `project-lifecycle`; standalone phase gates such as review, verification, commit readiness, and
  handoff packaging remain owned by their narrowest gate owner.
- Shared skill updates remain governed by `skill-creator` and require explicit user approval.

## Anti-Patterns

- Do not discuss multi-agent execution without listing the actual selected runtime agents in `agent_selection`.
- Do not treat harness activation as automatic delegation to a runtime orchestrator role.
- Do not guess model, cost, or tool availability for a delegated agent.
- Do not drop a phase owner after pruning a costly or cost-opaque agent.
- Do not assign overlapping write scopes to multiple agents without an ordering gate.
- Do not treat progress ledgers, transcripts, or model confidence as completion evidence.
- Do not treat this skill as general planning prose after it has been explicitly activated.

## Output

Return:

- `objective`: active goal and done criteria.
- `phase_plan`: current phases, gates, owners, and expected evidence.
- For long-running or repeated phases, include `work_unit`, `evaluator_contract`, `verified_progress`, `claim_boundary`, `stop_or_escalate`, `metric_type`, and `context_layers` in the phase entry.
- `agent_selection`: orchestration owner plus selected runtime agents and model/cost disclosure, or `orchestration_owner: none` with the reason.
- `handoffs`: bounded owner assignments, if delegation is needed.
- `verification`: selected verification gates and owner roles.
- `git_workspace_notes`: relevant status, checkpoint, or commit-readiness notes.
- `lifecycle_capture`: `project-lifecycle` handoff summary, or `none`.
- `capture_candidate`: lifecycle capture candidate, or `none`.
- `next_action`: one concrete next step.

In `agent_selection`, include:

- `orchestration_owner`: `main_agent`, runtime orchestrator role, or `none`
- `orchestration_reason`

For each delegated agent, include:

- `agent_id_or_role`
- `resolved_model_or_inheritance`
- `model_source`: `tool_metadata`, `local_config`, `user_instruction`, `inherited`, or `unknown`
- `cost_note`
- `purpose`
- `scope`
- `spawn_timing`
- `selection_reason`

## Version History

- v0.1.0 (2026-05-13): Initial optional execution envelope for orchestrator-suggested structured workflow.
- v0.1.1 (2026-05-18): Reposition harness as orchestration only and route ordinary task work to task skills.
- v0.1.2 (2026-05-21): Route ordinary code review to `code-review` while preserving harness ownership of multi-agent gates.
- v0.1.3 (2026-05-22): Add then-separate phase closeout gate routing while keeping learning capture separate.
- v0.1.4 (2026-05-23): Add runtime agent selection disclosure and activated orchestrator behavior.
- v0.1.5 (2026-05-28): Route lifecycle gates and capture candidates to `project-lifecycle`.
- v0.1.6 (2026-06-11): Align harness lifecycle signals and targets with signal-driven lifecycle
  capture.
- v0.1.7 (2026-06-21): Adopt phase-level loop contract fields for long-running and repeated work
  while keeping `/goal` authoring with `goal-context` and lifecycle capture with
  `project-lifecycle`.
- v0.1.8 (2026-06-21): Clarify sibling skill routing ownership versus loaded availability, so
  harness does not describe `goal-context` as unavailable when it is only not loaded.
- v0.1.9 (2026-06-29): Add suggested skill, redaction, and brief fields to bounded handoff reference guidance.
- v0.1.10 (2026-06-30): Remove retired lifecycle skill id from history wording.
- v0.1.11 (2026-07-03): Replace fixed agent-name references with runtime-resolved role capability
  wording.
- v0.1.12 (2026-07-03): Add explicit local-versus-delegated orchestration owner selection.
- v0.1.13 (2026-07-16): Separate harness activation from runtime delegation, centralize phase routing, and retain standalone orchestration and all-field handoff redaction guardrails.

## References

- `references/INDEX.md` - Navigation for harness references.
- `references/HARNESS_OVERVIEW.md`
- `references/ACTIVATION_MODEL.md`
- `references/PHASES_AND_GATES.md`
- `references/GIT_AND_WORKSPACE.md`
- `references/VERIFICATION_GATES.md`
- `references/HANDOFF_AND_STATE.md`
