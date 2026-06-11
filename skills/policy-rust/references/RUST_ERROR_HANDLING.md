---
globs: ["*.rs", "**/*.rs"]
description: "Rust error handling - Result, Option, custom errors, error propagation"
---

# RUST_ERROR_HANDLING.md - Rust Error Handling

This file defines error handling rules for Rust projects.
Violating these rules is incorrect output.

---

## HARD RULES

### 1) No Panic (application code)

- Avoid `panic!()`, `unwrap()`, `expect()` in application code
- Only allow in init or truly impossible cases (with comment)
- Tests may use `unwrap()`

### 2) Result vs Option

- Use `Result<T, E>` for fallible operations
- Use `Option<T>` for missing values
- Do not use `Option` to hide errors (use `Result` instead)

### 3) Error Propagation

- Prefer the `?` operator
- Avoid manual `match` with `return Err`
- Library code must define its own error types

### 4) Custom Error Types

- Use `thiserror` or `anyhow` (application code)
- Libraries use `thiserror`, applications may use `anyhow`
- Errors must implement `std::error::Error`

### 5) Error Messages

- Include enough context (file path, user ID, etc.)
- Use `.context()` or `.with_context()` (anyhow)

### 6) main() Error Handling

- `main()` should return `Result<(), E>`
- Use `anyhow::Result` for convenience
- Show a user-friendly error on exit

### 7) Avoid Excessive Nesting

- Use early returns to reduce nesting
- Use `if let` or `match` for simple cases
- Split complex logic into functions

### 8) Error Conversion

- Implement `From<E>` for automatic conversion
- `?` will call `From::from` automatically
- Avoid manual `.map_err(|e| ...)` unless adding context

### 9) Partial Functions

- Avoid `unwrap()` and `expect()` on fallible paths
- Use pattern matching or `ok_or()` to convert to `Result`
- For indexing (`vec[i]`), prefer `.get(i)` returning `Option`

### 10) Stop and Ask (mandatory)

If any of the following are unclear, stop and ask:

- Whether to use `Result` or `Option`
- How to design error types
- Whether custom errors are required
- How to handle multiple error types

---

## Recommended Crates

- `thiserror` - custom error types
- `anyhow` - application error handling
- `color-eyre` - better error reports (dev)

---

Violating these rules is incorrect output.
