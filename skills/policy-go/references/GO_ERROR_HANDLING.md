---
globs: ["*.go", "**/*.go"]
description: "Go error handling - panic rules, error wrapping, logging"
---

# GO_ERROR_HANDLING.md - Go Error Handling

This file defines error handling rules for Go projects.
Violating these rules is incorrect output.

---

## HARD RULES

### 1) No Panic (application code)
- Application code must not panic
- Only allow panic during initialization (`init()` or early `main()`) when execution cannot continue
- Even "impossible" cases should prefer returning errors

### 2) Recover Limits
- Do not use `recover()` to hide errors
- Use only at clear boundaries (e.g., HTTP middleware)
- After recover, you must log the full stack trace

### 3) Handle Errors Explicitly
- Errors must be returned or explicitly handled
- Do not ignore errors (unless there is a clear reason and `_` is used)
- Use `err` as the error variable name

### 4) Error Wrapping Rules
- Errors must be wrapped with context unless the layer rule below says otherwise
- If the project already uses `pkg/errors`, prefer it for message context and stack traces (`errors.WithMessage`, `errors.WithStack`) per the layer rules
- If the project does not use `pkg/errors`, prefer `fmt.Errorf("...: %w", err)` to preserve the original error
- Using `%v` is an exception and must be justified (e.g., to avoid leaking internal details)
  - Only when you are sure it is a specific error type or when responding externally in middleware
- Do not introduce new error dependencies without explicit approval (including `github.com/cockroachdb/errors`)
- If stack traces are required but no stack package exists, stop and ask whether to add `pkg/errors`

### 5) Layer Mapping (mandatory)
This policy uses conceptual layers:
- Interface/Delivery
- Application/Usecase
- Domain
- Infrastructure

Each project must map its packages to these layers in its architecture docs.
If no mapping exists, use this default mapping by naming convention:
- Interface/Delivery: handler, controller, router, transport
- Application/Usecase: service, usecase, app
- Domain: domain, entity, model
- Infrastructure: repository, dao, adapter, gateway, persistence

If a package does not clearly map to a conceptual layer, stop and ask.

### 6) Layered Error Handling (mandatory)
Apply the rules below using the conceptual layers defined above:
- **Infrastructure layer**: add stack trace only; do not add extra message
  - If using `pkg/errors`: `return errors.WithStack(err)`
  - If no stack package exists: return `err` unchanged
- **Domain layer**: convert to domain errors and add business context
  - Example: `return domainerrors.NewDatabaseExecuteError(err, "failed to create user")`
  - Domain errors must preserve the original cause (support `Unwrap()` or `%w`) so `errors.Is/As` continues to work
- **Usecase layer**: add business message without adding a second stack trace
  - If using `pkg/errors`: `return errors.WithMessage(err, "failed to execute user registration transaction")`
  - If no stack package exists: `return fmt.Errorf("failed to execute user registration transaction: %w", err)`

### 7) No Double Handling (hard rule)
- Do not both log and return an error at the same layer
- Do not log the same error multiple times in one layer (e.g., warn + error)
- Option A (default): log once at the boundary (Interface/Delivery); internal functions only return
- Option B: log internally and do not return (less common)

### 8) Error Checks
- Use `errors.Is(err, target)` to check error values
- Use `errors.As(err, &target)` to check error types and extract details
- Do not compare errors directly (`err == target`)

### 9) Custom Error Rules
- Define common errors as `var Err...`
- Define error types when extra data is needed (implement `Error()`)
- Use `errors.New()` for error variables

### 10) Client-Facing Errors (mandatory)
- Use stable machine error codes (uppercase English)
- Provide user-friendly messages in Traditional Chinese
- Do not expose internal details (stack traces, SQL, hostnames)

### 11) Error Message Format (internal errors)
- Internal error messages should be in English
- Start with lowercase (unless proper nouns)
- Do not end with punctuation
- Provide enough context (file path, user ID, etc.)

### 12) Stop and Ask (mandatory)
If any of the following are unclear, stop and ask:
- Which layer should log errors
- Whether a custom error type is needed
- How to map project packages to conceptual layers
- How to map internal errors to user-friendly messages
- Whether to add `pkg/errors` when stack traces are required but no stack package exists

---

## Detailed Guidance

Need full examples, patterns, or pitfalls? Use `$policy-go`.

---

Violating these rules is incorrect output.
