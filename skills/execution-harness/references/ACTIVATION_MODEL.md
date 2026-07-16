# Activation Model

Use an `orchestrator-suggested / user-approved` activation model.

## Suggest Trigger

The orchestrator may suggest this skill when the task needs explicit coordination across:

- Multiple phases.
- Multiple agents.
- Long-running or cross-session work.
- Git/workspace state.
- Verification gates.
- Diff sanity gates.
- Commit readiness.
- Learning capture.

## Apply Trigger

Apply the harness only when:

- The user explicitly asks for harness, structured workflow, multi-agent delivery, handoff, phase gates, or similar language.
- The user approves the orchestrator's suggestion to use the harness for the active task.
- The active task has already been scoped as a harness-managed workflow.

## Activated Mode

Activation and delegation are separate decisions. Once the harness is active, the assistant must choose who owns orchestration:

- Use `main_agent` for small or medium coordination, including one or two specialist lanes.
- Use a runtime orchestrator role only when available and coordination itself is the main work.
- Use `none` only when reporting that the harness was not activated and work stays direct.
- State that the response is an active harness plan.
- Keep every selected phase visible, even when the main agent owns it locally.
- Give each phase one owner, bounded scope, expected evidence, and acceptance gate.
- Do not perform implementation work in the same response unless the user explicitly requested execution after the harness plan.
- If delegated agents are proposed or used, include `agent_selection` with runtime-resolved model and cost disclosure.

## Delegated Orchestrator Trigger

Delegate orchestration to a runtime orchestrator role only when at least one applies:

- Three or more independent specialist lanes are needed.
- Sustained cross-session scheduling or long-running state tracking needs a dedicated coordinator.
- Specialist outputs need arbitration before execution can continue.
- The user explicitly asks for a runtime orchestrator or team.

A request for a harness, phase plan, or multi-agent delivery authorizes harness coordination; it does not by itself delegate a runtime orchestrator.

If no runtime orchestrator role is available, keep orchestration with `main_agent` and disclose the fallback.

## Agent Selection Disclosure

Do not keep a fixed list of agent names in the skill. Resolve selected agents at runtime from:

- Tool metadata.
- Local role configuration.
- Explicit user instruction.

If model or cost cannot be resolved, write `unknown` or `inherits parent`; do not guess. Treat inherited or unknown model cost as cost-opaque and disclose that before selecting the agent.

## Avoid Automatic Takeover

Do not apply this skill automatically for:

- Ordinary implementation.
- Ordinary review.
- Ordinary testing.
- Ordinary commit or PR preparation.
- Small docs, copy, config, or single-file changes.

When uncertain, suggest the harness and wait for approval instead of silently applying it.

## Anti-Patterns

- Hidden delegation or multi-agent discussion without `agent_selection`.
- Treating harness activation as automatic delegation to a runtime orchestrator role.
- Guessed model, cost, or tool availability.
- Dropping a phase owner after avoiding a costly or cost-opaque role.
- Overlapping delegated write scopes without an ordering gate.
