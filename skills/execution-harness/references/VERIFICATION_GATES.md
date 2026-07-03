# Verification Gates

The harness asks what evidence is needed before continuing. `policy-testing` owns the answer.

## Gate Types

- `reasoning`: logical verification and manual checklist only.
- `manual`: human-verifiable steps for UI, external services, or environment-dependent behavior.
- `command`: deterministic command with expected signal.
- `review`: risk-first review of changed behavior and nearby contracts.
- `challenge`: high-risk adversarial pass for security, data, architecture, concurrency, or migration risk.
- `smoke`: representative runtime check for service, CLI, MCP, script, or E2E behavior.

## Role Routing

- Main agent may perform reasoning and manual checklist drafting.
- A runtime command-execution role may execute selected commands only, when such a role exists.
- A runtime independent review role may handle ordinary risk-first review gates, when such a role
  exists.
- A runtime high-risk challenge role may handle explicit adversarial gates, when such a role exists.
- `policy-testing` selects the verification depth.

## Boundaries

- Do not run tests or programs unless the user or active policy allows execution.
- Do not let a command-execution role choose a broader strategy than the requested or selected
  command set.
- Do not let review gates replace command or smoke evidence when runtime behavior is the risk.
- Do not let command success replace review when the risk is contract, scope, or hidden regression.
