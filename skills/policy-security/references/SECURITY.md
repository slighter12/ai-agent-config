---
description: "Security baseline for all projects - secrets, PII, auth, crypto"
---

# SECURITY.md - Security Baseline

This file defines mandatory security rules for all projects.
This policy is domain-scoped and should be applied when security-relevant changes are in scope.
Ownership note: Security policy owns trust boundaries (authn/authz), secret/PII protection, cryptographic safety, and leak prevention in logs and responses. Use when exposure risk or access control is in scope. Avoid replacing API contract ownership for response envelope/status mapping; use the project API contract or optional `policy-api` detail for envelope structure, and enforce security redaction within that contract.

Violating these rules is incorrect output.

---

## Table of Contents

- 1) Secrets Management (hard rules)
- 1) Logging Safety (hard rules)
- 1) Error Message Safety (hard rules)
- 1) Privacy By Design (high-level rules)
- 1) Frontend Security (hard rules)
- 1) Authentication and Authorization (hard rules)
- 1) Cryptography (hard rules)
- 1) Input Validation (hard rules)
- 1) Dependency Security
- 1) When Uncertain (mandatory)

## 1) Secrets Management (hard rules)

- Never hardcode secrets or credentials.
- Never log secrets/tokens/PII (unless explicitly approved).
- Never commit secrets to Git.
- Security owns secret classification and handling constraints across exposure, rotation, and lifecycle.
- Policy-infra owns secret injection and storage mechanisms.

Do not:

- API keys, tokens, passwords
- Private keys, certificates
- Database credentials
- OAuth secrets

Required handling constraints:

- Classify secrets and apply least exposure to code, logs, and outputs.
- Define rotation expectations and revocation response for leaked credentials.
- Enforce lifecycle controls for creation, distribution, usage, and decommissioning.

---

## 2) Logging Safety (hard rules)

Never log:

- Secrets, tokens, API keys
- Passwords (even hashed)
- PII (personally identifiable information)
- Credit card numbers, government IDs
- Full request/response bodies (may contain sensitive data)

Allowed to log:

- request_id
- user_id (non-sensitive identifier)
- HTTP status codes
- Error codes (no sensitive details)
- Performance metrics

---

## 3) Error Message Safety (hard rules)

User-facing errors must not leak internal details:

Do not expose:

- Stack traces
- SQL queries
- Internal hostnames or IP addresses
- File paths
- Database schema details
- Internal error messages

Do:

- Use stable machine error codes
- Provide user-friendly messages
- Log full errors internally at the boundary (do not return them); avoid duplicate warn/error logs
- Ensure request traceability through the API contract (for example, documented request identifiers)
- If returning `details`, only include safe, client-correctable fields and ensure the endpoint contract documents any additional fields
- Use the project API contract, or optional `policy-api` detail, for response envelope and field-shape decisions.

Example:

```json
{
  "error": {
    "code": "VALIDATION_FAILED",
    "message": "請檢查輸入內容"
  }
}
```

This example only shows the safe error payload content. Use the project's existing API contract, or optional `policy-api` detail, for the complete response envelope, metadata fields, and request identifier shape.

---

## 4) Privacy By Design (high-level rules)

- Minimize sensitive data collected, stored, logged, copied into prompts, or sent to third-party tools.
- Prefer synthetic or redacted examples over real PII, tokens, private files, or full request/response bodies.
- When analytics, tracking, external processors, or retention behavior are in scope, surface the privacy risk and ask for project-specific requirements instead of inventing compliance policy.
- Keep privacy guidance high-level unless the project has explicit legal, regulatory, retention, or consent requirements.

---

## 5) Frontend Security (hard rules)

- The frontend must never receive or store secrets.
- Frontend authorization checks are UX only; the backend is the source of truth.
- Never run sensitive business logic on the client (can be bypassed).
- Use HTTPS in production.
- Set an appropriate CSP (Content Security Policy).

---

## 6) Authentication and Authorization (hard rules)

- All API endpoints require auth by default (unless explicitly public).
- Use existing auth mechanisms (do not roll your own).
- Sessions/tokens must expire.
- Sensitive operations require additional verification (e.g., re-auth, 2FA).

---

## 7) Cryptography (hard rules)

- Do not implement your own crypto algorithms.
- Use well-reviewed libraries.
- Encrypt sensitive data at rest.
- Use strong password hashes (bcrypt, Argon2, scrypt).

Do not:

- Use MD5 or SHA1 for password hashing
- Write custom crypto
- Use fixed encryption keys

---

## 8) Input Validation (hard rules)

- Never trust user input.
- Validate all inputs (type, format, range, length).
- Use parameterized queries (prevent SQL injection).
- Validate and sanitize HTML input (prevent XSS).
- Validate file uploads (type, size, content).

---

## 9) Dependency Security

- Regularly update dependencies to patch vulnerabilities.
- Avoid using packages with known vulnerabilities.
- Review permissions of third-party libraries.

---

## 10) When Uncertain (mandatory)

If any of the following are unclear, stop and ask:

- Auth/authz boundaries
- PII handling rules
- Compliance requirements (GDPR, CCPA, etc.)
- Secret classification, lifecycle, disclosure, or rotation constraints
- Secret storage/passing mechanisms (use optional `policy-infra` detail)
- Cryptography requirements

---

Violating these rules is incorrect output.
