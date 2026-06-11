---
description: "REST API design standards - status codes, error format, contracts"
---

# API.md - API Design Standards

This file defines mandatory REST API design standards.
It applies to any project that exposes HTTP APIs.
This policy is domain-scoped and should only be applied when API contract or handler behavior is in scope.
Ownership note: API policy owns wire contract semantics (status mapping, schema envelope, `request_id`, and stable machine-readable response shape). Use when these contracts are being defined or changed. Avoid using API policy alone to decide leak redaction/auth boundaries; use optional security policy detail for disclosure safety constraints.

Violating these rules is incorrect output.

---

## Table of Contents

- 1) HTTP Status Codes (mandatory)
- 1) Error Response Format (mandatory)
- 1) Success Response Format (mandatory)
- 1) API Contract (mandatory)
- 1) Request ID Propagation (mandatory)
- 1) API Versioning
- 1) Rate Limiting
- 1) CORS Header Contract
- 1) Security Boundary Handoff
- 1) When Uncertain (mandatory)

## 1) HTTP Status Codes (mandatory)

Use standard HTTP status codes:

| Status | Usage |
|--------|-------|
| **2xx Success** ||
| 200 OK | Successfully retrieved a resource or executed an action |
| 201 Created | Successfully created a resource |
| 204 No Content | Successfully executed with no body (e.g., DELETE) |
| **4xx Client Errors** ||
| 400 Bad Request | Request format error or validation failure |
| 401 Unauthorized | Not authenticated (missing or invalid token/session) |
| 403 Forbidden | Authenticated but not authorized |
| 404 Not Found | Resource does not exist |
| 409 Conflict | Resource state conflict (e.g., duplicate create) |
| 422 Unprocessable Entity | Syntax is valid but semantic validation fails |
| 429 Too Many Requests | Rate limit exceeded |
| **5xx Server Errors** ||
| 500 Internal Server Error | Unexpected server error |
| 503 Service Unavailable | Service temporarily unavailable |

Do not:

- Return 200 for all errors with an error field
- Use non-standard status codes
- Confuse 401 and 403

---

## 2) Error Response Format (mandatory)

Unified error format:

```json
{
  "error": {
    "code": "MACHINE_READABLE_CODE",
    "message": "請檢查輸入內容",
    "details": {}
  },
  "meta": {
    "request_id": "req_abc123"
  }
}
```

Rules:

- `code`: Required. Stable machine-readable error code in upper snake case (e.g., `VALIDATION_FAILED`). `code` is the primary contract for client logic and branching.
- `message`: Optional fallback user-facing text. Default to Traditional Chinese when locale behavior is unspecified by project rules.
- `message_key`: Optional i18n key (for example `errors.validation_failed`) for clients that localize by key/code mapping.
- `details`: Optional additional context (e.g., which field failed validation)
- `request_id`: Required for tracing and must appear in `meta.request_id`

For i18n-capable clients, prefer `code` (and optionally `message_key`) as the stable contract. Do not require clients to branch on `message`.

Details rules (mandatory):

- Only include `details` for client-correctable 4xx errors (especially 400/422)
- Do not include `details` for 5xx or auth/permission errors unless explicitly required by the endpoint contract
- `details` must be limited to contract-defined, client-correctable fields.
- Never include secrets, tokens, credentials, PII, stack traces, SQL queries, hostnames, file paths, schema internals, or internal error messages in `details`.
- If sensitive-data disclosure safety is unclear, apply the stricter redaction rule and use optional `policy-security` detail for trust-boundary decisions.
- Use the reference schema below by default; additional fields are allowed only if explicitly documented in the endpoint contract
- If the project or endpoint contract does not specify otherwise, use this standard format and do not ask whether to include `meta` or `details`

Reference schema (default):

