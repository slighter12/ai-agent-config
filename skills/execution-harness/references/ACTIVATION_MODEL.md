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

Once applied, the assistant must act as orchestrator for the response:

- State whether the current response is a harness plan or a single-agent fallback.
- Keep every selected phase visible, even when the main agent owns it locally.
- Give each phase one owner, bounded scope, expected evidence, and acceptance gate.
- Do not perform implementation work in the same response unless the user explicitly requested execution after the harness plan.
- If delegated agents are proposed or used, include `agent_selection` with runtime-resolved model and cost disclosure.

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
- Guessed model, cost, or tool availability.
- Dropping a phase owner after avoiding a costly or cost-opaque role.
- Overlapping delegated write scopes without an ordering gate.
