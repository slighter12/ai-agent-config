---
name: policy-security
description: "Apply security rules: authn/authz, secrets, input validation, crypto usage, and threat modeling. Keywords: authn, authz, secrets, crypto, input validation, XSS, CSRF, logging safety, dependency security. Use when threat exposure or security controls are in scope. Avoid when no security impact."
metadata:
  version: "0.1.1"
---

# Policy Guide

Use the referenced policy file for full rules. Keep output aligned with these rules.
Use when trust boundaries, leak prevention, authn/authz, secret handling, crypto usage, or sensitive logging/exposure risks are primary.
Avoid when only wire contract mechanics are in scope (for example exact status code or envelope shape) unless security constraints affect what can be disclosed; use optional `policy-api` detail for contract structure and keep security as the disclosure boundary owner.

## Version History

- v0.1.0 (2026-05-08): Initial portable policy release with converged security ownership boundaries for disclosure control, trust boundaries, and sensitive data handling.
- v0.1.1 (2026-05-29): Align routing wording with optional-depth policy handoffs.

## References

- `references/INDEX.md` - Use for navigation and file selection.
- `references/SECURITY.md`
