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
- `test-runner` executes selected commands only.
- `reviewer` handles ordinary risk-first review gates.
- `oracle` handles explicit high-risk challenge gates.
- `policy-testing` selects the verification depth.

## Boundaries

- Do not run tests or programs unless the user or active policy allows execution.
- Do not let test-runner choose a broader strategy than the requested or selected command set.
- Do not let review gates replace command or smoke evidence when runtime behavior is the risk.
- Do not let command success replace review when the risk is contract, scope, or hidden regression.
