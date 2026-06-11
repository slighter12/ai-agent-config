# Policy Boundaries Matrix

Use this file when multiple policy skills govern the same decision point and ownership or precedence is unclear. Choose the primary owner first, then add optional policy detail only where needed. Do not use this file for ordinary policy overlap when local policy wording already names the owner or handoff.

## Ownership and Trigger Matrix

| Policy | Owns (Primary Authority) | Use When | Avoid When | Optional Policy Detail |
| --- | --- | --- | --- | --- |
| `policy-core` | Global repo-wide baseline behavior, language/output rules, ambiguity handling, scope and execution guardrails | The task needs universal constraints or cross-policy alignment | A single narrower policy fully governs the request details | Use narrower policies only for their owned implementation, contract, security, testing, or infra decisions |
| `policy-go` | Go implementation boundaries, package/layer direction, concurrency/error/database patterns in Go code | The core implementation question is Go-specific | The issue is only wire contract, only deployment/runtime, or only cross-language baseline | Use `policy-api`, `policy-security`, or `policy-testing` detail only when those risks are in scope |
| `policy-rust` | Rust implementation boundaries, ownership/error/concurrency/unsafe patterns | The implementation is Rust-specific | The issue is only wire contract, only deployment/runtime, or only cross-language baseline | Use `policy-api`, `policy-security`, `policy-testing`, or `policy-infra` detail only when those risks are in scope |
| `policy-frontend` | Frontend framework/UI patterns, component/state/routing/i18n/accessibility behavior | The task is primarily frontend framework or UI implementation guidance | The task is backend-only, infra-only, or pure API wire design without UI concerns | Use `policy-api`, `policy-security`, `policy-testing`, or `policy-infra` detail only when those risks are in scope |
| `policy-api` | API wire contract: request/response schemas, validation boundaries, status/error semantics, versioning, API-facing limits | The main risk is external interface correctness or API contract consistency | Internal implementation details dominate and no API boundary decision is needed | Use `policy-security`, `policy-testing`, or language policy detail only when those risks are in scope |
| `policy-security` | Trust boundaries, data exposure/leak prevention, authn/authz, secrets handling, crypto correctness | The task touches confidentiality/integrity, attacker-controlled input, or credential/material handling | No meaningful security impact exists beyond routine coding hygiene | Use `policy-api`, `policy-frontend`, language policy, `policy-infra`, or `policy-testing` detail only when those risks are in scope |
| `policy-testing` | Verification strategy, risk-based test scope, deterministic unit/integration boundaries, fixtures/mocks | The user asks how to verify behavior or what tests are required | The task is purely design/spec discussion without verification planning | Use implementation or contract policy detail for the behavior being verified |
| `policy-infra` | Runtime/deploy/container behavior, environment and config injection, operational wiring | The task is about deployment/runtime configuration and environment-bound behavior | The task is only application logic/API semantics with no runtime/deploy concern | Use `policy-security` or `policy-testing` detail only when those risks are in scope |

## Process Skill Boundary

`execution-harness` is an optional process skill, not a policy authority. Use it to coordinate phases, owners, git/workspace state, verification gates, diff sanity gates, lifecycle gates, and capture candidates when the user asks for structured workflow or approves an orchestrator suggestion.

It does not replace policy ownership:

- `policy-core` owns baseline behavior and shared vocabulary.
- `policy-testing` owns verification strategy and gate selection.
- `project-lifecycle` owns capture-worthy lifecycle signal classification: accepted project
  decision/status synchronization, implementation pivots, documentation drift, capture-worthy
  handoff notes, loop active state, discussion records, and workflow learning triage. Standalone
  review, verification, commit readiness, and handoff packaging remain with their narrowest gate
  owner.
- `conventional-git-flow` owns branch, commit, push, PR, and dependency/version closure.
- `code-review` owns code review modes, including current-diff sanity, full review, security review, and current-diff architecture-risk review.
- `skill-creator` owns approved skill authoring and validation.
- Domain policies still own their implementation, API, security, frontend, infra, language, or runtime decisions.

## Precedence Outcomes (When Policies Overlap)

1. `policy-core` provides baseline constraints when loaded or when cross-policy alignment is needed. Trigger it for detected precedence/conflict, not for ordinary clear policy overlap where a narrower policy fully covers the task.
2. For language/framework implementation details, language/framework policy wins:
   - Go details -> `policy-go`
   - Rust details -> `policy-rust`
   - Frontend framework/UI details -> `policy-frontend`
3. For API wire contracts, `policy-api` is authoritative over language-specific implementation preferences.
4. For trust boundaries, leaks, secrets, and crypto concerns, `policy-security` is authoritative and may tighten any other policy outcome.
5. For verification scope and risk-based evidence, `policy-testing` is authoritative on test strategy and acceptance confidence.
6. For runtime/deploy/config injection concerns, `policy-infra` is authoritative for operational constraints.
7. If two authorities still conflict, keep the stricter safety/correctness constraint and document the assumption briefly.

## Policy Overlap Examples

- **Go handler changes API error payload shape**  
  Primary: `policy-api` when the externally visible response contract changes.  
  Use `policy-go` detail when handler/service placement or Go error mapping is being changed. Use `policy-security` detail only when details may leak sensitive information. Use `policy-testing` detail when contract/regression evidence is required.

- **Rust service adds token encryption and secret loading**  
  Primary: `policy-security` when crypto/secrets trust boundaries are being designed.  
  Use `policy-rust` detail when Rust ownership/error/runtime patterns affect the implementation. Use `policy-infra` detail when secret/config injection mechanics are changed. Use `policy-testing` detail when verification depth or risk gates must be decided.

- **Frontend form introduces auth flow and input validation messaging**  
  Primary: `policy-frontend` when UI framework behavior is the main decision.  
  Use `policy-security` detail when auth/input trust risks are in scope. Use `policy-api` detail when request/response contracts change. Use `policy-testing` detail when UI or boundary verification needs to be planned.

- **Container runtime env var mapping breaks service startup**  
  Primary: `policy-infra` when runtime/config injection mechanics are the main issue.  
  Use `policy-security` detail only when secret classification, exposure, or lifecycle constraints are involved. Use `policy-testing` detail when deployment/startup verification scope must be decided.
