---
description: "Risk-based testing strategy for all projects"
---

# TESTING.md - Testing Strategy

This file defines a risk-based testing strategy.
Applies to all projects and languages.
This policy is domain-scoped and should be applied when test strategy or validation scope is in focus.
Ownership note: Testing policy owns risk gates, validation strategy, and expected verification depth/coverage for a change. Use when deciding what evidence is required for confidence. Avoid treating this file as language runtime execution guidance; use optional language-specific policy detail for tool/command conventions while keeping risk classification and verification scope in testing policy.

Violating these rules is incorrect output.

---

## Table of Contents

- 1) Risk Level Classification (mandatory)
- 1) Testing Requirements (by risk)
- 1) Test Types
- 1) Test Naming Rules
- 1) TDD / Test-First Workflow
- 1) Test Data Management
- 1) Mocking Rules
- 1) Coverage Guidance
- 1) Verification Output Format (when code changes)
- 1) Test Environments
- 1) Default Behavior for No Test Execution
- 1) When Uncertain (mandatory)

## 1) Risk Level Classification (mandatory)

Every change must be classified by risk level:

| Risk Level | Definition | Examples |
|-----------|------------|----------|
| **Low** | Docs, copy, styling, refactors without behavior changes | README updates, CSS tweaks, renames |
| **Medium** | New features, non-critical logic, API behavior changes | New endpoint, validation, new UI page |
| **High** | Auth/authz, money/balance/orders, concurrency, migrations, data loss risk | Login logic, payments, DB migration, goroutine management |

---

## 2) Testing Requirements (by risk)

### Low Risk

- Testing: optional
- Must provide:
  - Manual verification checklist
  - Expected behavior notes

### Medium Risk

- Testing: recommended (unless user explicitly forbids)
- Minimum coverage:
  - Happy path
  - Validation failure
  - One meaningful edge case
- If tests cannot be run:
  - Provide a detailed manual verification checklist
  - List assumptions and risks

### High Risk

- Testing: mandatory (unless user explicitly forbids)
- Must cover:
  - All success paths
  - All error paths
  - Boundary conditions
  - Concurrency cases (if applicable)
  - Rollback strategy (if applicable)
- If tests cannot be run:
  - Must provide:
    - Detailed manual verification steps
    - Risk assessment
    - Rollback plan
    - Monitoring signals

---

## 3) Test Types

### Unit Tests

- Test a single function/method
- Mock external dependencies
- Fast to run (< 1 second)

### Integration Tests

- Test interactions between modules
- Can use real databases/services (test containers)
- Slower but closer to real scenarios

### E2E Tests

- Test complete user flows
- Usually slowest but highest confidence
- Prioritize critical paths

---

## 4) Test Naming Rules

Clearly describe what is being tested and expected:

Good names:

```go
func TestCreateUser_WithValidInput_ReturnsUser(t *testing.T)
func TestCreateUser_WithDuplicateEmail_ReturnsConflictError(t *testing.T)
```

Bad names:

```go
func TestCreateUser(t *testing.T)
func TestCase1(t *testing.T)
```

---

## 5) TDD / Test-First Workflow

Use this section when the user asks for TDD, test-first work, red-green-refactor, or a vertical validation loop.

Default loop:

1. Pick the smallest behavior slice that proves user value or contract behavior.
2. Write or describe the failing test first, using existing test conventions.
3. Implement only enough code to make that test pass.
4. Refactor only after the behavior is covered.
5. Repeat for the next slice.

Boundaries:

- Do not add tests when the repo policy or user request forbids test creation.
- Do not invent a new framework, test runner, fixture system, or dependency unless explicitly approved.
- Prefer integration-facing tests when the risk is contract behavior; prefer unit tests for isolated logic.
- Keep mocks near real trust boundaries: external services, time, filesystem, network, or slow dependencies.
- If tests are not run, provide the exact manual verification path and the expected failing/passing evidence.

---

## 6) Test Data Management

Principles:

- Each test should be independent
- Separate setup and teardown clearly
- Use meaningful test data (avoid `foo`, `bar`)
- Do not hardcode sensitive data (even in tests)

---

## 7) Mocking Rules

When to mock:

- External API calls
- Database operations (unit tests)
- Third-party services (payments, SMS)
- Time-dependent logic (use a controllable clock)

When not to mock:

- Integration tests should use real dependencies when possible
- Simple utility functions

---

## 8) Coverage Guidance

No fixed percentage requirement, but:

- High-risk code must have tests
- Critical business logic should have high coverage
- Coverage is a signal, not a goal (meaningful tests > high coverage)

---

## 9) Verification Output Format (when code changes)

Must provide:

1. Change summary
   - Which files were modified
   - Risk level

2. Verification plan
   - Tests (if any):
     - Command to run (e.g., `go test ./...`)
     - Expected outcome (all tests pass)
   - Manual verification (if no tests):
     - HTTP request examples
     - Expected status codes and responses
     - Edge case checklist

3. Risks and assumptions
   - List assumptions (e.g., "Assumes user is authenticated")
   - Known limitations

4. Rollback strategy (required for high risk)
   - How to revert changes
   - Data backup plan (if applicable)

If another active repo or task rule requires a different response structure, embed these items in that structure without duplicating sections.

---

## 10) Test Environments

Recommended:

- Local development: Docker Compose or test containers
- CI/CD: run all tests automatically
- Staging: mirror production
- Production: monitoring and alerting (do not test in production)

---

## 11) Default Behavior for No Test Execution

Unless the user explicitly requests otherwise or a more specific policy requires execution:

- Do not create test files
- Do not run tests
- Do not run programs

Scope boundary:

- This default is owned by testing policy as a verification-strategy baseline only.
- This section does not own language/runtime-specific execution conventions (test command selection, toolchain defaults, or language-specific run expectations).
- For language-specific execution defaults, use optional detail from the relevant language policy (for example `policy-go`, `policy-rust`, or frontend policy files).

But you must still provide:

- Manual verification checklist
- Risk assessment (if correctness depends on runtime behavior)

---

## 12) When Uncertain (mandatory)

If any of the following are unclear, stop and ask:

- Risk level classification (Low/Medium/High)
- Whether tests are required (user forbids or not)
- Test framework conventions (existing project patterns)
- How to configure the test environment (DB, mock services)

---

Violating these rules is incorrect output.
