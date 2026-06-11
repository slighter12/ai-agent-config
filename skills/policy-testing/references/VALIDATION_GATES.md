# Validation Gates

Use this reference when a task needs to choose the minimum evidence required before continuing.

This file does not change the default execution policy in `TESTING.md`: do not create tests, run tests, or run programs unless the user or active policy requires it.

## Gate Types

### Reasoning Gate

Use for low-risk docs, plans, copy, or simple config discussion.

Evidence:

- Logical verification.
- Manual checklist.
- Assumptions and residual risks.

### Manual Gate

Use when behavior depends on UI, environment, credentials, external services, or user-specific setup.

Evidence:

- Step-by-step manual verification.
- Expected observable results.
- Known limitations or environment assumptions.

### Command Gate

Use when a deterministic command is known and execution is allowed.

Evidence:

- Exact command.
- Working directory.
- Expected signal.
- Exit status and key output if executed.

`test-runner` may execute selected commands, but should not broaden the validation strategy.

### Review Gate

Use when regression, contract, scope, or adjacent behavior risk matters more than command output.

Evidence:

- Risk-first findings from reviewer or bounded current-diff review.
- Missing validation or residual risks.

### Challenge Gate

Use for high-risk security, data, concurrency, migration, architecture, or ambiguous high-impact decisions.

Evidence:

- Read-only oracle or adversarial pass.
- Concrete accept/reject criteria.
- Residual risks and rollback or mitigation notes when applicable.

### Smoke Gate

Use when runtime integration must be demonstrated, such as CLI, MCP, service startup, or E2E flows.

Evidence:

- Representative happy-path smoke.
- Representative error path.
- Dry-run boundary for mutating behavior when available.

## Selection Rules

- Low-risk changes usually need reasoning or manual gates.
- Medium-risk behavior changes usually need command, review, or manual gates.
- High-risk changes need explicit testing or challenge gates unless the user forbids execution.
- Tooling and agent-facing CLIs/MCPs should use `TOOL_SHIPCHECK.md`.
- If the required environment is unclear, stop and ask rather than inventing validation.