```json
{
  "error": {
    "code": "VALIDATION_FAILED",
    "message": "請檢查輸入內容",
    "details": {
      "fields": [
        {
          "field": "email",
          "code": "INVALID_FORMAT",
          "reason": "invalid_format",
          "allowed_values": "example@domain.com",
          "constraint": "must be a valid email address"
        }
      ]
    }
  },
  "meta": {
    "request_id": "req_abc123"
  }
}
```

Reference keys in `error.details.fields[]`: `field`, `code`, `reason`, `allowed_values`, `constraint`

Error code examples:

- `VALIDATION_FAILED` - Input validation failed
- `RESOURCE_NOT_FOUND` - Resource not found
- `PERMISSION_DENIED` - Permission denied
- `RATE_LIMIT_EXCEEDED` - Rate limit exceeded
- `RESOURCE_CONFLICT` - Resource conflict (e.g., duplicate key)

---

## 3) Success Response Format (mandatory)

Unified success format:

```json
{
  "data": {},
  "meta": {
    "request_id": "req_abc123"
  }
}
```

Pagination example:

```json
{
  "data": [],
  "meta": {
    "request_id": "req_abc123",
    "pagination": {
      "page": 1,
      "per_page": 20,
      "total": 100,
      "total_pages": 5
    }
  }
}
```

---

## 4) API Contract (mandatory)

Each endpoint must define:

- HTTP method
- URL path
- Request parameters (path/query/body)
- Request schema (data types, required fields)
- Response schema (success + error)
- Status codes (which scenarios return which codes)

---

## 5) Request ID Propagation (mandatory)

- Every request must have a unique `request_id`.
- If the client provides `X-Request-Id`, use it; otherwise generate a new one.
- `request_id` must:
  - Appear in the response body as `meta.request_id` (mandatory)
  - Optionally be returned in the response header `X-Request-Id` (does not replace body)
  - Appear in all related logs
  - Trace the full request chain

---

## 6) API Versioning

Recommended approaches (choose one):

- URL path: `/v1/users`, `/v2/users`
- Header: `Accept: application/vnd.api.v1+json`

Rules:

- Do not make breaking changes within the same version
- Adding optional fields is not breaking
- Removing or renaming fields requires a new version

---

## 7) Rate Limiting

Provide rate limit metadata in response headers:

```
X-RateLimit-Limit: 100
X-RateLimit-Remaining: 95
X-RateLimit-Reset: 1609459200
```

When exceeded, return 429 with a `Retry-After` header.

---

## 8) CORS Header Contract

API policy ownership in this section is limited to contract-visible header mechanics:

- Document which CORS headers are returned and under what request/response conditions.
- Keep header behavior stable and testable at the contract boundary (for example: origin echoing behavior, preflight response headers, and credential/header/method negotiation fields).
- Treat CORS header behavior as part of the API contract when clients depend on it.

API policy does not own trust/security decisions for CORS. Use optional `policy-security` detail for:

- Trusted-origin allowlist decisions
- Credentialed cross-origin security posture
- Any CORS decision that changes security boundaries or exposure risk

---

## 9) Security Boundary Handoff

API policy ownership in this section is limited to contract-visible security implications:

- Preserve contract-safe status/header/envelope behavior for auth/security outcomes (for example: correct 401 vs 403 mapping, stable error envelope, and `request_id` propagation).
- Keep security-relevant contract fields explicit and stable when endpoints require them (for example: authentication-related response headers or documented error codes).

API policy is not the owner of security controls. Do not expose sensitive/internal details while preserving the API contract. For any of the following, use optional `policy-security` detail:

- HTTPS/TLS requirements and transport protections
- Token/session/authentication mechanism selection
- Authorization controls and permission boundaries
- Secret handling, redaction, storage, rotation, and lifecycle constraints
- Query/header/body leakage controls for sensitive data

---

## 10) When Uncertain (mandatory)

If any of the following are unclear, stop and ask:

- Which status code to use
- Error code naming conventions
- API versioning strategy
- Impact scope of a breaking change

---

Violating these rules is incorrect output.
