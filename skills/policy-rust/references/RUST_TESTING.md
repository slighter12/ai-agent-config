---
globs: ["Cargo.toml", "**/Cargo.toml", "*.rs", "**/*.rs"]
description: "Rust testing and execution policy - default test runs and tooling expectations"
---

# RUST_TESTING.md - Rust Testing and Execution

This file defines testing and execution rules for Rust projects.
Violating these rules is incorrect output.

---

## HARD RULES

### 1) Default Test Execution (mandatory)

- For Rust code changes, run tests by default.
- This Rust-specific default applies unless user, repo, or provider instructions forbid execution.
- If the user explicitly forbids execution, do not run tests.
- Do not run non-test programs unless explicitly requested.

### 1.1) Relationship to testing policy (mandatory)

- `RUST_TESTING.md` only changes the default execution behavior for Rust code changes (run tests unless forbidden).
- When test strategy, quality gates, or coverage depth are the primary concern, follow `policy-testing` as the governing policy and apply this file as Rust-specific execution guidance.

### 2) When Execution Is Blocked (mandatory)

- If tests cannot be run due to missing toolchain, environment, or excessive runtime, stop and ask.
- Provide the minimal command(s) you would run once unblocked.

### 3) Test Selection

- Prefer the project's existing test commands.
- If no guidance exists, use `cargo test` or the closest scoped test.

### 4) Unsafe Code

- Follow `RUST_UNSAFE.md` requirements (Miri, fuzzing, safety tests).
- If required tools are unavailable, stop and ask.

### 5) Reporting

- Report which tests were run and expected outcome.
- Follow `policy-testing` output requirements when test strategy, quality gates, or coverage depth apply.

---

## Stop and Ask (mandatory)

If any of the following are unclear, stop and ask:

- Whether the user forbids execution
- Which test command the project uses
- Whether unsafe code paths were touched
