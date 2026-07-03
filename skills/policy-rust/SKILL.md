---
name: policy-rust
description: "Apply Rust language and architecture best practices: ownership, error handling, concurrency, unsafe patterns, crate/module/workspace boundaries, trait placement, and dependency direction. Use when Rust implementation details or Rust architecture boundaries are primary. Avoid when the task is language-agnostic or dominated by non-Rust policy concerns."
metadata:
  version: "0.2.0"
---

# Policy Guide

Use the reference files for full rules, examples, and patterns.

When Rust architecture or implementation work crosses API contracts, auth/secrets, test strategy, or runtime/container operations, use optional `policy-api`, `policy-security`, `policy-testing`, or `policy-infra` detail only when those risks are in scope. Keep `policy-rust` as the Rust-specific authority for crate/module/workspace boundaries, trait placement, dependency direction, and async/runtime boundaries.

## Version History

- v0.1.0 (2026-05-08): Initial portable policy release with converged Rust architecture boundaries across crate/module/workspace design and dependency direction.
- v0.1.1 (2026-05-29): Align routing wording with optional-depth policy handoffs.
- v0.2.0 (2026-07-03): Align Rust test execution defaults with global no-execution baseline and fix channel guidance.

## References

- `references/INDEX.md` - Use for navigation and file selection.
