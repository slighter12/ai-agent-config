---
name: policy-go
description: "Apply Go language best practices for concurrency, error handling, database usage, and Clean Architecture or layered architecture boundaries. Use when Go implementation details, handler/service/repository boundaries, delivery/usecase/repository flow, dependency direction, package boundaries, or Go-specific architecture rules are primary. Avoid when the task is not Go-specific or is only API contracts, testing strategy, security, infrastructure, frontend, or repo-domain guidance better handled by another policy or domain skill."
metadata:
  version: "0.1.1"
---

# Policy Guide

Use the reference files for full rules, examples, and patterns.

For Go tasks, use optional `policy-api` detail only when HTTP contracts, status codes, request/response schemas, or API error response semantics are in scope. Use optional `policy-security` detail only when authn/authz, secrets, PII, crypto, trust boundaries, or security validation are in scope. Use optional `policy-testing` detail only when tests, coverage, fixtures, mocks, or unit/integration/e2e verification strategy are in scope.

## Version History

- v0.1.0 (2026-05-08): Initial portable policy release with converged Go clean/layered architecture routing and explicit ownership boundaries across API, security, and testing.
- v0.1.1 (2026-05-29): Align routing wording with optional-depth policy handoffs.

## References

- `references/INDEX.md` - Use for navigation and file selection.
