---
name: policy-api
description: "Apply API boundary rules for layered architectures: contracts, validation, versioning, and error semantics. Keywords: API contract, HTTP status, error response, versioning, request ID, rate limiting, CORS header contract. Use when API contracts and boundary behavior are primary. Avoid when internal DB schemas or service internals are primary."
metadata:
  version: "0.1.1"
---

# Policy Guide

Use the referenced policy file for full rules. Keep output aligned with these rules.
Use when wire-level API contract behavior is primary: status codes, response schema, machine error code contract, request/response shape consistency, and `request_id` propagation.
Avoid when the primary question is data leak prevention or auth/authz policy; in those cases, use optional `policy-security` detail and let security rules constrain what API fields/messages are safe to return.

## Version History

- v0.1.0 (2026-05-08): Initial portable policy release with converged API ownership boundaries for contract semantics, versioning, and error-response behavior.
- v0.1.1 (2026-05-29): Align routing wording with optional-depth policy handoffs.

## References

- `references/INDEX.md` - Use for navigation and file selection.
- `references/API.md`
